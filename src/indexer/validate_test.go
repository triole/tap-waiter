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
