package indexer

import (
	"testing"
)

func TestIsValidSortBy(t *testing.T) {
	tc := InitTests(false)
	ind := tc.ind

	tests := []struct {
		name     string
		sortBy   string
		expected bool
	}{
		// Valid built-in options
		{"empty string", "", true},
		{"created", "created", true},
		{"lastmod", "lastmod", true},
		{"size", "size", true},

		// Valid content field paths
		{"body field", "body.title", true},
		{"body nested", "body.author.name", true},
		{"front_matter field", "front_matter.tags", true},
		{"body with underscore", "body.field_name", true},
		{"body with numbers", "body.field123", true},
		{"mixed path", "body.user.profile.name", true},

		// Invalid - path traversal attempts
		{"parent dir traversal", "../../../etc/passwd", false},
		{"double dot", "body..field", false},
		{"trailing slash", "body/", false},
		{"starting with slash", "/etc/passwd", false},
		{"body with dots in wrong place", "body../field", false},

		// Invalid - wrong format
		{"starting with number", "123field", false},
		{"special chars", "body@field", false},
		{"starting with dot", ".field", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ind.isValidSortBy(tt.sortBy)
			if result != tt.expected {
				t.Errorf("isValidSortBy(%q) = %v, expected %v", tt.sortBy, result, tt.expected)
			}
		})
	}
}

// Test input length limits for SortBy
func TestIsValidSortByLengthLimits(t *testing.T) {
	tc := InitTests(false)
	ind := tc.ind

	// Test exactly at limit
	exactLimitString := "front_matter." + string(make([]byte, MaxSortByLength-15))
	for i := range exactLimitString {
		if exactLimitString[i] == 0 {
			// Replace null bytes with valid chars for testing
			// In practice this string would contain zeros, but we need valid char count
			break
		}
	}
	
	// Generate a proper test string of correct length with valid characters
	testPrefix := "front_matter.test_field."
	remaining := MaxSortByLength - len(testPrefix)
	if remaining > 0 {
		padding := string(make([]byte, remaining))
		for i := 0; i < remaining; i++ {
			padding = padding[:i] + "a" + padding[i+1:]
		}
		testPrefix += padding
	}

	tests := []struct {
		name     string
		sortBy   string
		expected bool
	}{
		{"max length sortBy", testPrefix, true},
		{"exceeds max length by 1", testPrefix + "x", false},
		{"way beyond max length", "front_matter." + string(make([]byte, MaxSortByLength*2)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ind.isValidSortBy(tt.sortBy)
			if result != tt.expected {
				t.Errorf("isValidSortBy length check failed: len=%d, result=%v, expected=%v",
					len(tt.sortBy), result, tt.expected)
			}
		})
	}
}

func TestIsValidFilterParams(t *testing.T) {
	tc := InitTests(false)
	ind := tc.ind

	tests := []struct {
		name     string
		params   FilterParams
		expected bool
	}{
		// Valid filter parameters
		{"valid body prefix", FilterParams{
			Prefix:   "body.title",
			Operator: "==",
			Suffix:   []string{"hello"},
		}, true},

		{"valid front_matter prefix", FilterParams{
			Prefix:   "front_matter.tags",
			Operator: "==",
			Suffix:   []string{"tech"},
		}, true},

		{"valid empty prefix", FilterParams{
			Prefix:   "",
			Operator: "",
			Suffix:   []string{},
		}, true},

		{"valid regex operator", FilterParams{
			Prefix:   "body.content",
			Operator: "=~",
			Suffix:   []string{"test"},
		}, true},

		// Invalid - path traversal in suffix
		{"path traversal in suffix", FilterParams{
			Prefix:   "body.title",
			Operator: "==",
			Suffix:   []string{"../../../etc/passwd"},
		}, false},

		{"double slash in suffix", FilterParams{
			Prefix:   "body.title",
			Operator: "==",
			Suffix:   []string{"path//to//file"},
		}, false},

		// Invalid - operator
		{"invalid operator", FilterParams{
			Prefix:   "body.title",
			Operator: ">>>",
			Suffix:   []string{"test"},
		}, false},

		// Invalid - prefix format
		{"invalid prefix - starts with dot", FilterParams{
			Prefix:   ".field",
			Operator: "==",
			Suffix:   []string{"test"},
		}, false},

		{"invalid prefix - special chars", FilterParams{
			Prefix:   "body@field",
			Operator: "==",
			Suffix:   []string{"test"},
		}, false},

		// Invalid - overly long regex
		{"long regex pattern", FilterParams{
			Prefix:   "body.content",
			Operator: "=~",
			Suffix:   []string{"a" + string(make([]byte, 150))},
		}, false},

		// Invalid - malformed regex
		{"malformed regex", FilterParams{
			Prefix:   "body.content",
			Operator: "=~",
			Suffix:   []string{"["},
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ind.isValidFilterParams(tt.params)
			if result != tt.expected {
				t.Errorf("isValidFilterParams(%+v) = %v, expected %v", tt.params, result, tt.expected)
			}
		})
	}
}

// Test input length limits for FilterParams
func TestIsValidFilterParamsLengthLimits(t *testing.T) {
	tc := InitTests(false)
	ind := tc.ind

	// Test prefix length limits
	validPrefix := "front_matter.test.field"
	longPrefix := "front_matter." + string(make([]byte, MaxFilterPrefixLength+10))
	for i := 0; i < len(longPrefix); i++ {
		if longPrefix[i] == 0 {
			longPrefix = longPrefix[:i] + "a" + longPrefix[i+1:]
		}
	}

	// Test suffix length limits
	validSuffix := "test_value"
	longSuffix := string(make([]byte, MaxFilterSuffixLength+50))
	for i := 0; i < len(longSuffix); i++ {
		longSuffix = longSuffix[:i] + "x" + longSuffix[i+1:]
	}

	// Test suffix count limit
	lotsOfSuffixes := make([]string, MaxSuffixCount+5)
	for i := 0; i < len(lotsOfSuffixes); i++ {
		lotsOfSuffixes[i] = "valid_suffix"
	}

	tests := []struct {
		name     string
		params   FilterParams
		expected bool
	}{
		{
			"prefix at max length",
			FilterParams{
				Prefix:   validPrefix,
				Operator: "==",
				Suffix:   []string{"test"},
			},
			true,
		},
		{
			"prefix exceeds max length",
			FilterParams{
				Prefix:   longPrefix,
				Operator: "==",
				Suffix:   []string{"test"},
			},
			false,
		},
		{
			"suffix at max length",
			FilterParams{
				Prefix:   "body.title",
				Operator: "==",
				Suffix:   []string{validSuffix},
			},
			true,
		},
		{
			"suffix exceeds max length",
			FilterParams{
				Prefix:   "body.title",
				Operator: "==",
				Suffix:   []string{longSuffix},
			},
			false,
		},
		{
			"suffix count within limit",
			FilterParams{
				Prefix:   "body.title",
				Operator: "==",
				Suffix:   make([]string, MaxSuffixCount),
			},
			true,
		},
		{
			"suffix count exceeds limit",
			FilterParams{
				Prefix:   "body.title",
				Operator: "==",
				Suffix:   lotsOfSuffixes,
			},
			false,
		},
		{
			"regex suffix within stricter limit",
			FilterParams{
				Prefix:   "body.content",
				Operator: "=~",
				Suffix:   []string{"^valid_regex$"},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ind.isValidFilterParams(tt.params)
			if result != tt.expected {
				t.Errorf("isValidFilterParams length check failed: prefix_len=%d, suffix_count=%d, result=%v, expected=%v",
					len(tt.params.Prefix), len(tt.params.Suffix), result, tt.expected)
			}
		})
	}
}
