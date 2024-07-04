package util

import (
	"sort"
	"strings"

	"github.com/triole/logseal"
)

// keep-sorted start
func (ut Util) SlicesEqual(content, filter []string) (r bool) {
	content = ut.SliceToLowerCase(content)
	filter = ut.SliceToLowerCase(filter)
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
	ut.appliedFilterTraceMessage("equal slices", content, filter, r)
	return
}

func (ut Util) SlicesNotEqual(content, filter []string) (r bool) {
	r = !ut.SlicesEqual(content, filter)
	ut.appliedFilterTraceMessage("not equal slices", content, filter, r)
	return
}

func (ut Util) SliceContainsSlice(content, filter []string) (r bool) {
	content = ut.SliceToLowerCase(content)
	filter = ut.SliceToLowerCase(filter)
	r = true
	for _, fil := range filter {
		if !ut.SliceContainsString(content, fil) {
			r = false
			break
		}
	}
	ut.appliedFilterTraceMessage("contains slice", content, filter, r)
	return
}

func (ut Util) SliceNotContainsSlice(content, filter []string) (r bool) {
	content = ut.SliceToLowerCase(content)
	filter = ut.SliceToLowerCase(filter)
	r = true
	if len(content) >= len(filter) {
		for _, fil := range filter {
			if ut.SliceContainsString(content, fil) {
				r = false
				break
			}
		}
	}
	ut.appliedFilterTraceMessage("not contains slice", content, filter, r)
	return
}

func (ut Util) RxMatchSliceAll(content, filter []string) (r bool) {
	content = ut.SliceToLowerCase(content)
	r = true
	for _, fil := range filter {
		for _, con := range content {
			if !ut.RxMatch(fil, con) {
				r = false
				break
			}
		}
	}
	ut.appliedFilterTraceMessage("rx match slice completely", content, filter, r)
	return
}

func (ut Util) RxMatchSliceOnce(content, filter []string) (r bool) {
	content = ut.SliceToLowerCase(content)
	r = true
	for _, fil := range filter {
		if !ut.RxContains(content, fil) {
			r = false
			break
		}
	}
	ut.appliedFilterTraceMessage("rx match slice once", content, filter, r)
	return
}

func (ut Util) SliceContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (ut Util) RxContains(slice []string, item string) bool {
	for _, s := range slice {
		if ut.RxMatch(item, s) {
			return true
		}
	}
	return false
}

func (ut Util) SliceToLowerCase(arr []string) (newArr []string) {
	for _, el := range arr {
		newArr = append(newArr, strings.ToLower(el))
	}
	return
}

func (ut Util) appliedFilterTraceMessage(name string, content, filter []string, r bool) {
	ut.Lg.Trace("applied filter: "+name,
		logseal.F{"content": content, "filter": filter, "result": r},
	)
}

// keep sorted end
