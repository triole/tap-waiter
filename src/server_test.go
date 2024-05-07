package main

import (
	"testing"
)

func TestParseFilter(t *testing.T) {
	validateParseFilter(
		"front_matter.title==this+is+a+title",
		tIDXParamsFilter{
			Prefix: "front_matter.title", Operator: "==", Suffix: "this is a title",
		}, t,
	)
	validateParseFilter(
		"front_matter.title!=not+the+title",
		tIDXParamsFilter{
			Prefix: "front_matter.title", Operator: "!=", Suffix: "not the title",
		}, t,
	)
	validateParseFilter(
		"front_matter.tags*=title2",
		tIDXParamsFilter{
			Prefix: "front_matter.tags", Operator: "*=", Suffix: "title2",
		}, t,
	)
	validateParseFilter(
		"front_matter.tags!*=title",
		tIDXParamsFilter{
			Prefix: "front_matter.tags", Operator: "!*=", Suffix: "title",
		}, t,
	)
	validateParseFilter(
		"front_matter.title*=title[0-9]{1,2}",
		tIDXParamsFilter{
			Prefix: "front_matter.title", Operator: "*=", Suffix: "title[0-9]{1,2}",
		}, t,
	)
}

func validateParseFilter(filter string, exp tIDXParamsFilter, t *testing.T) {
	res := parseFilterString(filter)
	if res.Prefix != exp.Prefix || res.Operator != exp.Operator || res.Suffix != exp.Suffix {
		t.Errorf("parse filter failed, \nexp: %+v\ngot: %+v", exp, res)
	}
}
