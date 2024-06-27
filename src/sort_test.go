package main

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
	filename := fromTestFolder("specs/sort/spec.yaml")
	by, _, _ := readFile(filename)
	err := yaml.Unmarshal(by, &specs)
	if err != nil {
		t.Errorf("reading specs file failed: %q", filename)
	}
	return
}

func TestSort(t *testing.T) {
	specs := readSortTestSpecs(t)
	for _, spec := range specs {
		sortBy := "default"
		asc := true
		params := newTestParams(fromTestFolder(spec.ContentFolder), sortBy, asc)
		params.Endpoint.SortFileName = spec.SortFile
		params.Endpoint.IgnoreList = spec.IgnoreList
		idx := makeJoinerIndex(params)
		idx.applySortFileOrder(params)
		if !orderOK(idx, spec.Expectation, t) {
			t.Errorf(
				"sort failed: %s, asc: %v, \nexp: %v, got: %v",
				sortBy, asc, pprintr(spec.Expectation), getJoinerIndexPaths(idx),
			)
		}
	}
}
