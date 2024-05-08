package main

import "testing"

func TestEqualSlices(t *testing.T) {
	validateEqualSlices(
		[]string{"tag1", "tag2"},
		[]string{"tag1", "tag2"},
		true, t,
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
