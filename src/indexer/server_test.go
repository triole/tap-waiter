package indexer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestParseFilter(t *testing.T) {
	validateParseFilter(
		"front_matter.title==this+is+a+title",
		FilterParams{
			Prefix: "front_matter.title", Operator: "==", Suffix: []string{"this+is+a+title"},
		}, t,
	)
	validateParseFilter(
		"front_matter.title!=not+the+title",
		FilterParams{
			Prefix: "front_matter.title", Operator: "!=", Suffix: []string{"not+the+title"},
		}, t,
	)
	validateParseFilter(
		"front_matter.tags=*title2",
		FilterParams{
			Prefix: "front_matter.tags", Operator: "=*", Suffix: []string{"title2"},
		}, t,
	)
	validateParseFilter(
		"front_matter.tags!=*title",
		FilterParams{
			Prefix: "front_matter.tags", Operator: "!=*", Suffix: []string{"title"},
		}, t,
	)
	validateParseFilter(
		"front_matter.tags!==tag1,tag2",
		FilterParams{
			Prefix: "front_matter.tags", Operator: "!==", Suffix: []string{"tag1", "tag2"},
		}, t,
	)
	validateParseFilter(
		"front_matter.title=~title[0-9]{1,2}",
		FilterParams{
			Prefix: "front_matter.title", Operator: "=~", Suffix: []string{"title[0-9]{1,2}"},
		}, t,
	)
}

func validateParseFilter(filter string, exp FilterParams, t *testing.T) {
	ind, _, _ := prepareTests("", "", true)
	res := ind.parseFilterString(filter)
	r := true
	if res.Prefix != exp.Prefix || res.Operator != exp.Operator {
		r = false
		for i := 0; i < len(exp.Suffix); i++ {
			if res.Suffix[i] != exp.Suffix[i] {
				r = false
			}
		}
	}
	if !r {
		t.Errorf("parse filter failed, \nexp: %+v\ngot: %+v", exp, res)
	}
}

func TestServeContent(t *testing.T) {
	testSpecs := ut.Find(ut.FromTestFolder("specs/server"), "\\.yaml$")
	for _, specFile := range testSpecs {
		validateServeContent(specFile, t)
	}
}

func validateServeContent(specFile string, t *testing.T) {
	ind, _, _ := prepareTests("", "", true)
	var urls []string
	var exp []string
	var spec map[string][]string

	by, _, err := ut.ReadFile(specFile)
	if err != nil {
		t.Errorf("can not read spec file: %q", specFile)
	} else {
		err = yaml.Unmarshal(by, &spec)
		if err != nil {
			t.Errorf("can not unmarshal spec file: %q, %s", specFile, err)
		}
		urls = spec["urls"]
		exp = spec["expectation"]
	}

	testsrv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ind.ServeContent(w, r)
		}))
	defer testsrv.Close()

	for _, url := range urls {
		c := NewClient(testsrv.URL)
		res, err := http.Get(c.url + url)
		if err != nil {
			t.Errorf("test serve content, request failed: %s, %s", url, err)
		}
		defer res.Body.Close()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("test serve content failed, can not read body: %+v", err)
		} else {
			var ji JoinerIndex
			err = json.Unmarshal([]byte(bodyBytes), &ji)
			if err != nil {
				t.Errorf(
					"test joiner index failed, can not unmarshal server response: %+v", err,
				)
			}
			validateJoinerIndex(ji, url, exp, t)
		}
	}
}

func validateJoinerIndex(ji JoinerIndex, url string, exp []string, t *testing.T) {
	failed := false
	if len(ji) != len(exp) {
		t.Errorf(
			"validate joiner index failed: %q, lengths do not match: exp: %+v, got: %+v",
			url, len(exp), len(ji),
		)
	} else {
		for i := 0; i < len(ji); i++ {
			if ji[i].Path != exp[i] {
				failed = true
				break
			}
		}
	}
	if failed {
		t.Errorf(
			"validate joiner index failed: %q\n"+
				"exp, len: %d\n %+v,\n"+
				"got, len: %d\n%+v\n",
			url, len(exp), exp, len(ji), getJoinerIndexPaths(ji),
		)
	}
}

type Client struct {
	url string
}

func NewClient(url string) Client {
	return Client{url}
}

func getJoinerIndexPaths(ji JoinerIndex) (arr []string) {
	for _, el := range ji {
		arr = append(arr, el.Path)
	}
	return
}

func BenchmarkServer(b *testing.B) {
	ind, _, _ := prepareTests("", "", true)

	pos := ut.Trace()
	testsrv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ind.ServeContent(w, r)
		}))
	defer testsrv.Close()
	for url := range ind.Conf.API {
		c := NewClient(testsrv.URL)
		http.Get(c.url + url)
	}
	fmt.Printf("%s took %s with b.N = %d\n", pos, b.Elapsed(), b.N)
}
