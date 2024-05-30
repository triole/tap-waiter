package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v3"
)

type tJoinerEntry struct {
	Path      string      `json:"path"`
	SplitPath []string    `json:"split_path,omitempty"`
	Size      uint64      `json:"size,omitempty"`
	LastMod   int64       `json:"lastmod,omitempty"`
	Created   int64       `json:"created,omitempty"`
	Content   tContent    `json:"content,omitempty"`
	SortIndex interface{} `json:"-"`
}

type tJoinerIndex []tJoinerEntry

func (arr tJoinerIndex) Len() int {
	return len(arr)
}

func (arr tJoinerIndex) Less(i, j int) bool {
	switch arr[i].SortIndex.(type) {
	case float32, float64,
		int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if arr[i].SortIndex == arr[j].SortIndex {
			return arr[i].Path > arr[j].Path
		}
		return toFloat(arr[i].SortIndex) < toFloat(arr[j].SortIndex)
	default:
		if arr[i].SortIndex.(string) == arr[j].SortIndex.(string) {
			return arr[i].Path > arr[j].Path
		}
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
				var val []string
				if params.SortBy != "" {
					val = getContentVal(params.SortBy, li.Content)
				}
				if len(val) > 0 {
					li.SortIndex = strings.Join(val, ".")
				} else {
					prefix := ""
					if params.SortBy != "" {
						prefix = "zzzzz_"
					}
					li.SortIndex = fmt.Sprintf(
						"%s%05d_%s", prefix, getDepth(li.Path), li.Path,
					)
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
		var err error
		if params.Endpoint.SafFile != "" {
			joinerIndex, err = applySafFile(joinerIndex, params.Endpoint.SafFile)
		}
		if params.Endpoint.SafFile == "" || err != nil {
			if params.Ascending {
				sort.Sort(tJoinerIndex(joinerIndex))
			} else {
				sort.Sort(sort.Reverse(tJoinerIndex(joinerIndex)))
			}
			if params.Filter.Enabled {
				joinerIndex = filterJoinerIndex(joinerIndex, params)
			}
		}
	}
	return
}

func filterJoinerIndex(arr tJoinerIndex, params tIDXParams) (newArr tJoinerIndex) {
	newArr = []tJoinerEntry{}
	for _, el := range arr {
		val := getContentVal(params.Filter.Prefix, el.Content)
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
				match = rxMatchSliceAll(val, params.Filter.Suffix)
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

func applySafFile(ji tJoinerIndex, safPath string) (rji tJoinerIndex, err error) {
	var safList []string
	safList, err = readSafFile(safPath)
	if err == nil {
		for _, safEntry := range safList {
			for _, jiEntry := range ji {
				if safEntry == jiEntry.Path {
					rji = append(rji, jiEntry)
				}
			}
		}
	} else {
		lg.Error(
			"saf file reading failed, do not apply any sort or filter",
			logseal.F{"saf_file": safPath, "error": err})
		rji = ji
	}
	return
}

func readSafFile(filename string) (r []string, err error) {
	var by []byte
	var isTextfile bool
	by, isTextfile, err = readFile(filename)
	if err == nil && isTextfile {
		err = yaml.Unmarshal(by, &r)
	}
	return
}
