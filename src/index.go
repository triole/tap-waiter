package main

import (
	"fmt"
	"sort"

	"github.com/triole/logseal"
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

func (ji tJoinerIndex) Len() int {
	return len(ji)
}

func (ji tJoinerIndex) Less(i, j int) bool {
	switch ji[i].SortIndex.(type) {
	case float32, float64,
		int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if ji[i].SortIndex == ji[j].SortIndex {
			return ji[i].Path > ji[j].Path
		}
		return toFloat(ji[i].SortIndex) < toFloat(ji[j].SortIndex)
	default:
		if ji[i].SortIndex.(string) == ji[j].SortIndex.(string) {
			return ji[i].Path > ji[j].Path
		}
		return ji[i].SortIndex.(string) < ji[j].SortIndex.(string)
	}
}

func (ji tJoinerIndex) Swap(i, j int) {
	ji[i], ji[j] = ji[j], ji[i]
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
			li.SortIndex = joinerIndex.stringifySortIndex(
				[]interface{}{getDepth(li.Path), li.Path},
			)
			joinerIndex = append(joinerIndex, li)
			c++
			if c >= ln {
				close(chin)
				close(chout)
				break
			}
		}
		sort.Sort(tJoinerIndex(joinerIndex))

		// apply sort
		if params.Endpoint.SortFileName != "" {
			joinerIndex.applySortFileOrder(params)
		}

		switch params.SortBy {
		case "created":
			joinerIndex.sortByCreated()
		case "lastmod":
			joinerIndex.sortByLastMod()
		case "size":
			joinerIndex.sortBySize()
		default:
			joinerIndex.sortByOtherParams(params)
		}

		if params.Filter.Enabled {
			joinerIndex = joinerIndex.filterIndex(params)
		}

		if params.Ascending {
			sort.Sort(tJoinerIndex(joinerIndex))
		} else {
			sort.Sort(sort.Reverse(tJoinerIndex(joinerIndex)))
		}
	}
	return
}

func (ji tJoinerIndex) filterIndex(params tIDXParams) tJoinerIndex {
	var temp tJoinerIndex
	for _, el := range ji {
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
				temp = append(temp, el)
			}
		}
	}
	ji = temp
	return temp
}
