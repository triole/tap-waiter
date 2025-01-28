package indexer

import (
	"tap-waiter/src/conf"
	"tap-waiter/src/util"
	"testing"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v3"
)

type testContext struct {
	ind    Indexer
	params Params
	t      *testing.T
}

func InitTests(doIndex bool) (tc testContext) {
	lg := logseal.Init("warn")
	util := util.Init(lg)
	conf := conf.Init(ut.FromTestFolder("conf.yaml"), 16, util, lg)
	tc.ind.Lg = lg
	tc.ind.Util = util
	tc.ind.Conf = conf
	tc.ind = Init(conf, util, lg)
	tc.params.Endpoint = newTestEndpoint()
	tc.params.Endpoint.SourceType = "folder"
	if doIndex {
		tc.ind.updateTapIndex(tc.params)
	}
	return
}

func (tc testContext) readSpecs(filename string) (specs []interface{}) {
	var by []byte
	var err error
	absFilename := filename
	if !tc.ind.Util.IsAbs(absFilename) {
		absFilename = tc.ind.Util.FromTestFolder(filename)
	}
	tc.ind.Lg.Info("read specs file", logseal.F{"filename": absFilename})
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
		tc.ind.Lg.Fatal(
			"reading specs file failed",
			logseal.F{"path": absFilename, "error": err},
		)
		// tc.t.Errorf("reading specs file failed: %q", absFilename)
	} else {
		tc.ind.Lg.Info("got specs", logseal.F{"filename": absFilename, "specs": specs})
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
			len(exp), len(ti), exp, tc.getTapIndexFileNames(ti))
	} else {
		for i := 0; i <= len(exp)-1; i++ {
			if ti[i].Path != exp[i] {
				return false
			}
		}
	}
	return true
}

func (tc testContext) getTapIndexFileNames(ti TapIndex) (arr []string) {
	for _, el := range ti {
		arr = append(arr, el.Path)
	}
	return
}

func (tc testContext) getTapIndexPaths(ti TapIndex) (arr []string) {
	for _, el := range ti {
		arr = append(arr, el.Path)
	}
	return
}
