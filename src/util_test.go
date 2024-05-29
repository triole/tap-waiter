package main

import (
	"path/filepath"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

var (
	testFolder = "../testdata"
)

func fromTestFolder(s string) (r string) {
	t, err := filepath.Abs(testFolder)
	if err == nil {
		r = filepath.Join(t, s)
	}
	return
}

func TestGetFileSize(t *testing.T) {
	validateGetFileSize(fromTestFolder("dump/yaml/1.yaml"), 1009, t)
}

func validateGetFileSize(fil string, exp uint64, t *testing.T) {
	res := getFileSize(fil)
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
	res := rxFind(rx, str)
	if res != exp {
		t.Errorf(
			"error rx find, rx: %s, str: %s, exp: %s, got: %s",
			rx, str, exp, res,
		)
	}
}

func validateTestRxMatch(rx, str string, exp bool, t *testing.T) {
	res := rxMatch(rx, str)
	if res != exp {
		t.Errorf(
			"error rx match, rx: %s, str: %s, exp: %v, got: %v",
			rx, str, exp, res,
		)
	}
}

func readYAMLFile(filepath string) (r map[string]interface{}) {
	by, _, err := readFile(filepath)
	if err != nil {
		return
	} else {
		_ = yaml.Unmarshal(by, &r)
	}
	return
}

func itfArrTostrArr(itf []interface{}) (r []string) {
	for _, el := range itf {
		r = append(r, el.(string))
	}
	return
}
