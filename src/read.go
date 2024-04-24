package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func readDataFile(filename string, ps tEndpoint, chin chan string, chout chan tJoinerEntry) {
	chin <- filename
	pth := strings.TrimPrefix(
		strings.TrimPrefix(filename, ps.Folder), string(filepath.Separator),
	)
	lg.Trace("process file", logseal.F{
		"path": filename, "endpoint_config": fmt.Sprintf("%+v", ps),
	})
	je := tJoinerEntry{
		Depth:   len(strings.Split(pth, string(filepath.Separator))) - 1,
		Path:    pth,
		Ext:     filepath.Ext(filename),
		Content: readFileContent(filename, ps),
	}
	fm := readFileMeta(filename, ps)
	if !fm.LastMod.Time.IsZero() {
		je.FileMetadata.LastMod = fm.LastMod
	} else {
		lg.Trace("omit lastmod metadata", logseal.F{"path": filename})
	}
	if !fm.Created.Time.IsZero() {
		je.FileMetadata.Created = fm.Created
	} else {
		lg.Trace("omit created metadata", logseal.F{"path": filename})
	}
	chout <- je
	<-chin
}

func readFileMeta(filename string, ps tEndpoint) (fm tFileMeta) {
	switch filepath.Ext(filename) {
	case ".json":
		if !ps.Readers.Json.OmitMetadata {
			fm = getFileMeta(filename)
		}
	case ".md":
		if !ps.Readers.Markdown.OmitMetadata {
			fm = getFileMeta(filename)
		}
	case ".toml":
		if !ps.Readers.Toml.OmitContent {
			fm = getFileMeta(filename)
		}
	case ".yaml":
		if !ps.Readers.Yaml.OmitContent {
			fm = getFileMeta(filename)
		}
	}
	return
}

func readFileContent(filename string, ps tEndpoint) (data map[string]interface{}) {
	by, err := os.ReadFile(filename)
	lg.IfErrError(
		"can not read file", logseal.F{"path": filename, "error": err},
	)
	if err == nil {
		switch filepath.Ext(filename) {
		case ".json":
			if !ps.Readers.Json.OmitContent {
				data, err = readJson(by)
			}
		case ".md":
			data, err = readMarkdown(by, ps.Readers.Markdown)
		case ".toml":
			if !ps.Readers.Toml.OmitContent {
				data, err = readToml(by)
			}
		case ".yaml":
			if !ps.Readers.Yaml.OmitContent {
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

func readMarkdown(by []byte, omit tReaderMarkdown) (data map[string]interface{}, err error) {
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
		if !omit.OmitFrontMatter {
			data["front_matter"] = goldmarkmeta.Get(context)
		}
		if !omit.OmitBody {
			data["body"] = string(by)
		}
	}
	return
}
