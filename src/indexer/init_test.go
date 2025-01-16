package indexer

import (
	"testing"
	"tyson-tap/src/conf"
	"tyson-tap/src/util"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v3"
)

type testContext struct {
	ind    Indexer
	params Params
	t      *testing.T
}

func InitTests(doIndex bool) (tc testContext) {
	lg := logseal.Init()
	util := util.Init(lg)
	conf := conf.Init(ut.FromTestFolder("conf.yaml"), 16, util, lg)
	tc.ind.Lg = lg
	tc.ind.Util = util
	tc.ind.Conf = conf
	tc.ind = Init(conf, util, lg)
	tc.params.Endpoint = newTestEndpoint()
	tc.params.Endpoint.SourceType = "folder"
	if doIndex {
		tc.ind.UpdateTapIndex(tc.params)
	}
	return
}

func (tc testContext) readSpecs(filename string) (specs []interface{}) {
	var by []byte
	var err error
	absFilename := tc.ind.Util.FromTestFolder(filename)

	by, _, err = tc.ind.Util.ReadFile(absFilename)
	if err != nil {
		tc.t.Errorf("can not read specs file: %q", absFilename)
	}

	by, err = tc.ind.Conf.TemplateFile(by)
	if err != nil {
		tc.t.Errorf("can not expand variables in specs file: %q", absFilename)
	}

	err = yaml.Unmarshal(by, &specs)
	if err != nil {
		tc.t.Errorf("reading specs file failed: %q", filename)
	} else {
		tc.ind.Lg.Info("got specs", logseal.F{"filename": filename, "specs": specs})
	}
	return
}

func newTestEndpoint() conf.Endpoint {
	return conf.Endpoint{Return: conf.ReturnValues{
		Created:                  false,
		LastMod:                  false,
		Content:                  true,
		SplitMarkdownFrontMatter: true,
		Size:                     true,
	}}
}

func (tc testContext) orderOK(ti TapIndex, exp []string, t *testing.T) bool {
	if len(ti) != len(exp) {
		t.Errorf(
			"sort failed, lengths differ: %-4d != %-4d\n exp: %+v,\n got: %+v ",
			len(exp), len(ti), exp, tc.getTapIndexFileNames())
	} else {
		for i := 0; i <= len(exp)-1; i++ {
			if ti[i].Path != exp[i] {
				return false
			}
		}
	}
	return true
}

func (tc testContext) getTapIndexFileNames() (arr []string) {
	for _, el := range tc.ind.TapIndex {
		arr = append(arr, el.Path)
	}
	return
}

func (tc testContext) getTapIndexPaths() (arr []string) {
	for _, el := range tc.ind.TapIndex {
		arr = append(arr, el.Path)
	}
	return
}
