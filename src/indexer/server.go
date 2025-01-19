package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/triole/logseal"
)

func (ind Indexer) RunServer() {
	http.HandleFunc("/", ind.ServeContent)
	portstr := strconv.Itoa(ind.Conf.Port)
	ind.Lg.Info("run server, listen at :" + portstr + "/")
	err := http.ListenAndServe(":"+portstr, nil)
	if err != nil {
		panic(err)
	}
}

func (ind Indexer) ServeContent(w http.ResponseWriter, r *http.Request) {
	ind.Lg.Info("got request", logseal.F{"url": r.URL})
	params := Params{
		Ascending: true,
		Filter:    FilterParams{Enabled: false},
	}

	url, err := ind.decodeURL(r.URL.Path)
	if err == nil {
		queryParams := r.URL.Query()
		for key, values := range queryParams {
			lowKey := strings.ToLower(key)
			for _, val := range values {
				val, err = ind.decodeURL(val)
				if err == nil {
					lowVal := strings.ToLower(val)
					if lowKey == "sortby" {
						params.SortBy = lowVal
					}
					if lowKey == "order" && lowVal == "asc" {
						params.Ascending = true
					}
					if lowKey == "order" && lowVal == "desc" {
						params.Ascending = false
					}
					if lowKey == "filter" {
						params.Filter = ind.parseFilterString(val)
					}
				}
			}
		}
	}
	if val, ok := ind.Conf.API[url]; ok {
		start := time.Now()
		params.Endpoint = val
		ind.UpdateTapIndex(params)
		ind.Lg.Debug(
			"serve json",
			logseal.F{
				"url": url, "path": val.Source, "rxfilter": val.RxFilter, "duration": time.Since(start),
			},
		)
		w.Header().Add("Content Type", "application/json")
		ti := ind.getTapIndexCache(params.Endpoint.Source)
		json.NewEncoder(w).Encode(ti)
	} else {
		ind.return404(w)
	}
}

func (ind Indexer) parseFilterString(s string) (fil FilterParams) {
	fil.Prefix = ind.Util.RxFind("^[a-z0-9_\\-\\. ]+", s)
	fil.Operator = ind.Util.RxFind("^[^a-z0-9_\\-\\. ]+", strings.TrimPrefix(s, fil.Prefix))
	fil.Suffix = strings.Split(strings.TrimPrefix(s, fil.Prefix+fil.Operator), ",")
	sort.Strings(fil.Suffix)
	if fil.Prefix == "" {
		fil.Errors = append(
			fil.Errors, errors.New("error parsed filter: no match for prefix"),
		)
	}
	if fil.Operator == "" {
		fil.Errors = append(
			fil.Errors, errors.New("error parsed filter: no match for operator"),
		)
	}
	if fil.Prefix == "" {
		fil.Errors = append(
			fil.Errors, errors.New("error parsed filter: no match for suffix"),
		)
	}
	if len(fil.Errors) > 0 {
		for _, el := range fil.Errors {
			ind.Lg.Error(el)
		}
	} else {
		fil.Enabled = true
	}
	return
}

func (ind Indexer) decodeURL(s string) (t string, err error) {
	t, err = url.QueryUnescape(s)
	ind.Lg.IfErrError("can not decode url: %s, error: %s", s, err)
	return
}

func (ind Indexer) return404(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte(
		fmt.Sprintf("[ \"404 - %s\" ]", http.StatusText(404)),
	))
}
