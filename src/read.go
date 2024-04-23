package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/triole/logseal"
	"gopkg.in/yaml.v3"
)

func readDataFile(filename string, basePath string, chin chan string, chout chan tJoinerEntry) {
	chin <- filename

	je := tJoinerEntry{
		Path:     strings.TrimPrefix(strings.TrimPrefix(filename, basePath), "/"),
		Ext:      filepath.Ext(filename),
		FileMeta: getFileMeta(filename),
		Data:     readFile(filename),
	}

	chout <- je
	<-chin
}

func readFile(filename string) (data map[string]interface{}) {
	by, err := os.ReadFile(filename)
	lg.IfErrError(
		"can not read file", logseal.F{"path": filename, "error": err},
	)
	if err == nil {
		switch filepath.Ext(filename) {
		case ".json":
			data, err = readJson(by)
		case ".md":
			data, err = readMarkdown(by)
		case ".toml":
			data, err = readToml(by)
		case ".yaml":
			data, err = readYaml(by)
		}
		lg.IfErrError(
			"can not unmarshal data", logseal.F{"path": filename, "error": err},
		)
	}
	return
}

func readJson(by []byte) (data map[string]interface{}, err error) {
	err = json.Unmarshal(by, &data)
	return data, err
}

func readToml(by []byte) (data map[string]interface{}, err error) {
	err = toml.Unmarshal(by, &data)
	return data, err
}

func readYaml(by []byte) (data map[string]interface{}, err error) {
	err = yaml.Unmarshal(by, &data)
	return data, err
}

func readMarkdown(by []byte) (data map[string]interface{}, err error) {
	return
}
