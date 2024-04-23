package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/triole/logseal"
	"github.com/yuin/goldmark"
	goldmarkmeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v3"
)

func readDataFile(filename string, basePath string, chin chan string, chout chan tJoinerEntry) {
	chin <- filename
	pth := strings.TrimPrefix(
		strings.TrimPrefix(filename, basePath), string(filepath.Separator),
	)
	je := tJoinerEntry{
		Depth:    len(strings.Split(pth, string(filepath.Separator))),
		Path:     pth,
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
			if !CLI.Slim {
				data, err = readJson(by)
			}
		case ".md":
			data, err = readMarkdown(by)
		case ".toml":
			if !CLI.Slim {
				data, err = readToml(by)
			}
		case ".yaml":
			if !CLI.Slim {
				data, err = readYaml(by)
			}
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
	var buf bytes.Buffer
	markdown := goldmark.New(
		goldmark.WithExtensions(
			goldmarkmeta.Meta,
		),
	)
	context := parser.NewContext()
	err = markdown.Convert(by, &buf, parser.WithContext(context))
	if err == nil {
		data = make(map[string]interface{})
		data["front_matter"] = goldmarkmeta.Get(context)
		if !CLI.Slim {
			data["body"] = string(by)
		}
	}
	return
}
