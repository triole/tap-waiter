package util

import (
	"testing"
)

func TestGetFileSize(t *testing.T) {
	ut := Init(lg)
	validateGetFileSize(ut.FromTestFolder("dump/yaml/1.yaml"), 1009, t)
}

func validateGetFileSize(fil string, exp uint64, t *testing.T) {
	ut := Init(lg)
	res := ut.GetFileSize(fil)
	if res != exp {
		t.Errorf(
			"error get file size, file: %s, exp: %d, got: %d", fil, exp, res,
		)
	}
}

func TestRegex(t *testing.T) {
	validateTestRxFind("^[helo]+", "hello world", "hello", t)
	validateTestRxFind("lo.+r", "hello world", "lo wor", t)
	validateTestRxFind("lo.+r", "hello world", "lo wor", t)
	validateTestRxFind("[^w]+$", "hello world", "orld", t)

	validateTestRxMatch("^[helo]+", "hello world", true, t)
	validateTestRxMatch("^[^helo]+", "hello world", false, t)
	validateTestRxMatch("world", "hello world", true, t)
	validateTestRxMatch("mars", "hello world", false, t)
}

func validateTestRxFind(rx, str, exp string, t *testing.T) {
	ut := Init(lg)
	res := ut.RxFind(rx, str)
	if res != exp {
		t.Errorf(
			"error rx find, rx: %s, str: %s, exp: %s, got: %s",
			rx, str, exp, res,
		)
	}
}

func validateTestRxMatch(rx, str string, exp bool, t *testing.T) {
	ut := Init(lg)
	res := ut.RxMatch(rx, str)
	if res != exp {
		t.Errorf(
			"error rx match, rx: %s, str: %s, exp: %v, got: %v",
			rx, str, exp, res,
		)
	}
}
