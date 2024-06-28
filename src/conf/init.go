package conf

import (
	"tyson-tap/src/util"

	"github.com/triole/logseal"
)

func Init(confFile string, threads int, util util.Util, lg logseal.Logseal) (conf Conf) {
	conf = Conf{
		FileName: confFile,
		Threads:  threads,
		Lg:       lg,
		Util:     util,
	}
	content := conf.readConfig()
	conf.Port = content.Port
	conf.API = content.API
	return
}
