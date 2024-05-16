package main

import (
	"testing"

	yaml "gopkg.in/yaml.v3"
)

type tSpecGetMapValTest struct {
	ContentFile string `yaml:"content_file"`
	Content     map[string]interface{}
	Key         string   `yaml:"key"`
	Exp         []string `yaml:"exp"`
	Res         string   `yaml:"res"`
	Ep          tEndpoint
}

func readGetMapValSpecs(t *testing.T) (specs []tSpecGetMapValTest) {
	filename := fromTestFolder("specs/mapval/spec.yaml")
	by, _, _ := readFile(filename)
	err := yaml.Unmarshal(by, &specs)
	if err != nil {
		t.Errorf("reading specs file failed: %q", filename)
	}
	return
}

func TestGetMapVal(t *testing.T) {
	ep := newTestEndpoint()
	specs := readGetMapValSpecs(t)
	for _, spec := range specs {
		spec.Content = readFileContent(
			fromTestFolder(spec.ContentFile), ep,
		)
		validateGetMapVal(spec.Key, spec.Content, spec.Exp, t)
	}
}

func validateGetMapVal(key string, mp map[string]interface{}, exp []string, t *testing.T) {
	b := false
	res := getMapVal(key, mp)
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
			"error get map val,\nkey: %+v,\nmap: %+v,\nexp: %v,\nres: %v\n\n",
			key, mp, exp, res,
		)
	}
}
