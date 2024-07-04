package indexer

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

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

func (ji JoinerIndex) collectSortFiles(params Params) (sfs tSortFiles) {
	for _, el := range ji {
		sortFile := ji.readSortFile(
			filepath.Join(params.Endpoint.Folder, el.Path),
			params.Endpoint.SortFileName,
		)
		if sortFile.IsSortFile {
			sfs = append(sfs, sortFile)
		}
	}
	return
}

func (ji JoinerIndex) applySortFileOrderAndExclusion(params Params) (r JoinerIndex) {
	sortFiles := ji.collectSortFiles(params)
	var sortIndex int
	var exclude bool
	for idx, indexEl := range ji {
		exclude = false
		sortIndex = 999999
		relevantSortFile := ji.getRelevantSortFile(
			filepath.Join(params.Endpoint.Folder, indexEl.Path), sortFiles,
		)
		for orderIndex, orderEl := range relevantSortFile.Content.Order {
			if orderEl == indexEl.Path {
				sortIndex = orderIndex
				exclude = relevantSortFile.Content.Exclusive
			}
		}
		ji[idx].SortIndex = ji.stringifySortIndex(
			[]interface{}{
				sortIndex,
				ut.GetPathDepth(indexEl.Path),
				indexEl.Path,
				exclude,
			},
		)
		if relevantSortFile.Content.Exclusive {
			if ut.RxSliceContainsString(relevantSortFile.Content.Order, indexEl.Path) {
				r = append(r, ji[idx])
			}
		} else {
			r = append(r, ji[idx])
		}
	}
	sort.Sort(JoinerIndex(ji))
	return
}

func (ji JoinerIndex) getRelevantSortFile(folder string, sfs tSortFiles) (sf tSortFile) {
	for _, el := range sfs {
		if strings.HasPrefix(folder, el.Folder) {
			sf = el
		}
	}
	return
}

func (ji JoinerIndex) stringifySortIndex(li []interface{}) (r string) {
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

func (ji JoinerIndex) sortByCreated() {
	for idx, el := range ji {
		el.SortIndex = el.Created
		ji[idx] = el
	}
}

func (ji JoinerIndex) sortByLastMod() {
	for idx, el := range ji {
		el.SortIndex = el.LastMod
		ji[idx] = el
	}
}

func (ji JoinerIndex) sortBySize() {
	for idx, el := range ji {
		el.SortIndex = el.Size
		ji[idx] = el
	}
}

func (ji JoinerIndex) sortByOtherParams(params Params) {
	for idx, el := range ji {
		var val []string
		if params.SortBy != "" {
			val = ji.getContentVal(params.SortBy, el.Content)
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
		ji[idx] = el
	}
}

func (ji JoinerIndex) readSortFile(filename, sortFileName string) (sf tSortFile) {
	sf.IsSortFile = strings.HasSuffix(filename, sortFileName)
	if !sf.IsSortFile {
		return
	}
	sf.Path = filename
	sf.Folder = filepath.Dir(sf.Path)
	var by []byte
	var isTextfile bool
	by, isTextfile, sf.Error = ut.ReadFile(sf.Path)
	if sf.Error != nil && isTextfile {
		return
	} else {
		sf.Error = yaml.Unmarshal(by, &sf.Content)
	}
	sf.Folder = filepath.Dir(sf.Path)
	return
}
