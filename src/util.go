package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/triole/logseal"
)

func absPath(str string) string {
	p, err := filepath.Abs(str)
	lg.IfErrFatal("invalid file path", logseal.F{"path": str, "error": err})
	return p
}

func getFileSize(filename string) (siz uint64) {
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
			siz = uint64(stat.Size())
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

func rxFind(rx string, content string) string {
	temp, _ := regexp.Compile(rx)
	return temp.FindString(content)
}

func rxMatch(rx string, content string) bool {
	temp, _ := regexp.Compile(rx)
	return temp.MatchString(content)
}
