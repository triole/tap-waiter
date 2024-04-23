package main

import (
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/triole/logseal"
)

var (
	lg = logseal.Init("debug", nil, true, false)
)

func main() {
	parseArgs()
	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)

	dataPath := absPath(CLI.Path)
	outJSON := absPath(CLI.Output)
	lg.Info("start "+appName, logseal.F{"path": CLI.Path})
	if _, err := os.Stat(outJSON); !os.IsNotExist(err) && !CLI.Force {
		lg.Warn("exit, output json file exists", logseal.F{"path": outJSON})
		lg.Info("either choose a different output target or use -f/--force to overwrite")
		os.Exit(0)
	}

	if CLI.Watch {
		watch(dataPath, outJSON)
	} else {
		makeLunrIndex(dataPath, outJSON, CLI.Threads, true)
	}
}

func makeLunrIndex(dataPath string, outFile string, threads int, showProgressBar bool) {
	start := time.Now()

	var bar *progressbar.ProgressBar
	var lunrIndex tJoinerIndex

	dataFiles := find(dataPath, ".(toml|yaml|json)$")
	ln := len(dataFiles)

	if ln < 1 {
		lg.Warn("no data files found", logseal.F{"path": dataPath})
	} else {
		chin := make(chan string, threads)
		chout := make(chan tJoinerEntry, threads)

		potentialEmptyLine()
		lg.Info("md files to process", logseal.F{"no": ln, "threads": threads})
		potentialEmptyLine()

		if showProgressBar == true {
			bar = progressbar.Default(int64(ln))
		}

		for _, fil := range dataFiles {
			go readDataFile(fil, dataPath, chin, chout)
		}

		c := 0
		for li := range chout {
			if showProgressBar == true {
				bar.Add(1)
			}
			lunrIndex = append(lunrIndex, li)
			c++
			if c >= ln {
				close(chin)
				close(chout)
				break
			}
		}

		potentialEmptyLine()
		writeJSON(lunrIndex, outFile)

		lg.Info("done", logseal.F{"duration": time.Since(start)})
		potentialEmptyLine()
	}

}
