package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/triole/logseal"
)

func absPath(str string) (p string, err error) {
	p, err = filepath.Abs(str)
	lg.IfErrFatal("invalid file path", logseal.F{"path": str, "error": err})
	return p, err
}

func getDepth(pth string) int {
	return len(strings.Split(pth, string(filepath.Separator))) - 1
}

func find(basedir string, rxFilter string) []string {
	inf, err := os.Stat(basedir)
	if err != nil {
		lg.IfErrFatal(
			"unable to access md folder", logseal.F{
				"path": basedir, "error": err,
			},
		)
	}
	if !inf.IsDir() {
		lg.Fatal(
			"not a folder, please provide a directory to look for md files.",
			logseal.F{"path": basedir},
		)
	}

	filelist := []string{}
	rxf, _ := regexp.Compile(rxFilter)

	err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if rxf.MatchString(path) {
			inf, err := os.Stat(path)
			if err == nil {
				if !inf.IsDir() {
					filelist = append(filelist, path)
				}
			} else {
				lg.IfErrInfo("stat file failed", logseal.F{"path": path})
			}
		}
		return nil
	})
	lg.IfErrFatal("find files failed", logseal.F{"path": basedir, "error": err})
	return filelist
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
	s := pprintr(i)
	fmt.Println(string(s))
}

func pprintr(i interface{}) (r string) {
	s, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		lg.Error("prety print failed, can not marshal json", logseal.F{"error": err})
	} else {
		r = string(s)
	}
	return
}

func rxFind(rx string, content string) string {
	temp, _ := regexp.Compile(rx)
	return temp.FindString(content)
}

func rxMatch(rx string, content string) bool {
	temp, _ := regexp.Compile(rx)
	return temp.MatchString(content)
}

func toFloat(inp interface{}) (fl float64) {
	switch val := inp.(type) {
	case float32:
		fl = float64(val)
	case float64:
		fl = val
	case int:
		fl = float64(val)
	case int8:
		fl = float64(val)
	case int16:
		fl = float64(val)
	case int32:
		fl = float64(val)
	case int64:
		fl = float64(val)
	case uint:
		fl = float64(val)
	case uint8:
		fl = float64(val)
	case uint16:
		fl = float64(val)
	case uint32:
		fl = float64(val)
	case uint64:
		fl = float64(val)
	}
	return
}
