package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/triole/logseal"
)

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
	lg.Trace("got request", logseal.F{"endpoint": r.URL})
	key := r.URL.String()
	if val, ok := conf.API[key]; ok {
		start := time.Now()
		ji := makeJoinerIndex(val, CLI.Threads)
		lg.Info(
			"serve json",
			logseal.F{
				"url": key, "path": val.Folder, "rxfilter": val.RxFilter, "duration": time.Since(start),
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
