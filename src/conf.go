package main

import (
	"os"
	"path"

	"github.com/triole/logseal"
	"gopkg.in/yaml.v3"
)

type tConf struct {
	Port int                  `yaml:"port"`
	API  map[string]tEndpoint `yaml:"api"`
}

type tEndpoint struct {
	Folder       string        `yaml:"folder"`
	RxFilter     string        `yaml:"rxfilter"`
	ReturnValues tReturnValues `yaml:"return_values"`
}

type tReturnValues struct {
	Metadata                 bool `yaml:"metadata"`
	Content                  bool `yaml:"content"`
	Size                     bool `yaml:"size"`
	FileLastMod              bool `yaml:"file_lastmod"`
	FileCreated              bool `yaml:"file_created"`
	SplitMarkdownFrontMatter bool `yaml:"split_markdown_front_matter"`
}

func newConf() tConf {
	m := make(map[string]tEndpoint)
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
