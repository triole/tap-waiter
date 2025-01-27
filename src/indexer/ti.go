package indexer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/triole/logseal"
)

func (ind *Indexer) updateParams(params Params) Params {
	src := params.Endpoint.Source
	params.Method = params.Endpoint.Method
	params.Response = params.Endpoint.Response
	if ind.Util.IsURL(params.Endpoint.Source) {
		params.Type = "url"
	}
	if ind.Conf.Util.IsLocalPath(params.Endpoint.Source) {
		src, _ = ind.Conf.Util.AbsPath(params.Endpoint.Source)
		params.Type = ind.Util.FileOrFolder(src)
	}
	if len(params.Sources) < 1 {
		params.Sources = []string{src}
	}
	if params.Type == "url" {
		if params.Method == "" {
			params.Method = "GET"
		}
		params.Method = strings.ToUpper(params.Method)
		params.Method = strings.TrimPrefix(params.Method, "HTTP_")
	}
	return params
}

func (ind *Indexer) updateTapIndex(params Params) {
	params = ind.updateParams(params)
	ind.Lg.Debug(
		"start updating index",
		logseal.F{"index_params": fmt.Sprintf("%+v", params)},
	)
	var err error
	ti, tim := ind.getTapIndexCacheWithExpiration(params.Endpoint.ID)
	if len(ti) < 1 {
		if !ut.IsEmpty(params.Response) {
			content := ind.unmarshal([]byte(params.Response), params.Endpoint)

			te := TapEntry{
				Path:    params.Endpoint.ID,
				Content: content,
			}
			ti = append(ti, te)
		} else {
			switch params.Type {
			case "folder":
				params.Sources, err = ind.Util.Find(
					params.Endpoint.Source,
					params.Endpoint.RxFilter,
				)
				ind.Lg.IfErrError(
					"can not identify data sources", logseal.F{"error": err},
				)
			}
			ti = ind.assembleTapIndex(params)
			if params.Endpoint.Process.Strategy == "use_as_url_list" {
				ti = ind.applyJSONPath(ti, params.Endpoint.Process.JSONPath)
				urls := ind.returnRegexMatchFullIndex(ti, params.Endpoint.Process.RegexMatch)
				if len(urls) > 0 {
					params.Sources = urls
					params.Method = params.Endpoint.Process.Method
					params = ind.updateParams(params)
					ti = ind.assembleTapIndex(params)
				} else {
					ind.Lg.Warn("regex match was empty, skip building index")
				}
			}
			ti = ind.applyJSONPath(ti, params.Endpoint.Return.JSONPath)
		}

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
		ind.setTapIndexCache(params.Endpoint.ID, ti)
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
