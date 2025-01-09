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

func (ind Indexer) MakeJoinerIndex(params Params) (ji JoinerIndex) {
	ind.Lg.Debug(
		"make joiner index and start measure duration",
		logseal.F{"index_params": fmt.Sprintf("%+v", params)},
	)

	switch params.Endpoint.SourceType {
	case "file":
		ji = idx.gatherFiles([]string{params.Endpoint.Source}, params)
	case "folder":
		dataFiles, err := ind.Util.Find(params.Endpoint.Source, params.Endpoint.RxFilter)
		if err == nil {
			ji = idx.gatherFiles(dataFiles, params)
		}
	case "url":
		je := JoinerEntry{Path: params.Endpoint.Source}
		if params.Endpoint.ReturnValues.Content {
			resp, err := ind.req(
				params.Endpoint.Source,
				params.Endpoint.HTTPMethod,
			)
			// TODO: maybe later add the possibility to encode base64
			je.Content = ind.byteToBody(resp)
			je.Content.Error = err
			if je.Content.Error == nil {
				je.Content = ind.unmarshal(resp, params.Endpoint)
			}
		}
		ji = append(ji, je)
	}
	sort.Sort(JoinerIndex(ji))

	if params.Endpoint.SortFileName != "" {
		ji = ji.applySortFileOrderAndExclusion(params)
	} else {
		switch params.SortBy {
		case "created":
			ji.sortByCreated()
		case "lastmod":
			ji.sortByLastMod()
		case "size":
			ji.sortBySize()
		default:
			ji.sortByOtherParams(params)
		}
	}

	if params.Filter.Enabled {
		ji = ji.filterIndex(params)
	}

	if params.Ascending {
		sort.Sort(JoinerIndex(ji))
	} else {
		sort.Sort(sort.Reverse(JoinerIndex(ji)))
	}
	ji = ji.applyIgnoreList(params)
	return
}

func (ind Indexer) gatherFiles(dataFiles []string, params Params) (ji JoinerIndex) {
	ln := len(dataFiles)
	if ln < 1 {
		ind.Lg.Warn("no data files found", logseal.F{"path": params.Endpoint.Source})
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
			li.SortIndex = ji.stringifySortIndex(
				[]interface{}{ind.Util.GetPathDepth(li.Path), li.Path},
			)
			ji = append(ji, li)
			c++
			if c >= ln {
				close(chin)
				close(chout)
				break
			}
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
				match = ut.RxSliceMatchesSliceFully(val, params.Filter.Suffix)
			case "=~":
				match = ut.RxSliceContainsSliceFully(val, params.Filter.Suffix)
			}
			if match {
				temp = append(temp, el)
			}
		}
	}
	ji = temp
	return temp
}

func (ji JoinerIndex) applyIgnoreList(params Params) JoinerIndex {
	var nuJI JoinerIndex
	for _, el := range ji {
		if !ut.RxSliceContainsString(params.Endpoint.IgnoreList, el.Path) {
			nuJI = append(nuJI, el)
		}
	}
	return nuJI
}
