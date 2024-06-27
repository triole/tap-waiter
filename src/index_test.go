package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	yaml "gopkg.in/yaml.v3"
)

var (
	tempFolder     = filepath.Join(os.TempDir(), "tyson_tap_testdata")
	dummyTestFiles []string
)

func init() {
	dummyTestFiles = createDummyFiles()
}

type tSpecIndexTest struct {
	Folder      string   `yaml:"folder"`
	SortBy      string   `yaml:"sort_by"`
	Expectation []string `yaml:"expectation"`
	Ascending   bool
}

func readIndexTestSpecs(filename string, t *testing.T) (specs tSpecIndexTest) {
	by, _, _ := readFile(filename)
	err := yaml.Unmarshal(by, &specs)
	if err != nil {
		t.Errorf("reading specs file failed: %q", filename)
	}
	return
}
func TestMakeJoinerIndex(t *testing.T) {
	// validateMakeJoinerIndex(tempFolder, "created", dummyTestFiles, t)
	// validateMakeJoinerIndex(tempFolder, "created", dummyTestFiles, t)
	// sort.Strings(dummyTestFiles)
	// validateMakeJoinerIndex(tempFolder, "lastmod", dummyTestFiles, t)
	// validateMakeJoinerIndex(tempFolder, "lastmod", dummyTestFiles, t)

	testSpecs := find(fromTestFolder("specs/index"), "\\.yaml$")
	ascending := []bool{true, false}
	for _, el := range testSpecs {
		spec := readIndexTestSpecs(el, t)
		spec.Folder = fromTestFolder(spec.Folder)
		for _, asc := range ascending {
			spec.Ascending = asc
			if !asc {
				spec.Expectation = reverseArr(spec.Expectation)
			}
			validateMakeJoinerIndex(spec, t)
		}
	}
}

func validateMakeJoinerIndex(spec tSpecIndexTest, t *testing.T) {
	params := newTestParams(spec.Folder, spec.SortBy, spec.Ascending)
	idx := makeJoinerIndex(params)
	if !orderOK(idx, spec.Expectation, t) {
		t.Errorf(
			"sort failed: %s, by: %s, asc: %v,\n  exp: %v\n, got: %v",
			params.Endpoint.Folder,
			params.SortBy,
			params.Ascending,
			spec.Expectation, getJoinerIndexPaths(idx),
		)
	}
}

func orderOK(idx tJoinerIndex, exp []string, t *testing.T) bool {
	if len(idx) != len(exp) {
		t.Errorf("sort failed, lengths differ: %-4d != %-4d", len(idx), len(exp))
	}
	for i := 0; i <= len(exp)-1; i++ {
		if idx[i].Path != exp[i] {
			return false
		}
	}
	return true
}

func reverseArr(arr []string) []string {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
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

func createDummyFiles() (arr []string) {
	os.MkdirAll(tempFolder, os.ModePerm)
	for i := 1; i <= 3; i++ {
		name := filepath.Join(tempFolder, fmt.Sprintf("%03d", i)+".tmp")
		_, err := os.Stat(name)
		if err != nil {
			f, err := os.Create(name)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			time.Sleep(time.Duration(2) * time.Second)
		}
		arr = append(arr, filepath.Base(name))
	}
	return
}

func newTestEndpoint() tEndpoint {
	return tEndpoint{ReturnValues: tReturnValues{
		Created:                  true,
		LastMod:                  true,
		Content:                  true,
		SplitMarkdownFrontMatter: true,
		Size:                     true,
	}}
}
