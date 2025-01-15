package indexer

import (
	"strings"
)

func (ti TapIndex) getContentVal(key interface{}, content FileContent) (s []string) {
	keypath := ti.splitKey(key)
	if len(keypath) > 0 {
		switch keypath[0] {
		case "front_matter":
			s = ti.getMapVal(keypath[1:], content.FrontMatter)
		case "body":
			s = ti.getMapVal(keypath[1:], content.Body)
		default:
			s = ti.getMapVal(keypath, content.Body)
		}
	}
	return
}

func (ti TapIndex) getMapVal(keypath []string, itf interface{}) (s []string) {
	switch val := itf.(type) {
	case map[string]interface{}:
		val = ti.keyToLower(val)
		if len(keypath) > 0 {
			if dictval, ok := val[keypath[0]]; ok {
				s = ti.getMapVal(keypath[1:], dictval)
			}
		}
	case []interface{}:
		for _, x := range val {
			switch val2 := x.(type) {
			case map[string]interface{}:
				t := ti.getMapVal(keypath, val2)
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

func (ti TapIndex) keyToLower(dict map[string]interface{}) (r map[string]interface{}) {
	r = make(map[string]interface{})
	for key, val := range dict {
		r[strings.ToLower(key)] = val
	}
	return
}

func (ti TapIndex) splitKey(key interface{}) (arr []string) {
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
