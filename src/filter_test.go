package main

import "testing"

func TestEqualSlices(t *testing.T) {
	validateEqualSlices(
		[]string{"tag1", "tag2"},
		[]string{"tag1", "tag2"},
		true, t,
	)
	validateEqualSlices(
		[]string{"tag2", "tag1"},
		[]string{"tag1", "tag2"},
		true, t,
	)
	validateEqualSlices(
		[]string{"tag1", "tag2"},
		[]string{"tag2", "tag4"},
		false, t,
	)

}

func validateEqualSlices(s1, s2 []string, exp bool, t *testing.T) {
	res := equalSlices(s1, s2)
	if exp != res {
		t.Errorf(
			"error equal slices, slices: %+v %+v, exp: %v, got: %v",
			s1, s2, exp, res,
		)
	}
}

func TestContainsSlice(t *testing.T) {
	validateContainsSlice(
		[]string{"tag1", "tag2", "tag3", "tag4", "tag5"},
		[]string{"tag3"},
		true, t,
	)
	validateContainsSlice(
		[]string{"tag1", "tag2", "tag3", "tag4", "tag5"},
		[]string{"tag3", "tag5"},
		true, t,
	)
	validateContainsSlice(
		[]string{"tag1", "tag2", "tag3"},
		[]string{"tag5"},
		false, t,
	)
	validateContainsSlice(
		[]string{"tag1", "tag2", "tag3", "tag4", "tag5"},
		[]string{"tag7", "tag5"},
		false, t,
	)
}

func validateContainsSlice(s1, s2 []string, exp bool, t *testing.T) {
	res := containsSlice(s1, s2)
	if exp != res {
		t.Errorf(
			"error contains slice, slices: %+v %+v, exp: %v, got: %v",
			s1, s2, exp, res,
		)
	}
}

func TestNotContainsSlice(t *testing.T) {
	validateNotContainsSlice(
		[]string{"tag1", "tag2", "tag3"},
		[]string{"tag9"},
		true, t,
	)
	validateNotContainsSlice(
		[]string{"tag1", "tag2", "tag3"},
		[]string{"tag1", "tag2", "tag3", "tag4"},
		true, t,
	)
	validateNotContainsSlice(
		[]string{"tag1", "tag2", "tag3"},
		[]string{"tag1"},
		false, t,
	)
	validateNotContainsSlice(
		[]string{"tag1", "tag2", "tag3"},
		[]string{"tag1", "tag2"},
		false, t,
	)
}

func validateNotContainsSlice(s1, s2 []string, exp bool, t *testing.T) {
	res := notContainsSlice(s1, s2)
	if exp != res {
		t.Errorf(
			"error not contains slice, slices: %+v %+v, exp: %v, got: %v",
			s1, s2, exp, res,
		)
	}
}
