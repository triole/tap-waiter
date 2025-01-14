package indexer

import (
	"bytes"
	"encoding/json"
	"errors"
	"path"
	"path/filepath"
	"strings"
	"tyson-tap/src/conf"

	"github.com/c2h5oh/datasize"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/tidwall/gjson"
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
	if ps.Return.Size {
		je.Size = fileSize
	}
	if ps.MaxReturnSizeBytes > fileSize {
		if ps.Return.Content || ps.Return.SplitMarkdownFrontMatter {
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
	if ps.Return.SplitPath {
		je.SplitPath = strings.Split(pth, string(filepath.Separator))
	}
	if ps.Return.Created {
		je.Created = ind.Util.GetFileCreated(filename)
	}
	if ps.Return.LastMod {
		je.LastMod = ind.Util.GetFileLastMod(filename)
	}
	chout <- je
	<-chin
}

func (ind Indexer) readFileContent(filename string, ps conf.Endpoint) (content FileContent) {
	by, isTextfile, err := ind.Util.ReadFile(filename)
	if isTextfile && err == nil {
		switch filepath.Ext(filename) {
		case ".json", ".toml", ".yaml", ".yml":
			content = ind.unmarshal(by, ps)
		case ".md":
			content = ind.readMarkdown(by, ps.Return)
		default:
			content = ind.byteToBody(by)
		}
		ind.Lg.IfErrError(
			"error reading file",
			logseal.F{
				"path": filename, "error": err, "is_text_file": isTextfile,
			},
		)
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

func (ind Indexer) unmarshal(by []byte, ps conf.Endpoint) (content FileContent) {
	if ind.Util.IsTextData(by) {

		for _, el := range ps.Return.RegexReplace {
			by = []byte(
				ind.Util.RxReplaceAll(
					string(by),
					el[0],
					el[1],
				),
			)
		}

		if content = ind.unmarshalJSON(by); content.Error == nil {
			ind.Lg.Trace("json unmarshalled", logseal.F{"content": string(by)})
		}
		if content.Error != nil {
			if content = ind.unmarshalTOML(by); content.Error == nil {
				ind.Lg.Trace("toml unmarshalled", logseal.F{"content": string(by)})
				// return
			}
		}
		/* NOTE: responses like "404 page not found" are unmarshalled as yaml,
		find out later if this is only an inconsistency or a real problem */
		if content.Error != nil {
			if content = ind.unmarshalYAML(by); content.Error == nil {
				ind.Lg.Trace("yaml unmarshalled", logseal.F{"content": string(by)})
				// return
			}
		}
	}
	if !ind.Util.IsTextData(by) || content.Error != nil {
		content = ind.byteToBody(by)
		content.Error = errors.New("unmarshal failed, kept the plain data")
		ind.Lg.Trace(
			"could not unmarshal",
			logseal.F{"content": string(by), "err": content.Error},
		)
	}
	if ps.Return.JSONPath != "" {
		ind.Lg.Trace(
			"parse unmarshalled data using json path",
			logseal.F{
				"json_path": ps.Return.JSONPath,
			},
		)
		marsh, err := json.Marshal(content.Body)
		ind.Lg.IfErrError(
			"can not prepare to apply json path",
			logseal.F{"error": err},
		)
		if err == nil {
			result := gjson.GetBytes(marsh, ps.Return.JSONPath)
			if len(result.String()) < 1 {
				ind.Lg.Warn(
					"json path result is empty",
					logseal.F{"json_path": ps.Return.JSONPath},
				)
			} else {
				content = ind.unmarshalJSON([]byte(result.String()))
			}
		} else {
			content = FileContent{}
		}
	}

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
