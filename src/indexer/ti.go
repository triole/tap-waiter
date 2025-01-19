package indexer

import (
	"fmt"
	"sort"

	"github.com/triole/logseal"
)

func (ind *Indexer) UpdateTapIndex(params Params) {
	ind.Lg.Debug(
		"start indexing",
		logseal.F{"index_params": fmt.Sprintf("%+v", params)},
	)

	var err error
	ti, tim := ind.getTapIndexCacheWithExpiration(params.Endpoint.Source)
	if len(ti) < 1 {
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
		ti = ind.assembleTapIndex()
		sort.Sort(TapIndex(ti))

		if params.Endpoint.SortFileName != "" {
			ti = ti.applySortFileOrderAndExclusion(params)
		} else {
			switch params.SortBy {
			case "created":
				ti.sortByCreated()
			case "lastmod":
				ti.sortByLastMod()
			case "size":
				ti.sortBySize()
			default:
				ti.sortByOtherParams(params)
			}
		}

		if params.Filter.Enabled {
			ti = ti.filterIndex(params)
		}

		if params.Ascending {
			sort.Sort(TapIndex(ti))
		} else {
			sort.Sort(sort.Reverse(TapIndex(ti)))
		}
		ti = ti.applyIgnoreList(params)
		ind.setTapIndexCache(params.Endpoint.Source, ti)
	} else {
		ind.Lg.Debug(
			"return from cache",
			logseal.F{
				"key":             params.Endpoint.Source,
				"expiration_time": tim,
			})
	}
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
