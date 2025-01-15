package indexer

import (
	"fmt"
	"sort"

	"github.com/triole/logseal"
)

func (ind *Indexer) MakeTapIndex(params Params) {
	ind.Lg.Debug(
		"start indexing and measure duration",
		logseal.F{"index_params": fmt.Sprintf("%+v", params)},
	)

	var err error
	ind.DataSources.Params = params
	ind.DataSources.Type = params.Endpoint.SourceType
	switch ind.DataSources.Type {
	case "folder":
		ind.DataSources.Paths, err = ind.Util.Find(
			ind.DataSources.Params.Endpoint.Source,
			ind.DataSources.Params.Endpoint.RxFilter,
		)
		ind.Lg.IfErrError(
			"can not identify data sources", logseal.F{"error": err},
		)
	default:
		ind.DataSources.Paths = []string{
			ind.DataSources.Params.Endpoint.Source,
		}
	}
	ind.assembleTapIndex()
	sort.Sort(TapIndex(ind.TapIndex))

	if params.Endpoint.SortFileName != "" {
		ind.TapIndex = ind.TapIndex.applySortFileOrderAndExclusion(params)
	} else {
		switch params.SortBy {
		case "created":
			ind.TapIndex.sortByCreated()
		case "lastmod":
			ind.TapIndex.sortByLastMod()
		case "size":
			ind.TapIndex.sortBySize()
		default:
			ind.TapIndex.sortByOtherParams(params)
		}
	}

	if params.Filter.Enabled {
		ind.TapIndex = ind.TapIndex.filterIndex(params)
	}

	if params.Ascending {
		sort.Sort(TapIndex(ind.TapIndex))
	} else {
		sort.Sort(sort.Reverse(TapIndex(ind.TapIndex)))
	}
	ind.TapIndex = ind.TapIndex.applyIgnoreList(params)
}

func (ti TapIndex) filterIndex(params Params) TapIndex {
	var temp TapIndex
	for _, el := range ti {
		val := ti.getContentVal(params.Filter.Prefix, el.Content)
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
	ti = temp
	return temp
}

func (ti TapIndex) applyIgnoreList(params Params) TapIndex {
	var nuJI TapIndex
	for _, el := range ti {
		if !ut.RxSliceContainsString(params.Endpoint.IgnoreList, el.Path) {
			nuJI = append(nuJI, el)
		}
	}
	return nuJI
}
