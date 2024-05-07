package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/c2h5oh/datasize"
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
	je := tJoinerEntry{
		Path: pth,
	}
	fileSize := getFileSize(filename)
	if ps.ReturnValues.Size {
		je.Size = fileSize
	}
	if ps.MaxReturnSizeBytes > fileSize {
		if ps.ReturnValues.Content || ps.ReturnValues.SplitMarkdownFrontMatter {
			je.Content = readFileContent(filename, ps)
		}
	} else {
		lg.Trace(
			"do not display file content, size limit exceeded",
			logseal.F{
				"path":      filename,
				"file_size": datasize.ByteSize(fileSize).HumanReadable(),
				"max_size":  datasize.ByteSize(ps.MaxReturnSizeBytes).HumanReadable(),
			},
		)
	}
	if ps.ReturnValues.SplitPath {
		je.SplitPath = strings.Split(pth, string(filepath.Separator))
	}
	if ps.ReturnValues.Created {
		je.Created = getFileCreated(filename)
	}
	if ps.ReturnValues.LastMod {
		je.LastMod = getFileLastMod(filename)
	}
	chout <- je
	<-chin
}

func readFileContent(filename string, ps tEndpoint) (content map[string]interface{}) {
	by, isTextfile, err := readFile(filename)
	if isTextfile {
		if err == nil {
			switch filepath.Ext(filename) {
			case ".json":
				content, err = readJson(by)
			case ".toml":
				content, err = readToml(by)
			case ".yaml", ".yml":
				content, err = readYaml(by)
			case ".md":
				content, err = readMarkdown(by, ps.ReturnValues)
			default:
				content = byteToBody(by)
			}
			lg.IfErrError(
				"error reading file",
				logseal.F{
					"path": filename, "error": err, "is_text_file": isTextfile,
				},
			)
		}
	} else {
		lg.Debug(
			"no text file, skip reading",
			logseal.F{"path": filename, "is_text_file": isTextfile},
		)
	}
	return
}

func byteToBody(by []byte) (m map[string]interface{}) {
	m = make(map[string]interface{})
	m["body"] = string(by)
	return
}

func readFile(filename string) (by []byte, isTextfile bool, err error) {
	by, err = os.ReadFile(filename)
	isTextfile = utf8.ValidString(string(by))
	lg.IfErrError(
		"can not read file", logseal.F{"path": filename, "error": err},
	)
	return
}

func readJson(by []byte) (content map[string]interface{}, err error) {
	err = json.Unmarshal(by, &content)
	return content, err
}

func readToml(by []byte) (content map[string]interface{}, err error) {
	err = toml.Unmarshal(by, &content)
	return content, err
}

func readYaml(by []byte) (content map[string]interface{}, err error) {
	err = yaml.Unmarshal(by, &content)
	return content, err
}

func readMarkdown(by []byte, rv tReturnValues) (content map[string]interface{}, err error) {
	content = make(map[string]interface{})
	var buf bytes.Buffer
	markdown := goldmark.New(
		goldmark.WithExtensions(
			goldmarkmeta.Meta,
		),
	)
	context := parser.NewContext()
	err = markdown.Convert(by, &buf, parser.WithContext(context))
	if err == nil {
		if rv.Content {
			content["body"] = string(by)
		}
		if rv.SplitMarkdownFrontMatter {
			content["front_matter"] = goldmarkmeta.Get(context)
		}
	}
	return
}
