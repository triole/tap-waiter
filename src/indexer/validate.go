package indexer

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/triole/logseal"
)

const (
	// MaxURLPathLength limits request URI size to prevent DoS
	MaxURLPathLength = 2048
	
	// MaxQueryParamLength limits individual query parameters
	MaxQueryParamLength = 1024
	
	// MaxFilterStringLength limits total filter parameter size
	MaxFilterStringLength = 2048
	
	// MaxSortByLength limits sortby parameter size
	MaxSortByLength = 256
	
	// MaxFilterPrefixLength limits field path length
	MaxFilterPrefixLength = 256
	
	// MaxFilterSuffixLength limits individual value length
	MaxFilterSuffixLength = 512
	
	// MaxSuffixCount limits number of suffix values
	MaxSuffixCount = 50
	
	// MaxRegexSuffixLength stricter limit for regex patterns
	MaxRegexSuffixLength = 200
	
	// RegexCompilationTimeout prevents ReDoS during regex compilation
	RegexCompilationTimeout = 100 * time.Millisecond
)

// containsPathTraversalPrefix checks for path traversal patterns in field identifiers.
// Unlike containsPathTraversal, this allows dots and brackets which are valid in JSON field paths.
// It only blocks:
// - ".." sequences (path traversal)
// - Backslashes (Windows-style paths)
// - Forward slashes in contexts where they're dangerous
func containsPathTraversalPrefix(s string) bool {
	// Block path traversal sequences
	if strings.Contains(s, "..") {
		return true
	}
	// Block Windows-style path separators
	if strings.Contains(s, "\\") {
		return true
	}
	// Block forward slashes that indicate directory traversal (but allow them in valid field names like "body.content")
	if strings.Contains(s, "/") {
		return true
	}
	// Block other path-like patterns
	if strings.Contains(s, ":\\") || strings.HasPrefix(s, "/") {
		return true
	}
	return false
}

// containsPathTraversal checks for path traversal in general contexts (URLs, etc).
// This is more aggressive and blocks all path separators.
func containsPathTraversal(s string) bool {
	if strings.Contains(s, "..") ||
		strings.Contains(s, "\\") ||
		strings.Contains(s, "/") {
		return true
	}
	return false
}

// checkInputLength validates that input doesn't exceed maximum allowed length.
// Logs a warning if exceeded and returns false.
func checkInputLength(s string, maxLen int, fieldName string, lg logseal.Logseal) bool {
	if len(s) > maxLen {
		lg.Warn(
			"input length exceeds limit",
			logseal.F{
				"field":     fieldName,
				"length":    len(s),
				"max":       maxLen,
				"input":     s[:min(len(s), 100)],
			},
		)
		return false
	}
	return true
}

// isValidSortBy validates sortby parameters against a whitelist to prevent path traversal.
// Only allows specific safe values: "created", "lastmod", "size", or content field paths.
func (ind *Indexer) isValidSortBy(sortBy string) bool {
	// Check length limit first
	if !checkInputLength(sortBy, MaxSortByLength, "sortBy", ind.Lg) {
		return false
	}

	// Block path traversal attempts immediately
	if containsPathTraversal(sortBy) {
		ind.Lg.Warn(
			"sortBy parameter blocked: path traversal detected",
			logseal.F{"sortBy": sortBy},
		)
		return false
	}

	// Whitelist of safe built-in sort options
	allowedSorts := []string{
		"",
		"created",
		"lastmod",
		"size",
	}

	// Check if it's a built-in option
	for _, allowed := range allowedSorts {
		if sortBy == allowed {
			return true
		}
	}

	// For content field paths, they must start with "body." or "front_matter."
	// and contain only safe characters
	if strings.HasPrefix(sortBy, "body.") || strings.HasPrefix(sortBy, "front_matter.") {
		// Extract the field path part after the prefix
		fieldPath := sortBy
		if strings.HasPrefix(sortBy, "body.") {
			fieldPath = sortBy[5:] // Remove "body." prefix
		} else if strings.HasPrefix(sortBy, "front_matter.") {
			fieldPath = sortBy[15:] // Remove "front_matter." prefix
		}

		// Check that field path contains only safe characters
		// Allow: letters, numbers, underscores, dots, brackets [array access], dashes
		validFieldChars := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.\[\] -]*$`)
		if !validFieldChars.MatchString(fieldPath) {
			ind.Lg.Warn(
				"sortBy field path contains invalid characters",
				logseal.F{"sortBy": sortBy, "fieldPath": fieldPath},
			)
			return false
		}

		return true
	}

	ind.Lg.Warn(
		"invalid sortBy value not in whitelist",
		logseal.F{"sortBy": sortBy},
	)
	return false
}

// isValidFilterParams validates filter parameters to prevent injection attacks.
func (ind *Indexer) isValidFilterParams(params FilterParams) bool {
	// Check suffix count limit first
	if len(params.Suffix) > MaxSuffixCount {
		ind.Lg.Warn(
			"filter suffix count exceeds limit",
			logseal.F{
				"count": len(params.Suffix),
				"max":   MaxSuffixCount,
			},
		)
		return false
	}

	// Validate Prefix - must be safe field path
	if !checkInputLength(params.Prefix, MaxFilterPrefixLength, "filter.prefix", ind.Lg) {
		return false
	}

	// Block path traversal in prefix first
	if containsPathTraversalPrefix(params.Prefix) {
		ind.Lg.Warn(
			"path traversal detected in filter prefix",
			logseal.F{"prefix": params.Prefix},
		)
		return false
	}

	// Validate operators - whitelist approach
	allowedOperators := map[string]bool{
		"==":    true,
		"===":   true,
		"!=":    true,
		"!==":   true,
		"=~":    true,
		"==~":   true,
		"=*":    true,
		"!=*":   true,
	}

	if _, ok := allowedOperators[params.Operator]; !ok {
		ind.Lg.Warn(
			"filter operator not in whitelist",
			logseal.F{"operator": params.Operator},
		)
		return false
	}

	// Validate Suffix entries - ensure no path traversal or injection
	for i, suffix := range params.Suffix {
		// Check individual suffix length
		if !checkInputLength(suffix, MaxFilterSuffixLength, "filter.suffix", ind.Lg) {
			return false
		}

		// Block path traversal patterns
		if containsPathTraversalPrefix(suffix) {
			ind.Lg.Warn(
				"path traversal detected in filter suffix",
				logseal.F{"suffix_index": i, "suffix": suffix[:min(len(suffix), 100)]},
			)
			return false
		}

		// Additional validation for regex patterns
		if strings.HasPrefix(params.Operator, "=~") || strings.HasPrefix(params.Operator, "==~") {
			// Block potentially dangerous regex patterns that could cause ReDoS
			// Limit pattern length to prevent memory exhaustion
			if len(suffix) > MaxRegexSuffixLength {
				ind.Lg.Warn(
					"regex filter suffix too long",
					logseal.F{"length": len(suffix), "max": MaxRegexSuffixLength},
				)
				return false
			}

			// Attempt to compile the regex with timeout to prevent ReDoS
			ctx, cancel := context.WithTimeout(context.Background(), RegexCompilationTimeout)
			defer cancel()

			done := make(chan bool, 1)
			var compileErr error

			go func() {
				_, compileErr = regexp.Compile(suffix)
				done <- true
			}()

			select {
			case <-done:
				if compileErr != nil {
					ind.Lg.Warn(
						"invalid regex in filter suffix",
						logseal.F{"suffix_index": i, "suffix": suffix, "error": compileErr},
					)
					return false
				}
			case <-ctx.Done():
				ind.Lg.Warn(
					"regex compilation timeout (possible ReDoS)",
					logseal.F{"suffix_index": i, "suffix": suffix, "timeout": RegexCompilationTimeout},
				)
				return false
			}
		}
	}

	// Ensure prefix starts with valid field identifiers
	if params.Prefix != "" {
		// Allow: body.*, front_matter.*, or simple field names
		if !strings.HasPrefix(params.Prefix, "body.") &&
			!strings.HasPrefix(params.Prefix, "front_matter.") {
			// For simple field names, ensure they start with alphanumeric
			if len(params.Prefix) == 0 || !((params.Prefix[0] >= 'a' && params.Prefix[0] <= 'z') ||
				(params.Prefix[0] >= 'A' && params.Prefix[0] <= 'Z') ||
				params.Prefix[0] == '_') {
				ind.Lg.Warn(
					"filter prefix doesn't start with valid field identifier",
					logseal.F{"prefix": params.Prefix},
				)
				return false
			}
		}
	}

	return true
}
