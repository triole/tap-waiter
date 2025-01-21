package conf

import (
	"tap-waiter/src/util"

	"github.com/triole/logseal"
)

func Init(confFile string, threads int, util util.Util, lg logseal.Logseal) (conf Conf) {
	conf = Conf{
		FileName: confFile,
		Threads:  threads,
		Lg:       lg,
		Util:     util,
	}
	conf.API = make(map[string]Endpoint)
	conf.readConfig()
	return
}
