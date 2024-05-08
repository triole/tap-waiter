package main

import (
	"sort"

	"github.com/triole/logseal"
)

func equalSlices(content, filter []string) (r bool) {
	sort.Strings(content)
	r = false
	if len(filter) == len(content) {
		r = true
		for i := 0; i < len(filter); i++ {
			if filter[i] != content[i] {
				r = false
			}
		}
	}
	appliedFilterTraceMessage("equal slices", content, filter, r)
	return
}

func notEqualSlices(content, filter []string) (r bool) {
	appliedFilterTraceMessage("not equal slices", content, filter, r)
	return !equalSlices(content, filter)
}

func containsSlice(content, filter []string) (r bool) {
	r = true
	for _, fil := range filter {
		if !contains(content, fil) {
			r = false
		}
	}
	appliedFilterTraceMessage("contains slice", content, filter, r)
	return
}

func notContainsSlice(content, filter []string) (r bool) {
	r = true
	if len(content) >= len(filter) {
		for _, fil := range filter {
			if contains(content, fil) {
				r = false
			}
		}
	}
	appliedFilterTraceMessage("not contains slice", content, filter, r)
	return
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func appliedFilterTraceMessage(name string, content, filter []string, r bool) {
	lg.Trace("applied filter: "+name,
		logseal.F{"content": content, "filter": filter, "result": r},
	)
}
