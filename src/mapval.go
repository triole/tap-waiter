package main

import (
	"strings"
)

func getContentVal(key interface{}, content tContent) (s []string) {
	keypath := splitKey(key)
	if len(keypath) > 0 {
		if keypath[0] == "front_matter" {
			s = getMapVal(keypath[1:], content.FrontMatter)
		} else if keypath[0] == "body" {
			s = getMapVal(keypath[1:], content.Body)
		} else {
			s = getMapVal(keypath, content.Body)
		}
	}
	return
}

func getMapVal(keypath []string, itf interface{}) (s []string) {
	switch val := itf.(type) {
	case map[string]interface{}:
		val = keyToLower(val)
		if len(keypath) > 0 {
			if dictval, ok := val[keypath[0]]; ok {
				s = getMapVal(keypath[1:], dictval)
			}
		}
	case []interface{}:
		for _, x := range val {
			switch val2 := x.(type) {
			case map[string]interface{}:
				t := getMapVal(keypath, val2)
				if len(t) > 0 {
					s = t
				}
			default:
				s = append(s, x.(string))
			}
		}
	default:
		if val != nil {
			s = []string{val.(string)}
		} else {
			s = []string{}
		}
	}
	return
}

func keyToLower(dict map[string]interface{}) (r map[string]interface{}) {
	r = make(map[string]interface{})
	for key, val := range dict {
		r[strings.ToLower(key)] = val
	}
	return
}

func splitKey(key interface{}) (arr []string) {
	var tmp []string
	switch val := key.(type) {
	case string:
		tmp = strings.Split(val, ".")
	case []string:
		tmp = val
	}
	for _, el := range tmp {
		arr = append(arr, strings.ToLower(el))
	}
	return
}
