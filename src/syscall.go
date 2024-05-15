//go:build !windows && !freebsd && !darwin

package main

import (
	"syscall"

	"github.com/triole/logseal"
)

func getFileCreated(filename string) (uts int64) {
	var scs syscall.Stat_t
	if err := syscall.Stat(filename, &scs); err != nil {
		lg.Error(
			"syscall stat failed", logseal.F{"path": filename, "error": err},
		)
	}
	uts = int64(scs.Ctim.Sec)
	return
}
