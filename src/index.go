package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/triole/logseal"
)

type tJoinerEntry struct {
	Path      string                 `json:"path"`
	SplitPath []string               `json:"split_path,omitempty"`
	Size      uint64                 `json:"size,omitempty"`
	LastMod   int64                  `json:"lastmod,omitempty"`
	Created   int64                  `json:"created,omitempty"`
	Content   map[string]interface{} `json:"content,omitempty"`
	SortIndex interface{}            `json:"-"`
}

type tJoinerIndex []tJoinerEntry

func (arr tJoinerIndex) Len() int {
	return len(arr)
}

func (arr tJoinerIndex) Less(i, j int) bool {
	switch arr[i].SortIndex.(type) {
	case float32, float64,
		int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return toFloat(arr[i].SortIndex) < toFloat(arr[j].SortIndex)
	default:
		return arr[i].SortIndex.(string) < arr[j].SortIndex.(string)
	}
}

func (arr tJoinerIndex) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func getDepth(pth string) int {
	return len(strings.Split(pth, string(filepath.Separator))) - 1
}

func makeJoinerIndex(params tIDXParams) (joinerIndex tJoinerIndex) {
	lg.Debug(
		"make joiner index and start measure duration",
		logseal.F{"index_params": fmt.Sprintf("%+v", params)},
	)
	dataFiles := find(params.Endpoint.Folder, params.Endpoint.RxFilter)
	ln := len(dataFiles)

	if ln < 1 {
		lg.Warn("no data files found", logseal.F{"path": params.Endpoint.Folder})
	} else {
		chin := make(chan string, params.Threads)
		chout := make(chan tJoinerEntry, params.Threads)

		lg.Debug(
			"start to process files",
			logseal.F{"no": ln, "threads": params.Threads},
		)

		for _, fil := range dataFiles {
			go readDataFile(fil, params.Endpoint, chin, chout)
		}

		c := 0
		for li := range chout {
			switch params.SortBy {
			case "created":
				li.SortIndex = li.Created
			case "lastmod":
				li.SortIndex = li.LastMod
			case "size":
				li.SortIndex = li.Size
			default:
				val := getMapVal(params.SortBy, li.Content)
				if len(val) > 0 {
					li.SortIndex = strings.Join(val, ".")
				} else {
					prefix := ""
					if params.SortBy != "" {
						prefix = "zzzzz_"
					}
					li.SortIndex = fmt.Sprintf("%s%05d_%s", prefix, getDepth(li.Path), li.Path)
				}
			}
			joinerIndex = append(joinerIndex, li)
			c++
			if c >= ln {
				close(chin)
				close(chout)
				break
			}
		}
		if params.Filter.Enabled {
			joinerIndex = filterJoinerIndex(joinerIndex, params)
		}
		if params.Ascending {
			sort.Sort(tJoinerIndex(joinerIndex))
		} else {
			sort.Sort(sort.Reverse(tJoinerIndex(joinerIndex)))
		}
	}
	return
}

func filterJoinerIndex(arr tJoinerIndex, params tIDXParams) (newArr tJoinerIndex) {
	newArr = []tJoinerEntry{}
	for _, el := range arr {
		val := getMapVal(params.Filter.Prefix, el.Content)
		match := false
		if len(val) > 0 {
			switch params.Filter.Operator {
			case "===":
				match = equalSlices(val, params.Filter.Suffix)
			case "!==":
				match = notEqualSlices(val, params.Filter.Suffix)
			case "==":
				match = containsSlice(val, params.Filter.Suffix)
			case "!=":
				match = notContainsSlice(val, params.Filter.Suffix)
			case "==~":
				match = rxMatchSliceCompletely(val, params.Filter.Suffix)
			case "=~":
				match = rxMatchSliceOnce(val, params.Filter.Suffix)
			}
			if match {
				newArr = append(newArr, el)
			}
		}
	}
	return
}

func getMapVal(key string, dict map[string]interface{}) (s []string) {
	pth := strings.Split(key, ".")
	if val, ok := dict[pth[0]]; ok {
		switch vl := val.(type) {
		case map[string]interface{}:
			s = getMapVal(strings.Join(pth[1:], "."), vl)
		case []interface{}:
			for _, x := range vl {
				s = append(s, x.(string))
			}
		default:
			s = []string{vl.(string)}
		}
	}
	return
}
