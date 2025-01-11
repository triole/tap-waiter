package indexer

import (
	"tyson-tap/src/conf"
	"tyson-tap/src/util"

	"github.com/triole/logseal"
)

func prepareTests(folder, sortBy string, asc bool) (Indexer, JoinerIndex, Params) {
	lg := logseal.Init()
	ut := util.Init(lg)
	ind := Init(conf.Init(
		ut.FromTestFolder("conf.yaml"), 16, ut, lg), ut, logseal.Init(),
	)
	if folder == "" {
		folder = ut.FromTestFolder("dump")
	}
	params := newTestParams(folder, sortBy, asc)
	params.Endpoint.SourceType = "folder"
	ji := ind.MakeJoinerIndex(params)
	return ind, ji, params
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

func newTestParams(source, sortBy string, ascending bool) (p Params) {
	p.Endpoint = newTestEndpoint()
	p.Endpoint.Source = source
	p.Threads = 8
	p.Ascending = ascending
	p.SortBy = sortBy
	return
}
