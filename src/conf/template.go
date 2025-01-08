package conf

import (
	"bytes"
	"os"
	"os/user"
	"path/filepath"
	"text/template"

	"github.com/triole/logseal"
)

func (conf Conf) templateFile(inp []byte) (by []byte, err error) {
	ud := conf.getUserdataMap()
	buf := &bytes.Buffer{}
	templ, err := template.New("conf").Parse(string(inp))
	if err == nil {
		templ.Execute(buf, map[string]interface{}{
			"bindir":  conf.Util.GetBinDir(),
			"confdir": filepath.Dir(conf.Util.AbsPathSlim(conf.FileName)),
			"selfdir": filepath.Dir(conf.Util.AbsPathSlim(conf.FileName)),
			"workdir": conf.pwd(),
			"home":    ud["home"],
			"uid":     ud["uid"],
			"gid":     ud["gid"],
			"user":    ud["username"],
		})
		by = buf.Bytes()
	}
	return
}

func (conf Conf) getUserdataMap() map[string]string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	m := make(map[string]string)
	m["home"] = user.HomeDir + "/"
	m["uid"] = user.Uid
	m["gid"] = user.Gid
	m["username"] = user.Username
	m["name"] = user.Name
	return m
}

func (conf Conf) pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		conf.Lg.Error("can not determine current dir", logseal.F{"error": err})
		os.Exit(1)
	}
	return pwd
}
