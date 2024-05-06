package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	testFolder = "../testdata"
)

func TestMakeJoinerIndex(t *testing.T) {
	validateMakeJoinerIndex(
		nf("dump/yaml"), "created", true, nf("sort_validate/created_asc.json"), t,
	)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "created", false, nf("sort_validate/created_desc.json"), t,
	)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "lastmod", true, nf("sort_validate/lastmod_asc.json"), t,
	)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "lastmod", false, nf("sort_validate/lastmod_desc.json"), t,
	)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "size", true, nf("sort_validate/size_asc.json"), t,
	)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "size", false, nf("sort_validate/size_desc.json"), t,
	)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "path", true, nf("sort_validate/path_asc.json"), t,
	)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "path", false, nf("sort_validate/path_desc.json"), t,
	)
}

func validateMakeJoinerIndex(folder, sortBy string, ascending bool, val string, t *testing.T) {
	idx := makeJoinerIndex(newTestParams(folder, sortBy, ascending))
	exp := loadJSON(val)
	if !reflect.DeepEqual(idx, exp) {
		order := "asc"
		if ascending == false {
			order = "desc"
		}
		t.Errorf(
			"sort failed: %s %s, exp: %v, out: %v",
			sortBy, order, shortprintJI(idx), shortprintJI(exp),
		)
	}
}

func newTestParams(folder, sortBy string, ascending bool) (p tIDXParams) {
	p.Endpoint.Folder = folder
	p.Endpoint.ReturnValues.Content = true
	p.Endpoint.ReturnValues.Created = true
	p.Endpoint.ReturnValues.LastMod = true
	p.Endpoint.ReturnValues.Metadata = true
	p.Endpoint.ReturnValues.Size = true
	p.Threads = 8
	p.Ascending = ascending
	p.SortBy = sortBy
	return
}

func loadJSON(file string) (ji tJoinerIndex) {
	by, _, err := readFile(file)
	if err == nil {
		err := json.Unmarshal(by, &ji)
		if err == nil {
			return
		}
	}
	return
}

func nf(s string) string {
	return filepath.Join(testFolder, s)
}

func shortprintJI(ji tJoinerIndex) (s string) {
	for _, el := range ji {
		s += fmt.Sprintf("%s ", el.Path)
	}
	return
}
