package main

import (
	"path/filepath"
	"testing"
)

type tSpecFilterTest struct {
	Name string
	Pre  []string
	Suf  []string
	Exp  bool
	Res  bool
}

func readFilterSpecs(filename string, t *testing.T) (r []tSpecFilterTest) {
	specs := readYAMLFile(
		filepath.Join(fromTestFolder("specs/filter"), filename),
	)
	if len(specs) == 0 {
		t.Errorf("reading specs file failed: %q", filename)
	}
	for name, val := range specs {
		spec := val.(map[string]interface{})
		pre := itfArrTostrArr(spec["pre"].([]interface{}))
		suf := itfArrTostrArr(spec["suf"].([]interface{}))
		exp := spec["exp"].(bool)
		r = append(r, tSpecFilterTest{
			Name: name,
			Pre:  pre,
			Suf:  suf,
			Exp:  exp,
		})
	}
	return
}

func printTestFilterResult(spec tSpecFilterTest, t *testing.T) {
	if spec.Exp != spec.Res {
		t.Errorf("error filter test: %v", spec)
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
