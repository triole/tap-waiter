package indexer

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/triole/logseal"

	yaml "gopkg.in/yaml.v3"
)

type tSortFiles []tSortFile

type tSortFile struct {
	IsSortFile bool
	Path       string
	Folder     string
	Content    tSortFileContent
	Error      error
}

type tSortFileContent struct {
	Exclusive bool     `yaml:"exclusive"`
	Order     []string `yaml:"order"`
	Folder    string   `yaml:"-"`
}

func (ti TapIndex) collectSortFiles(params Params) (sfs tSortFiles) {
	baseDir := params.Endpoint.Source
	for _, el := range ti {
		sortFile := ti.readSortFile(
			filepath.Join(params.Endpoint.Source, el.Path),
			params.Endpoint.SortFileName,
			baseDir,
		)
		if sortFile.IsSortFile {
			sfs = append(sfs, sortFile)
		}
	}
	return
}

func (ti TapIndex) applySortFileOrderAndExclusion(params Params) (r TapIndex) {
	sortFiles := ti.collectSortFiles(params)
	var sortIndex int
	var exclude bool
	for idx, indexEl := range ti {
		exclude = false
		sortIndex = 999999
		relevantSortFile := ti.getRelevantSortFile(
			filepath.Join(params.Endpoint.Source, indexEl.Path), sortFiles,
		)
		for orderIndex, orderEl := range relevantSortFile.Content.Order {
			if orderEl == indexEl.Path {
				sortIndex = orderIndex
				exclude = relevantSortFile.Content.Exclusive
			}
		}
		ti[idx].SortIndex = ti.stringifySortIndex(
			[]interface{}{
				sortIndex,
				ut.GetPathDepth(indexEl.Path),
				indexEl.Path,
				exclude,
			},
		)
		if relevantSortFile.Content.Exclusive {
			if ut.RxSliceContainsString(relevantSortFile.Content.Order, indexEl.Path) {
				r = append(r, ti[idx])
			}
		} else {
			r = append(r, ti[idx])
		}
	}
	sort.Sort(TapIndex(ti))
	return
}

func (ti TapIndex) getRelevantSortFile(folder string, sfs tSortFiles) (sf tSortFile) {
	for _, el := range sfs {
		if strings.HasPrefix(folder, el.Folder) {
			sf = el
		}
	}
	return
}

func (ti TapIndex) stringifySortIndex(li []interface{}) (r string) {
	sep := "|"
	for _, itf := range li {
		switch val := itf.(type) {
		case int:
			r += sep + fmt.Sprintf("%06d", val)
		case bool:
			if val {
				r += sep + fmt.Sprintf("%06d", 1)
			} else {
				r += sep + fmt.Sprintf("%06d", 0)
			}
		default:
			r += sep + val.(string)
		}
	}
	if len(r) > 1 {
		r = r[1:]
	}
	return
}

func (ti TapIndex) sortByCreated() {
	for idx, el := range ti {
		el.SortIndex = el.Created
		ti[idx] = el
	}
}

func (ti TapIndex) sortByLastMod() {
	for idx, el := range ti {
		el.SortIndex = el.LastMod
		ti[idx] = el
	}
}

func (ti TapIndex) sortBySize() {
	for idx, el := range ti {
		el.SortIndex = el.Size
		ti[idx] = el
	}
}

func (ti TapIndex) sortByOtherParams(params Params) {
	for idx, el := range ti {
		var val []string
		if params.SortBy != "" {
			val = ti.getContentVal(params.SortBy, el.Content)
		}
		if len(val) > 0 {
			el.SortIndex = strings.Join(val, ".")
		} else {
			prefix := ""
			if params.SortBy != "" {
				prefix = "zzzzz_"
			}
			el.SortIndex = fmt.Sprintf(
				"%s%05d_%s", prefix, ut.GetPathDepth(el.Path), el.Path,
			)
		}
		ti[idx] = el
	}
}

// isValidSortFileName validates that the sort file name is safe to prevent path traversal.
// Only allows a specific filename without path components.
func (ind *Indexer) isValidSortFileName(sortFileName string) bool {
	if sortFileName == "" {
		return true // Empty means no sort file, which is valid
	}
	// Block path traversal patterns - check for .. , \, and /
	if containsPathTraversal(sortFileName) {
		ind.Lg.Warn(
			"path traversal detected in sort file name",
			logseal.F{"sortFileName": sortFileName},
		)
		return false
	}
	// Only allow alphanumeric, underscore, hyphen, dot, and common extensions
	// Maximum length of 100 characters to prevent DoS
	if len(sortFileName) > 100 {
		ind.Lg.Warn(
			"sort file name too long",
			logseal.F{"length": len(sortFileName)},
		)
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, sortFileName)
	return matched
}

func (ti TapIndex) readSortFile(filename string, sortFileName string, baseDir string) (sf tSortFile) {
	sf.IsSortFile = strings.HasSuffix(filename, sortFileName)
	if !sf.IsSortFile {
		return
	}
	sf.Path = filename

	// Validate path is within base directory to prevent traversal
	if !ut.IsPathSafe(sf.Path, baseDir) {
		sf.Error = fmt.Errorf("sort file outside allowed directory: %s", sf.Path)
		return
	}

	sf.Folder = filepath.Dir(sf.Path)
	var by []byte
	var isTextfile bool
	by, isTextfile, sf.Error = ut.ReadFile(sf.Path)
	if sf.Error != nil && isTextfile {
		return
	}
	if isTextfile && sf.Error == nil {
		sf.Error = yaml.Unmarshal(by, &sf.Content)
	}
	sf.Folder = filepath.Dir(sf.Path)
	return
}
