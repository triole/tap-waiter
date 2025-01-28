package conf

import (
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/c2h5oh/datasize"
	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v3"
)

func (conf *Conf) readConfig() {
	var content ConfContent
	by, err := os.ReadFile(conf.FileName)
	conf.Lg.IfErrFatal(
		"can not read file", logseal.F{"path": conf.FileName, "error": err},
	)

	by, err = conf.TemplateFile(by)
	conf.Lg.IfErrFatal(
		"can not expand config variables", logseal.F{"path": conf.FileName, "error": err},
	)

	err = yaml.Unmarshal(by, &content)
	conf.Lg.IfErrFatal(
		"can not unmarshal config", logseal.F{"path": conf.FileName, "error": err},
	)

	conf.ServerURL = content.ServerURL
	if conf.ServerURL == "" {
		conf.ServerURL = content.Bind
	}
	if !strings.HasPrefix(conf.ServerURL, "http://") && !strings.HasPrefix(conf.ServerURL, "https://") {
		conf.ServerURL = "https://" + conf.ServerURL
	}
	_, err = url.Parse(conf.ServerURL)
	if err != nil {
		conf.Lg.IfErrFatal(
			"invalid server url", logseal.F{"path": conf.FileName, "error": err},
		)
	}

	conf.Bind = content.Bind
	if content.DefaultCacheLifetimeStr == "" {
		content.DefaultCacheLifetimeStr = "5m"
	}
	conf.DefaultCacheLifetime, err = conf.Util.Str2Dur(
		content.DefaultCacheLifetimeStr,
	)
	conf.Lg.IfErrFatal(
		"can not parse cache lifetime setting",
		logseal.F{"error": err},
	)

	for key, val := range content.API {
		key = "/" + path.Clean(key)

		var v datasize.ByteSize
		if val.MaxReturnSize == "" {
			val.MaxReturnSize = "10K"
		}
		err = v.UnmarshalText([]byte(val.MaxReturnSize))
		if err == nil {
			val.MaxReturnSizeBytes = v.Bytes()
		} else {
			conf.Lg.Fatal(
				"unable to parse config's max_return_size", logseal.F{
					"error": err,
				},
			)
		}
		val.ID = key
		conf.API[val.ID] = val
	}
}
