package main

import (
	"sort"

	"github.com/triole/logseal"
)

var (
	lg   = logseal.Init("debug", nil, true, false)
	conf = tConf{}
)

func main() {
	parseArgs()
	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	conf = readConfig(CLI.Conf)
	runServer(conf)
}

func makeJoinerIndex(ps tPathSet, threads int) (joinerIndex tJoinerIndex) {
	dataFiles := find(ps.Folder, ps.RxFilter)
	ln := len(dataFiles)

	if ln < 1 {
		lg.Warn("no data files found", logseal.F{"path": ps.Folder})
	} else {
		chin := make(chan string, threads)
		chout := make(chan tJoinerEntry, threads)

		lg.Debug("files to index", logseal.F{"no": ln, "threads": threads})

		for _, fil := range dataFiles {
			go readDataFile(fil, ps.Slim, ps.Folder, chin, chout)
		}

		c := 0
		for li := range chout {
			joinerIndex = append(joinerIndex, li)
			c++
			if c >= ln {
				close(chin)
				close(chout)
				break
			}
		}

		sort.Sort(tJoinerIndex(joinerIndex))
		lg.Debug(
			"index created",
			logseal.F{"path": ps.Folder},
		)
	}
	return
}
