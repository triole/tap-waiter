package indexer

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/triole/logseal"
)

func (ind Indexer) makeAbsURL(targetURL string) (pURL *url.URL, err error) {
	pURL, err = url.Parse(targetURL)
	ind.Lg.IfErrError("can not parse url", logseal.F{"error": err})
	if !pURL.IsAbs() {
		var fullURL string
		fullURL, err = url.JoinPath(ind.Conf.ServerURL, targetURL)
		pURL, err = ind.makeAbsURL(fullURL)
	}
	return
}

func (ind Indexer) req(targetURL, method string) (data []byte, err error) {
	var parsURL *url.URL
	var requ *http.Request
	var resp *http.Response
	method = strings.ToUpper(method)
	parsURL, err = ind.makeAbsURL(targetURL)
	ind.Lg.Info("fire request", logseal.F{"url": parsURL, "method": method})
	if err == nil {
		client := &http.Client{}
		requ, err = http.NewRequest(method, parsURL.String(), nil)
		requ.Header.Set(
			"User-Agent",
			"Mozilla/5.0 (X11; Linux x86_64; rv:134.0) Gecko/20100101 Firefox/134.0",
		)
		ind.Lg.IfErrError("can not init request", logseal.F{"error": err})
		if err == nil {
			resp, err = client.Do(requ)
			ind.Lg.IfErrError("request failed", logseal.F{"error": err})
			if err == nil {
				data, err = io.ReadAll(resp.Body)
				ind.Lg.IfErrError(
					"unable to read request response", logseal.F{"error": err},
				)
			}
		}
	}
	return
}
