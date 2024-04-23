package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/triole/logseal"
)

func absPath(str string) string {
	p, err := filepath.Abs(str)
	lg.IfErrFatal(
		"invalid file path", logseal.F{
			"path": str, "error": err,
		},
	)
	return p
}

func getFileMeta(filename string) (fm tFileMeta) {
	fil, err := os.Stat(filename)
	if err != nil {
		lg.Error("can not stat file", logseal.F{"path": filename, "error": err})
		return
	}
	fm.LastMod = expandTime(fil.ModTime())
	fm.Created = expandTime(getLastMod(filename))
	return
}

func getLastMod(filename string) (t time.Time) {
	var scs syscall.Stat_t
	if err := syscall.Stat(filename, &scs); err != nil {
		lg.Error("syscall stat failed", logseal.F{"path": filename, "error": err})
	}
	ux := scs.Ctim.Sec
	t = time.Unix(int64(ux), 0)
	return t
}

func expandTime(t time.Time) (r tDateTime) {
	r.Time = t
	r.Unix = t.Unix()
	return r
}

func pprint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(s))
}
