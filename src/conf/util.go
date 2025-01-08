package conf

import (
	"net/url"
	"os"

	"github.com/triole/logseal"
)

func (conf Conf) isURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func (conf Conf) fileOrFolder(s string) string {
	if conf.isDir(s) {
		return "folder"
	}
	return "file"
}

func (conf Conf) isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		conf.Lg.Error("can not open path", logseal.F{"error": err})
		return false
	}
	return fileInfo.IsDir()
}
