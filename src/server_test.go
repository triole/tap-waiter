package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestParseFilter(t *testing.T) {
	validateParseFilter(
		"front_matter.title==this+is+a+title",
		tIDXParamsFilter{
			Prefix: "front_matter.title", Operator: "==", Suffix: []string{"this+is+a+title"},
		}, t,
	)
	validateParseFilter(
		"front_matter.title!=not+the+title",
		tIDXParamsFilter{
			Prefix: "front_matter.title", Operator: "!=", Suffix: []string{"not+the+title"},
		}, t,
	)
	validateParseFilter(
		"front_matter.tags=*title2",
		tIDXParamsFilter{
			Prefix: "front_matter.tags", Operator: "=*", Suffix: []string{"title2"},
		}, t,
	)
	validateParseFilter(
		"front_matter.tags!=*title",
		tIDXParamsFilter{
			Prefix: "front_matter.tags", Operator: "!=*", Suffix: []string{"title"},
		}, t,
	)
	validateParseFilter(
		"front_matter.tags!==tag1,tag2",
		tIDXParamsFilter{
			Prefix: "front_matter.tags", Operator: "!==", Suffix: []string{"tag1", "tag2"},
		}, t,
	)
	validateParseFilter(
		"front_matter.title=~title[0-9]{1,2}",
		tIDXParamsFilter{
			Prefix: "front_matter.title", Operator: "=~", Suffix: []string{"title[0-9]{1,2}"},
		}, t,
	)
}

func validateParseFilter(filter string, exp tIDXParamsFilter, t *testing.T) {
	res := parseFilterString(filter)
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
	testSpecs := find(fromTestFolder("specs/server"), "\\.yaml$")
	for _, specFile := range testSpecs {
		validateServeContent(specFile, t)
	}
}

func validateServeContent(specFile string, t *testing.T) {
	CLI.LogLevel = "trace"
	CLI.Threads = 16

	var urls []string
	var exp []string
	var spec map[string][]string
	conf = readConfig(fromTestFolder("conf.yaml"))

	by, _, err := readFile(specFile)
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

	svr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			serveContent(w, r)
		}))
	defer svr.Close()

	for _, url := range urls {
		c := NewClient(svr.URL)
		res, err := http.Get(c.url + url)
		if err != nil {
			t.Errorf("test serve content, request failed: %s, %s", url, err)
		}
		defer res.Body.Close()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("test serve content failed, can not read body: %+v", err)
		} else {
			var ji tJoinerIndex
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

func validateJoinerIndex(ji tJoinerIndex, url string, exp []string, t *testing.T) {
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
			url, len(exp), pprintr(exp), len(ji), pprintr(getJoinerIndexPaths(ji)),
		)
	}
}

type Client struct {
	url string
}

func NewClient(url string) Client {
	return Client{url}
}

func getJoinerIndexPaths(ji tJoinerIndex) (arr []string) {
	for _, el := range ji {
		arr = append(arr, el.Path)
	}
	return
}
