package util

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v2"
)

// keep-sorted start block=yes newline_separated=yes
func (ut Util) AbsPath(str string) (p string, err error) {
	p = str
	if !ut.IsAbs(str) {
		p, err = filepath.Abs(str)
		ut.Lg.IfErrFatal("invalid file path", logseal.F{"path": str, "error": err})
	}
	return p, err
}

func (ut Util) AbsPathSlim(str string) (p string) {
	p, _ = ut.AbsPath(str)
	return
}

func (ut Util) Find(basedir string, rxFilter string) (filelist []string, err error) {
	filelist = []string{}
	inf, err := os.Stat(basedir)
	if err != nil {
		ut.Lg.Error(
			"find files failed, unable to access folder", logseal.F{
				"path": basedir, "error": err,
			},
		)
		return
	}
	if !inf.IsDir() {
		ut.Lg.Error(
			"find files failed, not a folder, please provide a directory to look for files",
			logseal.F{"path": basedir},
		)
		return
	}

	rxf, _ := regexp.Compile(rxFilter)
	err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if rxf.MatchString(path) {
			inf, err := os.Stat(path)
			if err == nil {
				if !inf.IsDir() {
					filelist = append(filelist, path)
				}
			} else {
				ut.Lg.IfErrInfo("stat file failed", logseal.F{"path": path})
			}
		}
		return nil
	})
	ut.Lg.IfErrError("find files failed", logseal.F{"path": basedir, "error": err})
	return
}

func (ut Util) FromTestFolder(s string) (r string) {
	t, err := filepath.Abs("../../testdata")
	if err == nil {
		r = filepath.Join(t, s)
	}
	return
}

func (ut Util) GetBinDir() string {
	ex, err := os.Executable()
	ut.Lg.IfErrError(
		"unable to determine binary folder", logseal.F{"error": err},
	)
	return filepath.Dir(ex)
}

func (ut Util) GetFileLastMod(filename string) (uts int64) {
	fil, err := os.Stat(filename)
	if err != nil {
		ut.Lg.Error("get last mod failed, cannot stat file", logseal.F{"path": filename, "error": err})
		return
	}
	uts = fil.ModTime().Unix()
	return
}

func (ut Util) GetFileSize(filename string) (siz uint64) {
	file, err := os.Open(filename)
	ut.Lg.IfErrError(
		"can not open file to get file size",
		logseal.F{"path": filename, "error": err},
	)
	if err == nil {
		defer file.Close()
		stat, err := file.Stat()
		ut.Lg.IfErrError(
			"can not stat file to get file size",
			logseal.F{"path": filename, "error": err},
		)
		if err == nil {
			siz = uint64(stat.Size())
		}
	}
	return
}

func (ut Util) GetPathDepth(pth string) int {
	return len(strings.Split(pth, string(filepath.Separator))) - 1
}

func (ut Util) IsAbs(s string) bool {
	return filepath.IsAbs(s)
}

func (ut Util) IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		ut.Lg.Error("can not open path", logseal.F{"error": err})
		return false
	}
	return fileInfo.IsDir()
}

func (ut Util) IsLocalPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (ut Util) IsTextData(by []byte) bool {
	return utf8.ValidString(string(by))
}

func (ut Util) IsURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func (ut Util) ListInterfaceToListString(itf []interface{}) (r []string) {
	for _, el := range itf {
		r = append(r, fmt.Sprintf("%v", el))
	}
	return
}

func (ut Util) ReadFile(filename string) (by []byte, isTextfile bool, err error) {
	fn, err := ut.AbsPath(filename)
	if err == nil {
		by, err = os.ReadFile(fn)
		isTextfile = ut.IsTextData(by)
		ut.Lg.IfErrError(
			"can not read file", logseal.F{"path": filename, "error": err},
		)
	}
	return
}

func (ut Util) ReadYAMLFile(filepath string) (r map[string]interface{}) {
	by, _, err := ut.ReadFile(filepath)
	if err != nil {
		return
	} else {
		_ = yaml.Unmarshal(by, &r)
	}
	return
}

func (ut Util) ReverseArr(arr []string) []string {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func (ut Util) RxFind(rx string, content string) string {
	temp, _ := regexp.Compile(rx)
	return temp.FindString(content)
}

func (ut Util) RxMatch(rx string, content string) bool {
	temp, _ := regexp.Compile(rx)
	return temp.MatchString(content)
}

func (ut Util) RxReplaceAll(basestring, regex, newstring string) (r string) {
	rx := regexp.MustCompile(regex)
	r = rx.ReplaceAllString(basestring, newstring)
	return
}

func (ut Util) StringifySliceOfInterfaces(itf []interface{}) (r []string) {
	for _, el := range itf {
		r = append(r, el.(string))
	}
	return
}

func (ut Util) ToFloat(inp interface{}) (fl float64) {
	switch val := inp.(type) {
	case float32:
		fl = float64(val)
	case float64:
		fl = val
	case int:
		fl = float64(val)
	case int8:
		fl = float64(val)
	case int16:
		fl = float64(val)
	case int32:
		fl = float64(val)
	case int64:
		fl = float64(val)
	case uint:
		fl = float64(val)
	case uint8:
		fl = float64(val)
	case uint16:
		fl = float64(val)
	case uint32:
		fl = float64(val)
	case uint64:
		fl = float64(val)
	}
	return
}

func (ut Util) ToString(inp interface{}) string {
	return fmt.Sprintf("%s", inp)
}

func (ut Util) Trace() (r string) {
	pc, fullfile, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	file := ut.RxFind("src.*", fullfile)
	r = fmt.Sprintf("%s:%d %s", file, line, fn.Name())
	return
}

// keep-sorted end
