package util

import (
	"sort"
	"strings"
)

// keep-sorted start block=yes newline_separated=yes
func (ut Util) RxSliceContainsSliceFully(content, filter []string) bool {
	content = ut.SliceToLowerCase(content)
	for _, fil := range filter {
		if !ut.RxSliceContainsString(content, fil) {
			return false
		}
	}
	return true
}

func (ut Util) RxSliceContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if ut.RxMatch(item, s) {
			return true
		}
	}
	return false
}

func (ut Util) RxSliceMatchesSliceFully(content, filter []string) (r bool) {
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
	return
}

func (ut Util) SliceToLowerCase(arr []string) (newArr []string) {
	for _, el := range arr {
		newArr = append(newArr, strings.ToLower(el))
	}
	return
}

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
	return
}

func (ut Util) SlicesNotEqual(content, filter []string) (r bool) {
	r = !ut.SlicesEqual(content, filter)
	return
}

//keep-sorted end
