//go:build !windows && !freebsd && !darwin

package util

import (
	"syscall"

	"github.com/triole/logseal"
)

func (util Util) GetFileCreated(filename string) (uts int64) {
	var scs syscall.Stat_t
	if err := syscall.Stat(filename, &scs); err != nil {
		util.Lg.Error(
			"syscall stat failed", logseal.F{"path": filename, "error": err},
		)
	}
	uts = int64(scs.Ctim.Sec)
	return
}
