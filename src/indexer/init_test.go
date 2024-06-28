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
	params := newTestParams(folder, sortBy, asc)
	ji := ind.MakeJoinerIndex(params)
	return ind, ji, params
}
