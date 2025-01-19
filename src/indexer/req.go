package indexer

import (
	"io"
	"net/http"
	"net/url"

	"github.com/triole/logseal"
)

func (ind Indexer) req(targetURL, method string) (data []byte, err error) {
	ind.Lg.Debug("fire request", logseal.F{"url": targetURL, "method": method})
	url, err := url.Parse(targetURL)
	ind.Lg.IfErrError("can not parse url", logseal.F{"error": err})

	client := &http.Client{}

	request, err := http.NewRequest(method, url.String(), nil)
	request.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (X11; Linux x86_64; rv:134.0) Gecko/20100101 Firefox/134.0",
	)
	ind.Lg.IfErrError("can not init request", logseal.F{"error": err})

	response, err := client.Do(request)
	ind.Lg.IfErrError("request failed", logseal.F{"error": err})

	if err == nil {
		data, err = io.ReadAll(response.Body)
		ind.Lg.IfErrError(
			"unable to read request response", logseal.F{"error": err},
		)
	}
	return
}
