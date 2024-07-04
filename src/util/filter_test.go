package util

import (
	"path/filepath"
	"testing"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v3"
)

var (
	lg logseal.Logseal
)

func init() {
	lg = logseal.Init("debug", "stdout", true, false)
}

type tSpecFilterTest struct {
	Name string
	Pre  []string
	Suf  []string
	Exp  bool
	Res  bool
}

func readFilterSpecs(filename string, t *testing.T) (specs []tSpecFilterTest) {
	ut := Init(lg)
	by, _, _ := ut.ReadFile(
		filepath.Join(ut.FromTestFolder("specs/filter"), filename),
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
	ut := Init(lg)
	specs := readFilterSpecs("slice_equals.yaml", t)
	for _, spec := range specs {
		spec.Res = ut.SlicesEqual(spec.Pre, spec.Suf)
		printTestFilterResult(spec, t)
	}
}

func TestContainsSlice(t *testing.T) {
	ut := Init(lg)
	specs := readFilterSpecs("slice_contains.yaml", t)
	for _, spec := range specs {
		spec.Res = ut.SliceContainsSlice(spec.Pre, spec.Suf)
		printTestFilterResult(spec, t)
	}
}

func TestNotContainsSlice(t *testing.T) {
	ut := Init(lg)
	specs := readFilterSpecs("slice_not_contains.yaml", t)
	for _, spec := range specs {
		spec.Res = ut.SliceNotContainsSlice(spec.Pre, spec.Suf)
		printTestFilterResult(spec, t)
	}
}

func TestRxMatchSliceCompletely(t *testing.T) {
	ut := Init(lg)
	specs := readFilterSpecs("slice_rxmatch_all.yaml", t)
	for _, spec := range specs {
		spec.Res = ut.RxSliceMatchesSliceFully(spec.Pre, spec.Suf)
		printTestFilterResult(spec, t)
	}
}

func TestRxMatchSliceOnce(t *testing.T) {
	ut := Init(lg)
	specs := readFilterSpecs("slice_rxmatch_once.yaml", t)
	for _, spec := range specs {
		spec.Res = ut.RxSliceContainsSliceFully(spec.Pre, spec.Suf)
		printTestFilterResult(spec, t)
	}
}
