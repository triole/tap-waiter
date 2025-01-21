package indexer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"tap-waiter/src/conf"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/tidwall/gjson"
	"github.com/triole/logseal"
	"github.com/yuin/goldmark"
	goldmarkmeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	yaml "gopkg.in/yaml.v3"
)

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
		ind.returnJSONPath(content, ps.Return.JSONPath)
	}
	return
}
func (ind Indexer) applyJSONPath(ti TapIndex, jsonPath string) TapIndex {
	if jsonPath != "" {
		for idx, el := range ti {
			ti[idx].Content = ind.returnJSONPath(
				el.Content, jsonPath,
			)
		}
	}
	return ti
}

func (ind Indexer) returnJSONPath(content FileContent, jsonPath string) (r FileContent) {
	ind.Lg.Trace(
		"parse unmarshalled data using json path",
		logseal.F{
			"json_path": jsonPath,
		},
	)
	marsh, err := json.Marshal(content.Body)
	ind.Lg.IfErrError(
		"can not prepare to apply json path",
		logseal.F{"error": err},
	)
	if err == nil {
		result := gjson.GetBytes(marsh, jsonPath)
		if len(result.String()) < 1 {
			ind.Lg.Warn(
				"json path result is empty",
				logseal.F{"json_path": jsonPath},
			)
		} else {
			r = ind.unmarshalJSON([]byte(result.String()))
		}
	} else {
		r = FileContent{}
	}
	return
}

func (ind Indexer) returnRegexMatch(content FileContent, regex []string) (r []string) {
	for _, rx := range regex {
		r = append(
			r,
			ind.Util.RxFindAll(rx, fmt.Sprintf("%s", content))...,
		)
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
