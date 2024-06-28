package indexer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
	"tyson-tap/src/conf"

	yaml "gopkg.in/yaml.v3"
)

var (
	tempFolder           = filepath.Join(os.TempDir(), "tyson-tap-testdata")
	globalDummyTestFiles []string
)

type tSpecIndexTest struct {
	Folder      string   `yaml:"folder"`
	SortBy      string   `yaml:"sort_by"`
	Expectation []string `yaml:"expectation"`
	Ascending   bool
}

func readIndexTestSpecs(filename string, t *testing.T) (specs tSpecIndexTest) {
	ind, _, _ := prepareTests("", "", true)
	by, _, _ := ind.Util.ReadFile(filename)
	err := yaml.Unmarshal(by, &specs)
	if err != nil {
		t.Errorf("reading specs file failed: %q", filename)
	}
	return
}
func TestMakeJoinerIndex(t *testing.T) {
	ind, _, _ := prepareTests("", "", true)
	globalDummyTestFiles = createDummyFiles()
	specsArr := []string{"created", "lastmod"}
	for _, el := range specsArr {
		var spec tSpecIndexTest
		spec.Folder = tempFolder
		spec.SortBy = el
		ascending := []bool{true, false}
		for _, asc := range ascending {
			spec.Expectation = globalDummyTestFiles
			spec.Ascending = asc
			if !spec.Ascending {
				spec.Expectation = reverseArr(spec.Expectation)
			}
			validateMakeJoinerIndex(spec, t)
		}
		sort.Strings(globalDummyTestFiles)
	}

	testSpecs := ind.Util.Find(ind.Util.FromTestFolder("specs/index"), "\\.yaml$")
	ascending := []bool{true, false}
	for _, el := range testSpecs {
		spec := readIndexTestSpecs(el, t)
		spec.Folder = ind.Util.FromTestFolder(spec.Folder)
		for _, asc := range ascending {
			spec.Ascending = asc
			if !spec.Ascending {
				spec.Expectation = reverseArr(spec.Expectation)
			}
			validateMakeJoinerIndex(spec, t)
		}
	}
}

func validateMakeJoinerIndex(spec tSpecIndexTest, t *testing.T) {
	_, ji, params := prepareTests("", "", true)
	if !orderOK(ji, spec.Expectation, t) {
		t.Errorf(
			"sort failed: %s, by: %s, asc: %v,\n  exp: %v\n, got: %v",
			params.Endpoint.Folder,
			params.SortBy,
			params.Ascending,
			spec.Expectation, idx,
		)
	}
}

func orderOK(idx JoinerIndex, exp []string, t *testing.T) bool {
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

func newTestParams(folder, sortBy string, ascending bool) (p Params) {
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

func newTestEndpoint() conf.Endpoint {
	return conf.Endpoint{ReturnValues: conf.ReturnValues{
		Created:                  true,
		LastMod:                  true,
		Content:                  true,
		SplitMarkdownFrontMatter: true,
		Size:                     true,
	}}
}
