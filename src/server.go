package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/triole/logseal"
)

type tIDXParams struct {
	Endpoint  tEndpoint
	Filter    tIDXParamsFilter
	SortBy    string
	Ascending bool
	Threads   int
}

type tIDXParamsFilter struct {
	Prefix   string
	Operator string
	Suffix   string
	Errors   []error
}

func runServer(conf tConf) {
	http.HandleFunc("/", serveContent)
	portstr := strconv.Itoa(conf.Port)
	lg.Info("run server, listen at :" + portstr + "/")
	err := http.ListenAndServe(":"+portstr, nil)
	if err != nil {
		panic(err)
	}
}

func serveContent(w http.ResponseWriter, r *http.Request) {
	lg.Info("got request", logseal.F{"url": r.URL})
	idxParams := tIDXParams{
		SortBy:    "path",
		Ascending: true,
		Threads:   CLI.Threads,
	}

	url := r.URL.Path
	params := r.URL.Query()
	for key, values := range params {
		lowKey := strings.ToLower(key)
		for _, val := range values {
			lowVal := strings.ToLower(val)
			if lowKey == "sortby" {
				idxParams.SortBy = lowVal
			}
			if lowKey == "order" && lowVal == "asc" {
				idxParams.Ascending = true
			}
			if lowKey == "order" && lowVal == "desc" {
				idxParams.Ascending = false
			}
			if lowKey == "filter" {
				idxParams.Filter = parseFilterString(val)
			}
		}
	}

	if val, ok := conf.API[url]; ok {
		idxParams.Endpoint = val
		start := time.Now()
		ji := makeJoinerIndex(idxParams)
		lg.Debug(
			"serve json",
			logseal.F{
				"url": url, "path": val.Folder, "rxfilter": val.RxFilter, "duration": time.Since(start),
			},
		)
		w.Header().Add("Content Type", "application/json")
		json.NewEncoder(w).Encode(ji)
	} else {
		return404(w)
	}
}

func parseFilterString(s string) (fil tIDXParamsFilter) {
	fil.Prefix = decodeURL(rxFind("^[a-z0-9_\\-\\. ]+", s))
	fil.Operator = decodeURL(
		rxFind("^[^a-z0-9_\\-\\. ]+", strings.TrimPrefix(s, fil.Prefix)),
	)
	fil.Suffix = decodeURL(strings.TrimPrefix(s, fil.Prefix+fil.Operator))
	if fil.Prefix == "" {
		fil.Errors = append(
			fil.Errors, errors.New("error parsed filter: prefix empty"),
		)
	}
	if fil.Operator == "" {
		fil.Errors = append(
			fil.Errors, errors.New("error parsed filter: operator empty"),
		)
	}
	if fil.Prefix == "" {
		fil.Errors = append(
			fil.Errors, errors.New("error parsed filter: suffix empty"),
		)
	}
	return
}

func decodeURL(s string) (t string) {
	t, _ = url.QueryUnescape(s)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	return
}

func return404(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte(
		fmt.Sprintf("[ \"404 - %s\" ]", http.StatusText(404)),
	))
}
