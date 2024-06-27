package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

var (
	tempFolder     = filepath.Join(os.TempDir(), "tyson_tap_testdata")
	dummyTestFiles []string
)

func init() {
	dummyTestFiles = createDummyFiles()
}

func TestMakeJoinerIndex(t *testing.T) {
	validateMakeJoinerIndex(tempFolder, "created", dummyTestFiles, t)
	validateMakeJoinerIndex(tempFolder, "created", dummyTestFiles, t)
	sort.Strings(dummyTestFiles)
	validateMakeJoinerIndex(tempFolder, "lastmod", dummyTestFiles, t)
	validateMakeJoinerIndex(tempFolder, "lastmod", dummyTestFiles, t)

	testSpecs := find(fromTestFolder("specs/index"), "\\.yaml$")
	for _, el := range testSpecs {
		specs := readYAMLFile(el)
		folder := fromTestFolder(specs["folder"].(string))
		sortby := specs["sortby"].(string)
		expectation := itfArrTostrArr(specs["expectation"].([]interface{}))
		validateMakeJoinerIndex(folder, sortby, expectation, t)
		validateMakeJoinerIndex(folder, sortby, expectation, t)
	}
}

func validateMakeJoinerIndex(folder, sortBy string, exp []string, t *testing.T) {
	ascending := []bool{true, false}
	for _, asc := range ascending {
		exp = reverseArr(exp)
		idx := makeJoinerIndex(newTestParams(folder, sortBy, asc))
		if !orderOK(idx, exp, t) {
			t.Errorf(
				"sort failed: %s, asc: %v,\nexp: %v, got: %v",
				sortBy, asc, pprintr(exp), getJoinerIndexPaths(idx),
			)
		}
	}
}

func orderOK(idx tJoinerIndex, exp []string, t *testing.T) bool {
	if len(idx) != len(exp) {
		t.Errorf("sort failed, lengths differ: %-4d != %-4d", len(idx), len(exp))
	}
	for i := 0; i <= len(exp)-1; i++ {
		fmt.Printf("%s === %+v\n", idx[i].Path, exp[i])
		if idx[i].Path != exp[i] {
			return false
		}
	}
	// for i, el := range idx {
	// 	if el.Path != exp[i] {
	// 		return true
	// 	}
	// }
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
