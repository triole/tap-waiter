package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type tJoinerEntry struct {
	Path        string                 `json:"path"`
	Size        int64                  `json:"size,omitempty"`
	FileLastMod int64                  `json:"file_lastmod,omitempty"`
	FileCreated int64                  `json:"file_created,omitempty"`
	FrontMatter map[string]interface{} `json:"front_matter,omitempty"`
	Content     map[string]interface{} `json:"content,omitempty"`
}

type tJoinerIndex []tJoinerEntry

func (arr tJoinerIndex) Len() int {
	return len(arr)
}

func getDepth(pth string) int {
	return len(strings.Split(pth, string(filepath.Separator))) - 1
}

func (arr tJoinerIndex) Less(i, j int) bool {
	si1 := fmt.Sprintf("%05d_%s", getDepth(arr[i].Path), arr[i].Path)
	si2 := fmt.Sprintf("%05d_%s", getDepth(arr[j].Path), arr[j].Path)
	return si1 < si2
}

func (arr tJoinerIndex) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}
