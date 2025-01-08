package indexer

import (
	"bytes"
	"encoding/json"
	"path"
	"path/filepath"
	"strings"
	"tyson-tap/src/conf"

	"github.com/c2h5oh/datasize"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/triole/logseal"
	"github.com/yuin/goldmark"
	goldmarkmeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	yaml "gopkg.in/yaml.v3"
)

func (ind Indexer) readDataFile(filename string, ps conf.Endpoint, chin chan string, chout chan JoinerEntry) {
	chin <- filename
	pth := path.Base(filename)
	if !strings.EqualFold(filename, ps.Source) {
		pth = strings.TrimPrefix(
			strings.TrimPrefix(filename, ps.Source), string(filepath.Separator),
		)
	}
	je := JoinerEntry{
		Path: pth,
	}
	fileSize := ind.Util.GetFileSize(filename)
	if ps.ReturnValues.Size {
		je.Size = fileSize
	}
	if ps.MaxReturnSizeBytes > fileSize {
		if ps.ReturnValues.Content || ps.ReturnValues.SplitMarkdownFrontMatter {
			je.Content = ind.readFileContent(filename, ps)
		}
	} else {
		ind.Lg.Trace(
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
		je.Created = ind.Util.GetFileCreated(filename)
	}
	if ps.ReturnValues.LastMod {
		je.LastMod = ind.Util.GetFileLastMod(filename)
	}
	chout <- je
	<-chin
}

func (ind Indexer) readFileContent(filename string, ps conf.Endpoint) (content FileContent) {
	by, isTextfile, err := ind.Util.ReadFile(filename)
	if isTextfile {
		if err == nil {
			switch filepath.Ext(filename) {
			case ".json":
				content = ind.unmarshalJSON(by)
			case ".toml":
				content = ind.unmarshalTOML(by)
			case ".yaml", ".yml":
				content = ind.unmarshalYAML(by)
			case ".md":
				content = ind.readMarkdown(by, ps.ReturnValues)
			default:
				content = ind.byteToBody(by)
			}
			ind.Lg.IfErrError(
				"error reading file",
				logseal.F{
					"path": filename, "error": err, "is_text_file": isTextfile,
				},
			)
		}
	} else {
		ind.Lg.Debug(
			"no text file, skip reading",
			logseal.F{"path": filename, "is_text_file": isTextfile},
		)
	}
	return
}

func (ind Indexer) byteToBody(by []byte) (content FileContent) {
	content.Body = string(by)
	return
}

func (ind Indexer) unmarshalJSON(by []byte) (content FileContent) {
	var unmarsh interface{}
	content.Error = json.Unmarshal(by, &unmarsh)
	content.Body = unmarsh
	return content
}

func (ind Indexer) unmarshalTOML(by []byte) (content FileContent) {
	var unmarsh interface{}
	content.Error = toml.Unmarshal(by, &unmarsh)
	content.Body = unmarsh
	return content
}

func (ind Indexer) unmarshalYAML(by []byte) (content FileContent) {
	var unmarsh interface{}
	content.Error = yaml.Unmarshal(by, &unmarsh)
	content.Body = unmarsh
	return content
}

func (ind Indexer) readMarkdown(by []byte, rv conf.ReturnValues) (content FileContent) {
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
