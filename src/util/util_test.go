package util

import (
	"testing"

	"github.com/triole/logseal"
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

func TestIsEmpty(t *testing.T) {
	validateIsEmpty("hello world", false, t)
	validateIsEmpty("___", false, t)
	validateIsEmpty(" ", true, t)
	validateIsEmpty("   ", true, t)
}

func validateIsEmpty(str string, exp bool, t *testing.T) {
	ut := Init(lg)
	res := ut.IsEmpty(str)
	if res != exp {
		t.Errorf(
			"error is empty: str: %s, exp: %v, got: %v",
			str, exp, res,
		)
	}
}

func TestIsPathSafe(t *testing.T) {
	ut := Init(logseal.Init())

	// Test cases for path traversal protection
	tests := []struct {
		name     string
		path     string
		baseDir  string
		expected bool
	}{
		{
			name:     "valid path within base",
			path:     "/allowed/path/file.txt",
			baseDir:  "/allowed/path",
			expected: true,
		},
		{
			name:     "path equals base directory",
			path:     "/allowed/path",
			baseDir:  "/allowed/path",
			expected: true,
		},
		{
			name:     "path traversal with ..",
			path:     "/allowed/path/../secret/file.txt",
			baseDir:  "/allowed/path",
			expected: false,
		},
		{
			name:     "path traversal escaping base",
			path:     "/allowed/some/../../../etc/passwd",
			baseDir:  "/allowed/some",
			expected: false,
		},
		{
			name:     "path traversal with absolute ..",
			path:     "/allowed/../../etc/passwd",
			baseDir:  "/allowed",
			expected: false,
		},
		{
			name:     "path traversal to sibling directory",
			path:     "/allowed/dir1/../dir2/file.txt",
			baseDir:  "/allowed/dir1",
			expected: false,
		},
		{
			name:     "nested path within base",
			path:     "/allowed/path/subdir/deep/file.txt",
			baseDir:  "/allowed/path",
			expected: true,
		},
		{
			name:     "similar prefix attack",
			path:     "/allowed-fake/path/file.txt",
			baseDir:  "/allowed",
			expected: false,
		},
		{
			name:     "path without separator after base",
			path:     "/allowedextra/file.txt",
			baseDir:  "/allowed",
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			baseDir:  "/allowed/path",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ut.IsPathSafe(tt.path, tt.baseDir)
			if result != tt.expected {
				t.Errorf(
					"IsPathSafe(%q, %q) = %v, expected %v",
					tt.path, tt.baseDir, result, tt.expected,
				)
			}
		})
	}
}
