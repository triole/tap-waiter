package indexer

import (
	"tyson-tap/src/conf"
	"tyson-tap/src/util"

	"github.com/triole/logseal"
)

var (
	ut util.Util
)

type Indexer struct {
	TapIndex    TapIndex
	DataSources DataSources
	Conf        conf.Conf
	Util        util.Util
	Lg          logseal.Logseal
}

func Init(conf conf.Conf, util util.Util, lg logseal.Logseal) (idx Indexer) {
	idx = Indexer{
		Conf: conf,
		Util: util,
		Lg:   lg,
	}
	return idx
}
