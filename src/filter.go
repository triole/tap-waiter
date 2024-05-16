package main

import (
	"sort"
	"strings"

	"github.com/triole/logseal"
)

func equalSlices(content, filter []string) (r bool) {
	content = listToLower(content)
	filter = listToLower(filter)
	sort.Strings(content)
	r = false
	if len(filter) == len(content) {
		r = true
		for i := 0; i < len(filter); i++ {
			if filter[i] != content[i] {
				r = false
				break
			}
		}
	}
	appliedFilterTraceMessage("equal slices", content, filter, r)
	return
}

func notEqualSlices(content, filter []string) (r bool) {
	r = !equalSlices(content, filter)
	appliedFilterTraceMessage("not equal slices", content, filter, r)
	return
}

func containsSlice(content, filter []string) (r bool) {
	content = listToLower(content)
	filter = listToLower(filter)
	r = true
	for _, fil := range filter {
		if !contains(content, fil) {
			r = false
			break
		}
	}
	appliedFilterTraceMessage("contains slice", content, filter, r)
	return
}

func notContainsSlice(content, filter []string) (r bool) {
	content = listToLower(content)
	filter = listToLower(filter)
	r = true
	if len(content) >= len(filter) {
		for _, fil := range filter {
			if contains(content, fil) {
				r = false
				break
			}
		}
	}
	appliedFilterTraceMessage("not contains slice", content, filter, r)
	return
}

func rxMatchSliceAll(content, filter []string) (r bool) {
	content = listToLower(content)
	r = true
	for _, fil := range filter {
		for _, con := range content {
			if !rxMatch(fil, con) {
				r = false
				break
			}
		}
	}
	appliedFilterTraceMessage("rx match slice completely", content, filter, r)
	return
}

func rxMatchSliceOnce(content, filter []string) (r bool) {
	content = listToLower(content)
	r = true
	for _, fil := range filter {
		if !rxContains(content, fil) {
			r = false
			break
		}
	}
	appliedFilterTraceMessage("rx match slice once", content, filter, r)
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

func rxContains(slice []string, item string) bool {
	for _, s := range slice {
		if rxMatch(item, s) {
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

func listToLower(arr []string) (newArr []string) {
	for _, el := range arr {
		newArr = append(newArr, strings.ToLower(el))
	}
	return
}
