package main

import (
	"path/filepath"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

type tSpecFilterTest struct {
	Name string
	Pre  []string
	Suf  []string
	Exp  bool
	Res  bool
}

func readFilterSpecs(filename string, t *testing.T) (specs []tSpecFilterTest) {
	by, _, _ := readFile(
		filepath.Join(fromTestFolder("specs/filter"), filename),
	)
	err := yaml.Unmarshal(by, &specs)
	if err != nil {
		t.Errorf("reading specs file failed: %q", filename)
	}
	return
}

func printTestFilterResult(spec tSpecFilterTest, t *testing.T) {
	if spec.Exp != spec.Res {
		t.Errorf("error filter test: %+v", spec)
	}
}

func TestEqualSlices(t *testing.T) {
	specs := readFilterSpecs("slice_equals.yaml", t)
	for _, spec := range specs {
		spec.Res = equalSlices(spec.Pre, spec.Suf)
		printTestFilterResult(spec, t)
	}
}

func TestContainsSlice(t *testing.T) {
	specs := readFilterSpecs("slice_contains.yaml", t)
	for _, spec := range specs {
		spec.Res = containsSlice(spec.Pre, spec.Suf)
		printTestFilterResult(spec, t)
	}
}

func TestNotContainsSlice(t *testing.T) {
	specs := readFilterSpecs("slice_not_contains.yaml", t)
	for _, spec := range specs {
		spec.Res = notContainsSlice(spec.Pre, spec.Suf)
		printTestFilterResult(spec, t)
	}
}

func TestRxMatchSliceCompletely(t *testing.T) {
	specs := readFilterSpecs("slice_rxmatch_all.yaml", t)
	for _, spec := range specs {
		spec.Res = rxMatchSliceAll(spec.Pre, spec.Suf)
		printTestFilterResult(spec, t)
	}
}

func TestRxMatchSliceOnce(t *testing.T) {
	specs := readFilterSpecs("slice_rxmatch_once.yaml", t)
	for _, spec := range specs {
		spec.Res = rxMatchSliceOnce(spec.Pre, spec.Suf)
		printTestFilterResult(spec, t)
	}
}
