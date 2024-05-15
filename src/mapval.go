package main

import (
	"strings"

	"github.com/triole/logseal"
)

func getMapVal(key interface{}, dict map[string]interface{}) (s []string) {
	keypath := splitKey(key)
	dict = keyToLower(dict)
	if len(keypath) > 0 {
		if dictval, ok := dict[keypath[0]]; ok {
			switch val1 := dictval.(type) {
			case map[string]interface{}:
				s = getMapVal(keypath[1:], val1)
			case []interface{}:
				for _, x := range val1 {
					switch val2 := x.(type) {
					case map[string]interface{}:
						t := getMapVal(keypath[1:], val2)
						if len(t) > 0 {
							s = t
						}
					default:
						s = append(s, x.(string))
					}
				}
			default:
				s = []string{val1.(string)}
			}
		}
	} else {
		lg.Warn("can not parse given map key", logseal.F{"key": key})
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
