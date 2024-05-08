package main

import (
	"path/filepath"
	"testing"
)

var (
	testFolder = "../testdata"
)

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

func fromTestFolder(s string) (r string) {
	t, err := filepath.Abs(testFolder)
	if err == nil {
		r = filepath.Join(t, s)
	}
	return
}
