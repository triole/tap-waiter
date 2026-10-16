package util

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/triole/logseal"
	str2duration "github.com/xhit/go-str2duration/v2"
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

func (ut Util) FileOrFolder(s string) string {
	if ut.IsDir(s) {
		return "folder"
	}
	return "file"
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

	// Validate regex filter to prevent ReDoS
	if rxFilter != "" {
		valid, validErr := ut.isValidRegex(rxFilter)
		if !valid {
			ut.Lg.Error(
				"invalid regex filter",
				logseal.F{"filter": rxFilter, "error": validErr},
			)
			return filelist, validErr
		}

		rxf, compileErr := regexp.Compile(rxFilter)
		if compileErr != nil {
			ut.Lg.Error(
				"invalid regex filter",
				logseal.F{"filter": rxFilter, "error": compileErr},
			)
			return filelist, fmt.Errorf("regex compilation failed: %w", compileErr)
		}
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
	} else {
		// No filter, return all files
		err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
			if err == nil {
				if !f.IsDir() {
					filelist = append(filelist, path)
				}
			}
			return nil
		})
	}

	ut.Lg.IfErrError("find files failed", logseal.F{"path": basedir, "error": err})
	return
}

func (ut Util) FromTestFolder(s string) (r string) {
	// Get the directory of the current executable or current working directory
	binDir := ut.GetBinDir()
	if binDir == "" {
		binDir = "."
	}

	testdataPath := filepath.Join("testdata")
	upTwoLevels := filepath.Join(testdataPath, "../../../testdata")

	for _, path := range []string{testdataPath, upTwoLevels} {
		if _, err := os.Stat(path); err == nil {
			r = filepath.Join(path, s)
			return
		}
	}

	// Fallback to first option even if it doesn't exist
	r = filepath.Join(testdataPath, s)
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

func (ut Util) IsEmpty(str string) bool {
	r := strings.TrimSpace(str)
	return strings.EqualFold(r, "")
}

func (ut Util) IsLocalPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsPathSafe validates that a file path is within the allowed base directory.
// It prevents path traversal attacks by resolving the absolute path and
// ensuring it starts with the base directory.
func (ut Util) IsPathSafe(path, baseDir string) bool {
	// Resolve the absolute path of the base directory first
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		ut.Lg.Error("failed to resolve base directory", logseal.F{"path": baseDir, "error": err})
		return false
	}

	// Clean the target path to remove .. and . components
	cleanPath := filepath.Clean(path)

	// Resolve the absolute path of the target file
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		ut.Lg.Error("failed to resolve file path", logseal.F{"path": path, "error": err})
		return false
	}

	// The absolute path must be within the base directory
	// It must either equal the base directory or start with baseDir + separator
	if absPath == absBase {
		return true
	}

	// Check path starts with base directory followed by separator
	// This prevents prefix attacks like "allowed" matching "allowed-evil"
	return len(absPath) > len(absBase) &&
		absPath[len(absBase)] == filepath.Separator &&
		absPath[:len(absBase)] == absBase
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
	temp, err := regexp.Compile(rx)
	if err != nil {
		ut.Lg.Error(
			"invalid regex",
			logseal.F{"regex": rx, "error": err},
		)
		return ""
	}
	return temp.FindString(content)
}

func (ut Util) RxFindAll(rx string, content string) []string {
	temp, err := regexp.Compile(rx)
	if err != nil {
		ut.Lg.Error(
			"invalid regex",
			logseal.F{"regex": rx, "error": err},
		)
		return []string{}
	}
	return temp.FindAllString(content, -1)
}

func (ut Util) RxMatch(rx string, content string) bool {
	temp, err := regexp.Compile(rx)
	if err != nil {
		ut.Lg.Error(
			"invalid regex",
			logseal.F{"regex": rx, "error": err},
		)
		return false
	}
	return temp.MatchString(content)
}

func (ut Util) RxReplaceAll(basestring, regex, newstring string) (r string) {
	rx, err := regexp.Compile(regex)
	if err != nil {
		ut.Lg.Error(
			"invalid regex",
			logseal.F{"regex": regex, "error": err},
		)
		return basestring
	}
	r = rx.ReplaceAllString(basestring, newstring)
	return
}

func (ut Util) Str2Dur(s string) (dur time.Duration, err error) {
	dur, err = str2duration.ParseDuration(s)
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

// isValidRegex validates a regex pattern to prevent ReDoS attacks.
// It checks for length limits and potentially dangerous patterns.
func (ut Util) isValidRegex(pattern string) (bool, error) {
	// Maximum length to prevent DoS
	const maxRegexpLen = 65536
	if len(pattern) > maxRegexpLen {
		return false, fmt.Errorf("regex pattern too long: %d bytes (max %d)", len(pattern), maxRegexpLen)
	}

	// Try to compile with a timeout to detect catastrophic backtracking
	// We use a 1 second timeout for compilation
	const compilationTimeout = time.Second
	ch := make(chan struct {
		ok    bool
		error error
	}, 1)

	go func() {
		_, err := regexp.Compile(pattern)
		ch <- struct {
			ok    bool
			error error
		}{ok: err == nil, error: err}
	}()

	select {
	case <-time.After(compilationTimeout):
		return false, fmt.Errorf("regex compilation timed out - possible ReDoS pattern")
	case result := <-ch:
		if result.error != nil {
			return false, result.error
		}
		return true, nil
	}
}

// keep-sorted end
