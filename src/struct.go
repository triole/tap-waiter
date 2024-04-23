package main

type tJoinerEntry struct {
	Path string                 `json:"path"`
	Data map[string]interface{} `json:"data"`
}

type tJoinerIndex []tJoinerEntry
