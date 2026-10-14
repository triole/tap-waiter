package indexer

import (
	"regexp"
	"strings"

	"github.com/triole/logseal"
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

// isValidSortBy validates the SortBy parameter against a whitelist to prevent path traversal.
// Only allows specific safe values: "created", "lastmod", "size", or content field paths.
func (ind *Indexer) isValidSortBy(sortBy string) bool {
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
	// Validate Prefix - must be safe field path
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
	for _, suffix := range params.Suffix {
		// Block path traversal patterns
		if containsPathTraversal(suffix) {
			ind.Lg.Warn(
				"path traversal detected in filter suffix",
				logseal.F{"suffix": suffix},
			)
			return false
		}
		// For regex operators, validate the regex pattern is safe
		if strings.HasPrefix(params.Operator, "=~") || strings.HasPrefix(params.Operator, "==~") {
			// Block potentially dangerous regex patterns that could cause ReDoS
			// Limit pattern length to prevent memory exhaustion
			if len(suffix) > 100 {
				ind.Lg.Warn(
					"regex filter suffix too long",
					logseal.F{"length": len(suffix)},
				)
				return false
			}
			// Try to compile as a sanity check
			if _, err := regexp.Compile(suffix); err != nil {
				ind.Lg.Warn(
					"invalid regex in filter suffix",
					logseal.F{"suffix": suffix, "error": err},
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
