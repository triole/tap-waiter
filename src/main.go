package main

import (
	"fmt"
	"tyson-tap/src/conf"
	"tyson-tap/src/indexer"
	"tyson-tap/src/util"

	"github.com/triole/logseal"
)

func main() {
	parseArgs()
	lg := logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	util := util.Init(lg)
	conf := conf.Init(CLI.Conf, CLI.Threads, util, lg)
	ind := indexer.Init(conf, util, lg)

	lg.Info(
		"run "+appName, logseal.F{
			"config": CLI.Conf, "log_level": CLI.LogLevel,
		},
	)
	lg.Debug("full configuration layout", logseal.F{"config": fmt.Sprintf("%+v", conf)})
	// if CLI.ValidateConf {
	// 	pprint(conf)
	// 	os.Exit(0)
	// }
	ind.RunServer()
}
