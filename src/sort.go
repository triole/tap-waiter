package main

import (
	"fmt"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type tSortFile struct {
	Exclusive bool     `yaml:"exclusive"`
	Order     []string `yaml:"order"`
	Folder    string   `yaml:"-"`
}

func applySortFilesOrder(ji tJoinerIndex, params tIDXParams) (rji tJoinerIndex, err error) {
	for _, el := range ji {
		if strings.HasSuffix(el.Path, params.Endpoint.SortFiles) {
			sortFile, err := readSortFile(
				filepath.Join(params.Endpoint.Folder, el.Path),
			)
			fmt.Printf("%+v\n", sortFile)
			if err == nil {
				if sortFile.Exclusive {
					rji = sortExclusive(ji, sortFile)
				} else {
					rji = sortNonExclusive(ji, sortFile)
				}
			}
		}
	}
	return
}

func sortExclusive(ji tJoinerIndex, sf tSortFile) (rji tJoinerIndex) {
	rji = ji
	return
}

func sortNonExclusive(ji tJoinerIndex, sf tSortFile) (rji tJoinerIndex) {
	rji = ji
	return
}

func readSortFile(filename string) (sf tSortFile, err error) {
	var by []byte
	var isTextfile bool
	by, isTextfile, err = readFile(filename)
	if err == nil && isTextfile {
		err = yaml.Unmarshal(by, &sf)
	}
	sf.Folder = filepath.Dir(filename)
	return
}

func (arr tJoinerIndex) sortByCreated() {
	for idx, el := range arr {
		el.SortIndex = el.Created
		arr[idx] = el
	}
}

func (arr tJoinerIndex) sortByLastMod() {
	for idx, el := range arr {
		el.SortIndex = el.LastMod
		arr[idx] = el
	}
}

func (arr tJoinerIndex) sortBySize() {
	for idx, el := range arr {
		el.SortIndex = el.Size
		arr[idx] = el
	}
}

func (arr tJoinerIndex) sortByOtherParams(params tIDXParams) {
	for idx, el := range arr {
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
		arr[idx] = el
	}
	return
}
