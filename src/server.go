package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/triole/logseal"
)

type tIDXParams struct {
	Endpoint  tEndpoint
	SortBy    string
	Ascending bool
	Threads   int
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
	idxParams := tIDXParams{
		SortBy:    "path",
		Ascending: true,
		Threads:   CLI.Threads,
	}

	lg.Trace("got request", logseal.F{"endpoint": r.URL})
	url := r.URL.Path

	params := r.URL.Query()
	for key, values := range params {
		for _, value := range values {
			if key == "sortby" {
				idxParams.SortBy = value
			}
			if key == "order" && value == "asc" {
				idxParams.Ascending = true
			}
			if key == "order" && value == "desc" {
				idxParams.Ascending = false
			}
		}
	}
	fmt.Printf("%+v\n", url)
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

func return404(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte(
		fmt.Sprintf("[ \"404 - %s\" ]", http.StatusText(404)),
	))
}
