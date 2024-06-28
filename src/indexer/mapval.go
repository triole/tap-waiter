package indexer

import (
	"strings"
)

func (ji JoinerIndex) getContentVal(key interface{}, content FileContent) (s []string) {
	keypath := ji.splitKey(key)
	if len(keypath) > 0 {
		if keypath[0] == "front_matter" {
			s = ji.getMapVal(keypath[1:], content.FrontMatter)
		} else if keypath[0] == "body" {
			s = ji.getMapVal(keypath[1:], content.Body)
		} else {
			s = ji.getMapVal(keypath, content.Body)
		}
	}
	return
}

func (ji JoinerIndex) getMapVal(keypath []string, itf interface{}) (s []string) {
	switch val := itf.(type) {
	case map[string]interface{}:
		val = ji.keyToLower(val)
		if len(keypath) > 0 {
			if dictval, ok := val[keypath[0]]; ok {
				s = ji.getMapVal(keypath[1:], dictval)
			}
		}
	case []interface{}:
		for _, x := range val {
			switch val2 := x.(type) {
			case map[string]interface{}:
				t := ji.getMapVal(keypath, val2)
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

func (ji JoinerIndex) keyToLower(dict map[string]interface{}) (r map[string]interface{}) {
	r = make(map[string]interface{})
	for key, val := range dict {
		r[strings.ToLower(key)] = val
	}
	return
}

func (ji JoinerIndex) splitKey(key interface{}) (arr []string) {
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
