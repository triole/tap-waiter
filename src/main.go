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
		makeJoinerIndex(dataPath, outJSON, CLI.Threads, true)
	}
}

func makeJoinerIndex(dataPath string, outFile string, threads int, showProgressBar bool) {
	start := time.Now()

	var bar *progressbar.ProgressBar
	var joinerIndex tJoinerIndex

	dataFiles := find(dataPath, CLI.Rxfilter)
	ln := len(dataFiles)

	if ln < 1 {
		lg.Warn("no data files found", logseal.F{"path": dataPath})
	} else {
		chin := make(chan string, threads)
		chout := make(chan tJoinerEntry, threads)

		lg.Info("files to index", logseal.F{"no": ln, "threads": threads})

		if showProgressBar {
			bar = progressbar.Default(int64(ln))
		}

		for _, fil := range dataFiles {
			go readDataFile(fil, dataPath, chin, chout)
		}

		c := 0
		for li := range chout {
			if showProgressBar {
				bar.Add(1)
			}
			joinerIndex = append(joinerIndex, li)
			c++
			if c >= ln {
				close(chin)
				close(chout)
				break
			}
		}

		if !CLI.DryRun {
			writeJSON(joinerIndex, outFile)
		} else {
			pprint(joinerIndex)
		}

		lg.Info("done", logseal.F{"duration": time.Since(start)})
	}

}
