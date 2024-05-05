package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestMakeJoinerIndex(t *testing.T) {
	validateMakeJoinerIndex("created", false, t)
	validateMakeJoinerIndex("created", false, t)
	validateMakeJoinerIndex("lastmod", false, t)
	validateMakeJoinerIndex("lastmod", false, t)
	validateMakeJoinerIndex("size", true, t)
	validateMakeJoinerIndex("size", false, t)
}

func validateMakeJoinerIndex(sortBy string, ascending bool, t *testing.T) {
	validateArr := []interface{}{}
	idx := makeJoinerIndex(newTestParams(sortBy, ascending))
	for _, je := range idx {
		v := reflect.ValueOf(je)
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldName := strings.ToLower(v.Type().Field(i).Name)
			if fieldName == sortBy {
				validateArr = append(validateArr, field.Interface())
			}
		}
	}
	fail := false
	for i := 1; i <= len(validateArr)-1; i++ {
		lastVal := validateArr[i-1]
		thisVal := validateArr[i]
		switch thisVal.(type) {
		case string:
			c := strings.Compare(lastVal.(string), thisVal.(string))
			fmt.Printf("%+v\n", c)
		case uint64:
			if ascending && lastVal.(uint64) > thisVal.(uint64) {
				fail = true
			}
			if !ascending && lastVal.(uint64) < thisVal.(uint64) {
				fail = true
			}
		}
	}
	if fail {
		t.Errorf("sort failed, sortby: %s, asc: %v", sortBy, ascending)
	}
}

func newTestParams(sortBy string, ascending bool) (p tIDXParams) {
	p.Endpoint.Folder = "../testdata"
	p.Endpoint.ReturnValues.Content = true
	p.Endpoint.ReturnValues.Created = true
	p.Endpoint.ReturnValues.LastMod = true
	p.Endpoint.ReturnValues.Metadata = true
	p.Endpoint.ReturnValues.Size = true
	p.Threads = 8
	p.Ascending = ascending
	p.SortBy = sortBy
	return
}
