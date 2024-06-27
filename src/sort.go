package main

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

func (ji tJoinerIndex) collectSortFiles(params tIDXParams) (sfs tSortFiles) {
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

func (ji tJoinerIndex) applySortFileOrder(params tIDXParams) {
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
				getDepth(indexEl.Path),
				indexEl.Path,
				exclude,
			},
		)
	}
	sort.Sort(tJoinerIndex(ji))

	fmt.Printf("\n\n\n%+v\n", "DONE DONE DONE")
	for _, el := range ji {
		fmt.Printf("%20s --- %40s\n", el.Path, el.SortIndex)
	}
}

func (ji tJoinerIndex) fileOnIgnoreList(str string, list []string) bool {
	for _, el := range list {
		fmt.Printf("%+v %+v\n", el, str)
		if strings.HasSuffix(str, el) {
			return true
		}
	}
	return false
}

func (ji tJoinerIndex) applyExclusions(ignoreList []string) (newji tJoinerIndex) {
	for _, el := range ji {
		if strings.HasSuffix(el.SortIndex.(string), "|000000") && !ji.fileOnIgnoreList(el.Path, ignoreList) {
			newji = append(newji, el)
		}
	}
	return
}

func (ji tJoinerIndex) getRelevantSortFile(folder string, sfs tSortFiles) (sf tSortFile) {
	for _, el := range sfs {
		if strings.HasPrefix(folder, el.Folder) {
			sf = el
		}
	}
	return
}

func (ji tJoinerIndex) stringifySortIndex(li []interface{}) (r string) {
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

func (ji tJoinerIndex) sortByCreated() {
	for idx, el := range ji {
		el.SortIndex = el.Created
		ji[idx] = el
	}
}

func (ji tJoinerIndex) sortByLastMod() {
	for idx, el := range ji {
		el.SortIndex = el.LastMod
		ji[idx] = el
	}
}

func (ji tJoinerIndex) sortBySize() {
	for idx, el := range ji {
		el.SortIndex = el.Size
		ji[idx] = el
	}
}

func (ji tJoinerIndex) sortByOtherParams(params tIDXParams) {
	for idx, el := range ji {
		var val []string
		if params.SortBy != "" {
			val = getContentVal(params.SortBy, el.Content)
		}
		if len(val) > 0 {
			el.SortIndex = strings.Join(val, ".")
		} else {
			prefix := ""
			if params.SortBy != "" {
				prefix = "zzzzz_"
			}
			el.SortIndex = fmt.Sprintf(
				"%s%05d_%s", prefix, getDepth(el.Path), el.Path,
			)
		}
		ji[idx] = el
	}
}

func (ji tJoinerIndex) readSortFile(filename, sortFileName string) (sf tSortFile) {
	sf.IsSortFile = strings.HasSuffix(filename, sortFileName)
	if !sf.IsSortFile {
		return
	}
	sf.Path = filename
	sf.Folder = filepath.Dir(sf.Path)
	var by []byte
	var isTextfile bool
	by, isTextfile, sf.Error = readFile(sf.Path)
	if sf.Error != nil && isTextfile {
		return
	} else {
		sf.Error = yaml.Unmarshal(by, &sf.Content)
	}
	sf.Folder = filepath.Dir(sf.Path)
	return
}
