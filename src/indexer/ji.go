package indexer

import (
	"fmt"
	"sort"

	"github.com/triole/logseal"
)

func (ji JoinerIndex) Len() int {
	return len(ji)
}

func (ji JoinerIndex) Less(i, j int) bool {
	switch ji[i].SortIndex.(type) {
	case float32, float64,
		int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if ji[i].SortIndex == ji[j].SortIndex {
			return ji[i].Path > ji[j].Path
		}
		return ut.ToFloat(ji[i].SortIndex) < ut.ToFloat(ji[j].SortIndex)
	default:
		if ji[i].SortIndex.(string) == ji[j].SortIndex.(string) {
			return ji[i].Path > ji[j].Path
		}
		return ji[i].SortIndex.(string) < ji[j].SortIndex.(string)
	}
}

func (ji JoinerIndex) Swap(i, j int) {
	ji[i], ji[j] = ji[j], ji[i]
}

func (ind Indexer) MakeJoinerIndex(params Params) (joinerIndex JoinerIndex) {
	ind.Lg.Debug(
		"make joiner index and start measure duration",
		logseal.F{"index_params": fmt.Sprintf("%+v", params)},
	)
	dataFiles := ind.Util.Find(params.Endpoint.Folder, params.Endpoint.RxFilter)
	ln := len(dataFiles)

	if ln < 1 {
		ind.Lg.Warn("no data files found", logseal.F{"path": params.Endpoint.Folder})
	} else {
		chin := make(chan string, params.Threads)
		chout := make(chan JoinerEntry, params.Threads)

		ind.Lg.Debug(
			"start to process files",
			logseal.F{"no": ln, "threads": params.Threads},
		)

		for _, fil := range dataFiles {
			go ind.readDataFile(fil, params.Endpoint, chin, chout)
		}

		c := 0
		for li := range chout {
			li.SortIndex = joinerIndex.stringifySortIndex(
				[]interface{}{ind.Util.GetPathDepth(li.Path), li.Path},
			)
			joinerIndex = append(joinerIndex, li)
			c++
			if c >= ln {
				close(chin)
				close(chout)
				break
			}
		}
		sort.Sort(JoinerIndex(joinerIndex))

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
			sort.Sort(JoinerIndex(joinerIndex))
		} else {
			sort.Sort(sort.Reverse(JoinerIndex(joinerIndex)))
		}
	}
	return
}

func (ji JoinerIndex) filterIndex(params Params) JoinerIndex {
	var temp JoinerIndex
	for _, el := range ji {
		val := ji.getContentVal(params.Filter.Prefix, el.Content)
		match := false
		if len(val) > 0 {
			switch params.Filter.Operator {
			case "===":
				match = ut.SlicesEqual(val, params.Filter.Suffix)
			case "!==":
				match = ut.SlicesNotEqual(val, params.Filter.Suffix)
			case "==":
				match = ut.SliceContainsSlice(val, params.Filter.Suffix)
			case "!=":
				match = ut.SliceNotContainsSlice(val, params.Filter.Suffix)
			case "==~":
				match = ut.RxMatchSliceAll(val, params.Filter.Suffix)
			case "=~":
				match = ut.RxMatchSliceOnce(val, params.Filter.Suffix)
			}
			if match {
				temp = append(temp, el)
			}
		}
	}
	ji = temp
	return temp
}
