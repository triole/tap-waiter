package indexer

import (
	"testing"
	"tyson-tap/src/conf"

	yaml "gopkg.in/yaml.v3"
)

type tSpecGetContentValTest struct {
	ContentFile string `yaml:"content_file"`
	Content     FileContent
	Key         string   `yaml:"key"`
	Exp         []string `yaml:"exp"`
	Res         string   `yaml:"res"`
	Ep          conf.Endpoint
}

func readGetContentValSpecs(t *testing.T) (specs []tSpecGetContentValTest) {
	ind, _, _ := prepareTests("", "", true)
	filename := ind.Util.FromTestFolder("specs/mapval/spec.yaml")
	by, _, _ := ind.Util.ReadFile(filename)
	err := yaml.Unmarshal(by, &specs)
	if err != nil {
		t.Errorf("reading specs file failed: %q", filename)
	}
	return
}

func TestGetContentVal(t *testing.T) {
	ind, _, _ := prepareTests("", "", true)
	ep := newTestEndpoint()
	specs := readGetContentValSpecs(t)
	for _, spec := range specs {
		spec.Content = ind.readFileContent(ind.Util.FromTestFolder(spec.ContentFile), ep)
		validateGetContentVal(spec.Key, spec.Content, spec.Exp, t)
	}
}

func validateGetContentVal(key string, mp FileContent, exp []string, t *testing.T) {
	_, ji, _ := prepareTests("", "", true)
	b := false
	res := ji.getContentVal(key, mp)
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
