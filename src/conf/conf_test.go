package conf

import (
	"testing"
	"tyson-tap/src/util"

	"github.com/triole/logseal"
)

func TestConf(t *testing.T) {
	lg := logseal.Init()
	util := util.Init(lg)
	fn := util.FromTestFolder("conf.yaml")
	threads := 16
	conf := Init(fn, threads, util, lg)
	if conf.FileName != fn {
		t.Errorf("%s: %s", util.Trace(), "filename does not match")
	}
	if conf.Threads != threads {
		t.Errorf("%s: %s", util.Trace(), "threads do not match")
	}
}
