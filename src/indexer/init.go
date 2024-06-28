package indexer

import (
	"tyson-tap/src/conf"
	"tyson-tap/src/util"

	"github.com/triole/logseal"
)

type Indexer struct {
	JI   JoinerIndex
	Conf conf.Conf
	Util util.Util
	Lg   logseal.Logseal
}

var (
	ut  util.Util
	idx Indexer
)

func Init(conf conf.Conf, util util.Util, lg logseal.Logseal) Indexer {
	ut = util
	idx = Indexer{
		Conf: conf,
		Util: util,
		Lg:   lg,
	}
	return idx
}
