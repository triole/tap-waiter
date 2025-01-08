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
	SpecFile    string
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
	specs.SpecFile = filename
	return
}
func TestMakeJoinerIndex(t *testing.T) {
	var ind Indexer
	var ji JoinerIndex
	var params Params
	globalDummyTestFiles = createDummyFiles()
	specsArr := []string{"created", "lastmod"}
	for _, el := range specsArr {
		var spec tSpecIndexTest
		spec.Folder = tempFolder
		spec.SortBy = el
		ascending := []bool{true, false}
		for _, asc := range ascending {
			ind, ji, params = prepareTests(tempFolder, "", asc)
			spec.Expectation = globalDummyTestFiles
			spec.Ascending = asc
			if !spec.Ascending {
				spec.Expectation = ut.ReverseArr(spec.Expectation)
			}
			validateMakeJoinerIndex(ji, params, spec, t)
		}
		sort.Strings(globalDummyTestFiles)
	}

	testSpecs, _ := ind.Util.Find(ind.Util.FromTestFolder("specs/index"), "\\.yaml$")
	ascending := []bool{true, false}
	for _, el := range testSpecs {
		spec := readIndexTestSpecs(el, t)
		spec.Folder = ind.Util.FromTestFolder(spec.Folder)
		for _, asc := range ascending {
			spec.Ascending = asc
			ind, ji, params = prepareTests(spec.Folder, spec.SortBy, spec.Ascending)
			if !spec.Ascending {
				spec.Expectation = ut.ReverseArr(spec.Expectation)
			}
			validateMakeJoinerIndex(ji, params, spec, t)
		}
	}
}

func validateMakeJoinerIndex(ji JoinerIndex, params Params, specs tSpecIndexTest, t *testing.T) {
	if !orderOK(ji, specs.Expectation, t) {
		t.Errorf(
			"sort failed: %s, by: %s, asc: %v,\nspec file: %s,\n  exp: %v\n, got: %v",
			params.Endpoint.Source,
			params.SortBy,
			params.Ascending,
			specs.SpecFile,
			specs.Expectation, getFileNamesOfJI(ji),
		)
	}
}

func orderOK(ji JoinerIndex, exp []string, t *testing.T) bool {
	if len(ji) != len(exp) {
		t.Errorf(
			"sort failed, lengths differ: %-4d != %-4d\n exp: %+v,\n got: %+v ",
			len(exp), len(ji), exp, getFileNamesOfJI(ji))
	} else {
		for i := 0; i <= len(exp)-1; i++ {
			if ji[i].Path != exp[i] {
				return false
			}
		}
	}
	return true
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

func getFileNamesOfJI(ji JoinerIndex) (arr []string) {
	for _, el := range ji {
		arr = append(arr, el.Path)
	}
	return
}
