package main

import (
	"tap-waiter/src/conf"
	"tap-waiter/src/indexer"
	"tap-waiter/src/util"

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
	ind.RunServer()
}
