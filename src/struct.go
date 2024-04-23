package main

import (
	"fmt"
	"time"
)

type tJoinerEntry struct {
	Path     string                 `json:"path"`
	Depth    int                    `json:"depth"`
	Ext      string                 `json:"ext"`
	FileMeta tFileMeta              `json:"file_metadata"`
	Data     map[string]interface{} `json:"data"`
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
	Time time.Time
	Unix int64
}

type tFileMeta struct {
	LastMod tDateTime
	Created tDateTime
}
