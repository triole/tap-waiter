package main

import (
	"testing"

	yaml "gopkg.in/yaml.v3"
)

type tSpecGetContentValTest struct {
	ContentFile string `yaml:"content_file"`
	Content     tContent
	Key         string   `yaml:"key"`
	Exp         []string `yaml:"exp"`
	Res         string   `yaml:"res"`
	Ep          tEndpoint
}

func readGetContentValSpecs(t *testing.T) (specs []tSpecGetContentValTest) {
	filename := fromTestFolder("specs/mapval/spec.yaml")
	by, _, _ := readFile(filename)
	err := yaml.Unmarshal(by, &specs)
	if err != nil {
		t.Errorf("reading specs file failed: %q", filename)
	}
	return
}

func TestGetContentVal(t *testing.T) {
	ep := newTestEndpoint()
	specs := readGetContentValSpecs(t)
	for _, spec := range specs {
		spec.Content = readFileContent(fromTestFolder(spec.ContentFile), ep)
		validateGetContentVal(spec.Key, spec.Content, spec.Exp, t)
	}
}

func validateGetContentVal(key string, mp tContent, exp []string, t *testing.T) {
	b := false
	res := getContentVal(key, mp)
	if len(exp) == len(res) {
		for i, x := range res {
			if x != exp[i] {
				b = true
			}
		}
	} else {
		b = true
	}
	if b {
		t.Errorf(
			"test fail get content val,\nkey: %+v,\nmap: %+v,\nexp: %v,\nres: %v\n\n",
			key, mp, exp, res,
		)
	}
}
