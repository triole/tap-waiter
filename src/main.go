package main

import (
	"fmt"
	"os"

	"github.com/triole/logseal"
)

var (
	lg   = logseal.Init("info", nil, true, false)
	conf = tConf{}
)

func main() {
	parseArgs()
	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	conf = readConfig(CLI.Conf)
	lg.Info(
		"run "+appName, logseal.F{
			"config": CLI.Conf, "log_level": CLI.LogLevel,
		},
	)
	lg.Debug("full configuration layout", logseal.F{"config": fmt.Sprintf("%+v", conf)})
	if CLI.ValidateConf {
		pprint(conf)
		os.Exit(0)
	}
	runServer(conf)
}
