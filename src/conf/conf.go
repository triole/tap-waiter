package conf

import (
	"os"
	"path"

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
	conf.Port = content.Port
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

		val.EpURL = key
		conf.API[key] = val
	}
}
