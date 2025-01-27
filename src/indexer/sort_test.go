package indexer

import (
	"testing"
)

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
		ti := tc.ind.getTapIndexCache(tc.params.Endpoint.ID)
		if !tc.orderOK(ti, exp, t) {
			t.Errorf(
				"sort failed: %s, asc: %v, \n  exp: %v\n, got: %v",
				spec["content_folder"].(string),
				tc.params.Ascending,
				exp,
				tc.getTapIndexPaths(ti),
			)
		}
	}
}
