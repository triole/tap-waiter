package main

import (
	"testing"
)

func TestParseFilter(t *testing.T) {
	validateParseFilter(
		"front_matter.title==this+is+a+title",
		tIDXParamsFilter{
			Prefix: "front_matter.title", Operator: "==", Suffix: []string{"this+is+a+title"},
		}, t,
	)
	validateParseFilter(
		"front_matter.title!=not+the+title",
		tIDXParamsFilter{
			Prefix: "front_matter.title", Operator: "!=", Suffix: []string{"not+the+title"},
		}, t,
	)
	validateParseFilter(
		"front_matter.tags=*title2",
		tIDXParamsFilter{
			Prefix: "front_matter.tags", Operator: "=*", Suffix: []string{"title2"},
		}, t,
	)
	validateParseFilter(
		"front_matter.tags!=*title",
		tIDXParamsFilter{
			Prefix: "front_matter.tags", Operator: "!=*", Suffix: []string{"title"},
		}, t,
	)
	validateParseFilter(
		"front_matter.tags!==tag1,tag2",
		tIDXParamsFilter{
			Prefix: "front_matter.tags", Operator: "!==", Suffix: []string{"tag1", "tag2"},
		}, t,
	)
	validateParseFilter(
		"front_matter.title=~title[0-9]{1,2}",
		tIDXParamsFilter{
			Prefix: "front_matter.title", Operator: "=~", Suffix: []string{"title[0-9]{1,2}"},
		}, t,
	)
}

func validateParseFilter(filter string, exp tIDXParamsFilter, t *testing.T) {
	res := parseFilterString(filter)
	r := true
	if res.Prefix != exp.Prefix || res.Operator != exp.Operator {
		r = false
		for i := 0; i < len(exp.Suffix); i++ {
			if res.Suffix[i] != exp.Suffix[i] {
				r = false
			}
		}
	}
	if !r {
		t.Errorf("parse filter failed, \nexp: %+v\ngot: %+v", exp, res)
	}
}
