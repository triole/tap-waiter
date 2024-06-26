package main

import (
	"path/filepath"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

type tSpecSortTest struct {
	ContentFolder string   `yaml:"content_folder"`
	SortFile      string   `yaml:"sort_file"`
	Expectation   []string `yaml:"expectation"`
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
		sortFile, _ := readSortFile(
			filepath.Join(fromTestFolder(spec.ContentFolder), spec.SortFile),
		)
		sortBy := "default"
		asc := true
		idx := makeJoinerIndex(
			newTestParams(fromTestFolder(spec.ContentFolder), sortBy, asc),
		)
		exclusive := false
		if sortFile.Exclusive {
			exclusive = true
			idx = sortExclusive(idx, sortFile)
		} else {
			idx = sortNonExclusive(idx, sortFile)
		}
		t.Errorf("validate sort, exclusive: %v", exclusive)
		if !orderNotOK(idx, spec.Expectation, t) {
			t.Errorf(
				"sort failed: %s, asc: %v, exclusive: %v,\nexp: %v, got: %v",
				sortBy, asc, exclusive, pprintr(spec.Expectation), getJoinerIndexPaths(idx),
			)
		}
	}
}
