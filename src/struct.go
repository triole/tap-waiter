package main

import (
	"fmt"
	"time"
)

type tJoinerEntry struct {
	Path         string                 `json:"path"`
	Depth        int                    `json:"depth"`
	Ext          string                 `json:"ext"`
	FileMetadata tFileMeta              `json:"file_metadata,omitempty"`
	Content      map[string]interface{} `json:"content,omitempty"`
}

type tJoinerIndex []tJoinerEntry

func (arr tJoinerIndex) Len() int {
	return len(arr)
}

func (arr tJoinerIndex) Less(i, j int) bool {
	si1 := fmt.Sprintf("%05d_%s", arr[i].Depth, arr[i].Path)
	si2 := fmt.Sprintf("%05d_%s", arr[j].Depth, arr[j].Path)
	return si1 < si2
}

func (arr tJoinerIndex) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

type tDateTime struct {
	Time time.Time `json:"time,omitempty"`
	Unix int64     `json:"unix,omitempty"`
}

type tFileMeta struct {
	LastMod tDateTime `json:"lastmod,omitempty"`
	Created tDateTime `json:"created,omitempty"`
}
