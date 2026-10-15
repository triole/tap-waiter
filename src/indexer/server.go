package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/triole/logseal"
)

func (ind Indexer) RunServer() {
	http.HandleFunc("/", ind.serveContent)
	ind.Lg.Info(
		"run server",
		logseal.F{
			"server_url":             ind.Conf.ServerURL,
			"bind":                   ind.Conf.Bind,
			"default_cache_lifetime": ind.Conf.DefaultCacheLifetime,
		},
	)
	err := http.ListenAndServe(ind.Conf.Bind, nil)
	if err != nil {
		ind.Lg.Fatal(
			"server failed",
			logseal.F{"error": err},
		)
	}
}

func (ind Indexer) serveContent(w http.ResponseWriter, r *http.Request) {
	// Set security headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	ind.Lg.Debug("got request", logseal.F{"url": r.URL})
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
						// Validate SortBy against whitelist to prevent path traversal
						if ind.isValidSortBy(lowVal) {
							params.SortBy = lowVal
							ind.Lg.Trace("got query param sort",
								logseal.F{
									"sort": lowVal,
								},
							)
						} else {
							ind.Lg.Warn("invalid sort parameter, blocked",
								logseal.F{"sort": lowVal},
							)
						}
					}
					if lowKey == "order" && lowVal == "asc" {
						params.Ascending = true
					}
					if lowKey == "order" && lowVal == "desc" {
						params.Ascending = false
					}
					if lowKey == "filter" {
						filterParams := ind.parseFilterString(val)
						if ind.isValidFilterParams(filterParams) {
							params.Filter = filterParams
							ind.Lg.Trace("got query param filter",
								logseal.F{"filter": params.Filter},
							)
						} else {
							ind.Lg.Warn("invalid filter parameter, blocked",
								logseal.F{"filter": lowVal},
							)
						}
					}
				}
			}
		}
	}
	if val, ok := ind.Conf.API[url]; ok {
		start := time.Now()
		params.Endpoint = val
		ind.updateTapIndex(params)
		ind.Lg.Info(
			"serve json",
			logseal.F{
				"url": url, "path": val.Source, "rxfilter": val.RxFilter, "duration": time.Since(start),
			},
		)
		w.Header().Add("Content Type", "application/json")
		ti := ind.lookupTapIndexCache(params.Endpoint.ID)
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
	if len(fil.Suffix) == 0 {
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

// decodeURL safely decodes a URL-encoded string while blocking dangerous control characters.
// Returns an empty string and error if decoding fails or dangerous characters are detected.
func (ind Indexer) decodeURL(s string) (t string, err error) {
	// Decode the URL-encoded string
	t, err = url.QueryUnescape(s)
	if err != nil {
		ind.Lg.IfErrError("can not decode url: %s, error: %s", s, err)
		return "", err
	}

	// Block control characters (ASCII 0-31 and 127) that could be used for injection
	for i, r := range t {
		if r < 32 || r == 127 {
			ind.Lg.Warn(
				"control character detected in URL parameter, blocking",
				logseal.F{"position": i, "char": r, "original": s, "decoded": t},
			)
			return "", fmt.Errorf("control character %U blocked at position %d", r, i)
		}
	}

	// Block null bytes which could be used for string truncation attacks
	if strings.Contains(t, "\x00") {
		ind.Lg.Warn(
			"null byte detected in URL parameter, blocking",
			logseal.F{"original": s, "decoded": t},
		)
		return "", fmt.Errorf("null byte blocked")
	}

	return t, nil
}

func (ind Indexer) return404(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write(fmt.Appendf(nil, "[ \"404 - %s\" ]", http.StatusText(404)))
}
