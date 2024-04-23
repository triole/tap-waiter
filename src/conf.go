package main

import (
	"os"
	"path"

	"github.com/triole/logseal"
	"gopkg.in/yaml.v3"
)

type tConf struct {
	Port int                 `yaml:"port"`
	API  map[string]tPathSet `yaml:"api"`
}

type tPathSet struct {
	Folder   string `yaml:"folder"`
	RxFilter string `yaml:"rxfilter"`
	Slim     bool   `yaml:"slim"`
}

func newConf() tConf {
	m := make(map[string]tPathSet)
	return tConf{
		Port: 0,
		API:  m,
	}
}

func readConfig(filename string) (conf tConf) {
	tempconf := newConf()
	by, err := os.ReadFile(filename)
	lg.IfErrFatal(
		"can not read file", logseal.F{"path": filename, "error": err},
	)
	err = yaml.Unmarshal(by, &tempconf)
	lg.IfErrFatal(
		"can not unmarshal config", logseal.F{"path": filename, "error": err},
	)
	conf = newConf()
	conf.Port = tempconf.Port
	for key, val := range tempconf.API {
		key = "/" + path.Clean(key)
		val.Folder = absPath(val.Folder)
		conf.API[key] = val
	}
	return
}
