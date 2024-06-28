package conf

import (
	"testing"
	"tyson-tap/src/util"

	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
	ut util.Util
)

func init() {
	lg = logseal.Init("debug", "stdout", true, false)
	ut = util.Init(lg)
}
func TestConf(t *testing.T) {
	fn := ut.FromTestFolder("conf.yaml")
	threads := 16
	conf := Init(fn, threads, ut, lg)
	if conf.FileName != fn {
		t.Errorf("%s: %s", ut.Trace(), "filename does not match")
	}
	if conf.Threads != threads {
		t.Errorf("%s: %s", ut.Trace(), "threads do not match")
	}
}
