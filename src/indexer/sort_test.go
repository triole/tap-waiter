package indexer

import (
	"strings"
	"testing"

	"tap-waiter/src/util"

	"github.com/triole/logseal"
)

func init() {
	// Initialize global ut for tests that use it directly
	lg := logseal.Init("warn")
	ut = util.Init(lg)
}

func TestSort(t *testing.T) {
	tc := InitTests(false)
	specs := tc.readSpecs("specs/sort/spec.yaml")

	for _, specItf := range specs {
		spec := specItf.(map[string]interface{})
		exp := tc.ind.Util.ListInterfaceToListString(
			spec["exp"].([]interface{}),
		)
		ign := tc.ind.Util.ListInterfaceToListString(
			spec["ignore_list"].([]interface{}),
		)
		tc.params.Ascending = true
		tc.params.Endpoint.Source = spec["content_folder"].(string)
		tc.params.Endpoint.SortFileName = spec["sort_file"].(string)
		tc.params.Endpoint.IgnoreList = ign
		tc.ind.flushCache()
		tc.ind.updateTapIndex(tc.params)
		ti := tc.ind.lookupTapIndexCache(tc.params.Endpoint.ID)
		if !tc.orderOK(ti, exp, spec, t) {
			t.Errorf(
				"sort failed: %s, asc: %v, \n  exp: %v\n, got: %v,\n spec: %v",
				spec["content_folder"].(string),
				tc.params.Ascending,
				exp,
				tc.getTapIndexPaths(ti),
				spec,
			)
		}
	}
}

func TestReadSortFile(t *testing.T) {
	ti := TapIndex{}

	tests := []struct {
		name         string
		filename     string
		sortFileName string
		baseDir      string
		expectError  bool
		errorMsg     string
	}{
		// Path traversal attempts
		{"path traversal ..", "/path/to/file.yaml", "../sort.yaml", "/path/to", true, "path traversal detected"},
		{"backslash traversal", "/path/to/file.yaml", ".\\..\\sort.yaml", "/path/to", true, "path traversal detected"},
		{"slash traversal", "/path/to/file.yaml", "./../../sort.yaml", "/path/to", true, "path traversal detected"},

		// Length DoS protection
		{"too long filename", "/path/to/file.yaml", strings.Repeat("a", 101) + ".yaml", "/path/to", true, "too long"},

		// Special character validation
		{"special char semicolon", "/path/to/file.yaml", "sort;yam.l", "/path/to", true, "special characters"},
		{"valid with underscore", "/path/to/file.yaml", "sort_file.yaml", "/path/to", false, ""},
		{"valid with hyphen", "/path/to/file.yaml", "sort-file.yaml", "/path/to", false, ""},
		{"valid with dot", "/path/to/file.yaml", "my.sort.file.yaml", "/path/to", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ti.readSortFile(tt.filename, tt.sortFileName, tt.baseDir)

			if tt.expectError {
				if result.Error == nil {
					t.Errorf("readSortFile expected error containing %q, got no error", tt.errorMsg)
				} else if !strings.Contains(result.Error.Error(), tt.errorMsg) {
					t.Errorf("readSortFile error mismatch: expected to contain %q, got %q", tt.errorMsg, result.Error.Error())
				}
			} else {
				if result.Error != nil {
					t.Errorf("readSortFile expected success, got error: %v", result.Error)
				}
			}
		})
	}
}
