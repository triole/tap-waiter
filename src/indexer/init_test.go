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
	ji := ind.MakeJoinerIndex(params)
	return ind, ji, params
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
