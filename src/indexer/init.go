package indexer

import (
	"time"
	"tyson-tap/src/conf"
	"tyson-tap/src/util"

	cache "github.com/patrickmn/go-cache"
	"github.com/triole/logseal"
)

var (
	ut util.Util
)

type Indexer struct {
	Cache       *cache.Cache
	DataSources DataSources
	Conf        conf.Conf
	Util        util.Util
	Lg          logseal.Logseal
}

func Init(conf conf.Conf, util util.Util, lg logseal.Logseal) (idx Indexer) {
	idx = Indexer{
		Cache: cache.New(5*time.Minute, 10*time.Minute),
		Conf:  conf,
		Util:  util,
		Lg:    lg,
	}
	return idx
}
