package indexer

import (
	"regexp"
	"strings"

	"github.com/triole/logseal"
)

// Input length limits for security and DoS prevention
const (
	MaxURLPathLength      = 2048   // Maximum URL path length
	MaxQueryParamLength   = 1024   // Maximum length for individual query parameter values
	MaxFilterStringLength = 2048   // Maximum total filter string length
	MaxFilterPrefixLength = 256    // Maximum filter prefix (field path) length
	MaxFilterSuffixLength = 512    // Maximum individual suffix value length
	MaxSortByLength       = 256    // Maximum SortBy parameter length
	MaxSuffixCount        = 50     // Maximum number of suffix values in filter
)

// containsPathTraversal checks if a string contains path traversal sequences.
// Returns true if the input contains any form of path traversal.
func containsPathTraversal(s string) bool {
	// Check for common path traversal patterns
	if strings.Contains(s, "..") || 
	   strings.Contains(s, "\\") || 
	   strings.Contains(s, "/") {
		return true
	}
	return false
}

// checkInputLength validates that an input string does not exceed the maximum allowed length.
// Returns true if the input is within acceptable limits.
func checkInputLength(s string, maxLen int, fieldName string, lg logseal.Logseal) bool {
	if len(s) > maxLen {
		lg.Warn(
			"input length exceeds limit",
			logseal.F{
				"field":     fieldName,
				"length":    len(s),
				"max":       maxLen,
				"input":     s[:min(len(s), 100)], // Log first 100 chars only
			},
		)
		return false
	}
	return true
}

// isValidSortBy validates the SortBy parameter against a whitelist to prevent path traversal.
// Only allows specific safe values: "created", "lastmod", "size", or content field paths.
func (ind *Indexer) isValidSortBy(sortBy string) bool {
	// Check length limit first
	if !checkInputLength(sortBy, MaxSortByLength, "sortBy", ind.Lg) {
		return false
	}

	// Block path traversal attempts immediately
	if containsPathTraversal(sortBy) {
		ind.Lg.Warn(
			"path traversal detected in sortBy parameter",
			logseal.F{"sortBy": sortBy},
		)
		return false
	}

	// Whitelist of safe built-in sort options
	safeBuiltIns := map[string]bool{
		"":        true,
		"created": true,
		"lastmod": true,
		"size":    true,
	}

	if safeBuiltIns[sortBy] {
		return true
	}

	// For content field paths, validate the format
	// Must start with "front_matter." or "body." and contain only safe characters
	// This prevents path traversal like "../../../etc/passwd"
	if strings.HasPrefix(sortBy, "front_matter.") || strings.HasPrefix(sortBy, "body.") {
		// Allow only alphanumeric, underscore, dot, and bracket characters for JSON path access
		// Pattern: field_name or field_name[nested_field]
		// Block any path traversal patterns in the field path as well
		if containsPathTraversal(sortBy) {
			return false
		}
		matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_.\[\]]+$`, sortBy)
		return matched
	}

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

	if strings.HasPrefix(params.Prefix, "front_matter.") || strings.HasPrefix(params.Prefix, "body.") {
		// Block path traversal in prefix first
		if containsPathTraversal(params.Prefix) {
			ind.Lg.Warn(
				"path traversal detected in filter prefix",
				logseal.F{"prefix": params.Prefix},
			)
			return false
		}
		matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_.\[\]]*$`, params.Prefix)
		if !matched {
			return false
		}
	} else if params.Prefix != "" {
		// If not a known prefix, it's invalid
		ind.Lg.Warn(
			"invalid filter prefix",
			logseal.F{"prefix": params.Prefix},
		)
		return false
	}

	// Validate Operator - whitelist safe operators
	safeOperators := map[string]bool{
		"":    true,
		"===": true,
		"!==": true,
		"==":  true,
		"!=":  true,
		"==~": true,
		"=~":  true,
		"=*":  true,
		"!=*": true,
	}
	if !safeOperators[params.Operator] {
		return false
	}

	// Validate Suffix entries - ensure no path traversal or injection
	for i, suffix := range params.Suffix {
		// Check individual suffix length
		if !checkInputLength(suffix, MaxFilterSuffixLength, "filter.suffix", ind.Lg) {
			return false
		}

		// Block path traversal patterns
		if containsPathTraversal(suffix) {
			ind.Lg.Warn(
				"path traversal detected in filter suffix",
				logseal.F{"suffix_index": i, "suffix": suffix[:min(len(suffix), 100)]},
			)
			return false
		}
		// For regex operators, validate the regex pattern is safe
		if strings.HasPrefix(params.Operator, "=~") || strings.HasPrefix(params.Operator, "==~") {
			// Block potentially dangerous regex patterns that could cause ReDoS
			// Limit pattern length to prevent memory exhaustion
			const MaxRegexSuffixLength = 200 // Stricter limit for regex patterns
			if len(suffix) > MaxRegexSuffixLength {
				ind.Lg.Warn(
					"regex filter suffix too long",
					logseal.F{"length": len(suffix), "max": MaxRegexSuffixLength},
				)
				return false
			}
			// Try to compile as a sanity check
			if _, err := regexp.Compile(suffix); err != nil {
				ind.Lg.Warn(
					"invalid regex in filter suffix",
					logseal.F{"suffix_index": i, "suffix": suffix, "error": err},
				)
				return false
			}
		}
	}

	// Validate Errors are empty (should be if parsing succeeded)
	if len(params.Errors) > 0 {
		return false
	}

	return true
}
