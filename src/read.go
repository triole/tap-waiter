package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/c2h5oh/datasize"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/triole/logseal"
	"github.com/yuin/goldmark"
	goldmarkmeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	yaml "gopkg.in/yaml.v3"
)

type tContent struct {
	Body        interface{} `json:"body,omitempty"`
	FrontMatter interface{} `json:"front_matter,omitempty"`
	Error       error       `json:"-"`
}

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

func readFileContent(filename string, ps tEndpoint) (content tContent) {
	by, isTextfile, err := readFile(filename)
	if isTextfile {
		if err == nil {
			switch filepath.Ext(filename) {
			case ".json":
				content = unmarshalJSON(by)
			case ".toml":
				content = unmarshalTOML(by)
			case ".yaml", ".yml":
				content = unmarshalYAML(by)
			case ".md":
				content = readMarkdown(by, ps.ReturnValues)
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

func byteToBody(by []byte) (content tContent) {
	content.Body = string(by)
	return
}

func readFile(filename string) (by []byte, isTextfile bool, err error) {
	fn, err := absPath(filename)
	if err == nil {
		by, err = os.ReadFile(fn)
		isTextfile = utf8.ValidString(string(by))
		lg.IfErrError(
			"can not read file", logseal.F{"path": filename, "error": err},
		)
	}
	return
}

func unmarshalJSON(by []byte) (content tContent) {
	var unmarsh interface{}
	content.Error = json.Unmarshal(by, &unmarsh)
	content.Body = unmarsh
	return content
}

func unmarshalTOML(by []byte) (content tContent) {
	var unmarsh interface{}
	content.Error = toml.Unmarshal(by, &unmarsh)
	content.Body = unmarsh
	return content
}

func unmarshalYAML(by []byte) (content tContent) {
	var unmarsh interface{}
	content.Error = yaml.Unmarshal(by, &unmarsh)
	content.Body = unmarsh
	return content
}

func readMarkdown(by []byte, rv tReturnValues) (content tContent) {
	var buf bytes.Buffer
	markdown := goldmark.New(goldmark.WithExtensions(goldmarkmeta.Meta))
	context := parser.NewContext()
	content.Error = markdown.Convert(by, &buf, parser.WithContext(context))
	if content.Error == nil {
		if rv.Content {
			content.Body = string(by)
		}
		if rv.SplitMarkdownFrontMatter {
			content.FrontMatter = goldmarkmeta.Get(context)
		}
	}
	return
}
