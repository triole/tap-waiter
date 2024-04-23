package main

import "time"

type tJoinerEntry struct {
	Path     string                 `json:"path"`
	Ext      string                 `json:"ext"`
	FileMeta tFileMeta              `json:"meta"`
	Data     map[string]interface{} `json:"data"`
}

type tJoinerIndex []tJoinerEntry

type tDateTime struct {
	Time time.Time
	Unix int64
}

type tFileMeta struct {
	LastMod tDateTime
	Created tDateTime
}
