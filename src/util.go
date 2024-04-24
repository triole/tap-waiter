package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/triole/logseal"
)

func absPath(str string) string {
	p, err := filepath.Abs(str)
	lg.IfErrFatal("invalid file path", logseal.F{"path": str, "error": err})
	return p
}

func getFileSize(filename string) (siz int64) {
	file, err := os.Open(filename)
	lg.IfErrError(
		"can not open file to get file size",
		logseal.F{"path": filename, "error": err},
	)
	if err == nil {
		defer file.Close()
		stat, err := file.Stat()
		lg.IfErrError(
			"can not stat file to get file size",
			logseal.F{"path": filename, "error": err},
		)
		if err == nil {
			siz = stat.Size()
		}
	}
	return
}

func getFileCreated(filename string) (uts int64) {
	var scs syscall.Stat_t
	if err := syscall.Stat(filename, &scs); err != nil {
		lg.Error("syscall stat failed", logseal.F{"path": filename, "error": err})
	}
	uts = scs.Ctim.Sec
	return
}

func getFileLastMod(filename string) (uts int64) {
	fil, err := os.Stat(filename)
	if err != nil {
		lg.Error("can not stat file", logseal.F{"path": filename, "error": err})
		return
	}
	uts = fil.ModTime().Unix()
	return
}

func pprint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(s))
}
