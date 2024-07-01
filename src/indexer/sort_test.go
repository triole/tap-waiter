package indexer

import (
	"testing"

	yaml "gopkg.in/yaml.v3"
)

type tSpecSortTest struct {
	ContentFolder string   `yaml:"content_folder"`
	SortFile      string   `yaml:"sort_file"`
	Expectation   []string `yaml:"expectation"`
	IgnoreList    []string `yaml:"ignore_list"`
}

func readSortTestSpecs(t *testing.T) (specs []tSpecSortTest) {
	filename := ut.FromTestFolder("specs/sort/spec.yaml")
	by, _, _ := ut.ReadFile(filename)
	err := yaml.Unmarshal(by, &specs)
	if err != nil {
		t.Errorf("reading specs file failed: %q", filename)
	}
	return
}

func TestSort(t *testing.T) {
	ind, _, _ := prepareTests("", "", true)
	specs := readSortTestSpecs(t)
	for _, spec := range specs {
		sortBy := "default"
		asc := true
		params := newTestParams(ut.FromTestFolder(spec.ContentFolder), sortBy, asc)
		params.Endpoint.SortFileName = spec.SortFile
		params.Endpoint.IgnoreList = spec.IgnoreList
		idx := ind.MakeJoinerIndex(params)
		idx.applySortFileOrder(params)
		if !orderOK(idx, spec.Expectation, t) {
			t.Errorf(
				"sort failed: %s, asc: %v, \n  exp: %v\n, got: %v",
				sortBy, asc, spec.Expectation, getJoinerIndexPaths(idx),
			)
		}
	}
}
