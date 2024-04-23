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

func readDataFile(filename string, slim bool, basePath string, chin chan string, chout chan tJoinerEntry) {
	chin <- filename
	pth := strings.TrimPrefix(
		strings.TrimPrefix(filename, basePath), string(filepath.Separator),
	)
	je := tJoinerEntry{
		Depth:    len(strings.Split(pth, string(filepath.Separator))),
		Path:     pth,
		Ext:      filepath.Ext(filename),
		FileMeta: getFileMeta(filename),
		Data:     readFile(filename, slim),
	}

	chout <- je
	<-chin
}

func readFile(filename string, slim bool) (data map[string]interface{}) {
	by, err := os.ReadFile(filename)
	lg.IfErrError(
		"can not read file", logseal.F{"path": filename, "error": err},
	)
	if err == nil {
		switch filepath.Ext(filename) {
		case ".json":
			if slim {
				data, err = readJson(by)
			}
		case ".md":
			data, err = readMarkdown(by, slim)
		case ".toml":
			if slim {
				data, err = readToml(by)
			}
		case ".yaml":
			if slim {
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

func readMarkdown(by []byte, slim bool) (data map[string]interface{}, err error) {
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
		if !slim {
			data["body"] = string(by)
		}
	}
	return
}
