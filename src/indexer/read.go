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
	// Validate path is safe to prevent path traversal attacks
	// For endpoints with a known source directory, validate against it
	if ps.Source != "" {
		baseDir, err := ind.Util.AbsPath(ps.Source)
		if err == nil {
			if !ind.Util.IsPathSafe(filename, baseDir) {
				ind.Lg.Warn("path traversal attempt blocked during file read", logseal.F{"path": filename, "base_dir": baseDir, "endpoint": ps.ID})
				return content // Return empty content
			}
		}
	}
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
		content = ind.extractWithJSONPath(content, ps.Return.JSONPath)
	}
	return
}
func (ind Indexer) applyJSONPath(ti TapIndex, jsonPath string) TapIndex {
	if jsonPath != "" {
		for idx, el := range ti {
			if el.Content.Error == nil {
				ti[idx].Content = ind.extractWithJSONPath(
					el.Content, jsonPath,
				)
			}
		}
	}
	return ti
}

// extractWithJSONPath applies a JSONPath expression to extract nested values from structured content.
// Returns the extracted data or an empty FileContent if extraction fails.
func (ind Indexer) extractWithJSONPath(content FileContent, jsonPath string) (r FileContent) {
	ind.Lg.Trace(
		"parse unmarshalled data using json path",
		logseal.F{
			"json_path": jsonPath,
		},
	)
	// Marshal content.Body once to get valid JSON bytes
	marsh, err := json.Marshal(content.Body)
	if err != nil {
		ind.Lg.Error(
			"can not marshal content for json path extraction",
			logseal.F{"error": err},
		)
		return FileContent{}
	}

	result := gjson.GetBytes(marsh, jsonPath)
	if !result.Exists() {
		ind.Lg.Warn(
			"json path result not found",
			logseal.F{"json_path": jsonPath},
		)
		return FileContent{}
	}

	r.Body = result.Value()
	return
}

func (ind Indexer) extractWithRegexFromIndex(ti TapIndex, regex []string) (arr []string) {
	for _, el := range ti {
		matches := ind.extractWithRegex(el.Content, regex)
		if len(matches) > 0 {
			arr = append(arr, matches...)
		}
	}
	return
}

// extractWithRegex applies regex patterns to content string and returns all matches.
func (ind Indexer) extractWithRegex(content FileContent, regex []string) (r []string) {
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
