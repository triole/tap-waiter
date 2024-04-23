package main

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func readDataFile(filename string, basePath string, chin chan string, chout chan tJoinerEntry) {
	chin <- filename
	by := readFile(filename)

	je := tJoinerEntry{
		Path:     strings.TrimPrefix(strings.TrimPrefix(filename, basePath), "/"),
		FileMeta: getFileMeta(filename),
		Data:     readYaml(by),
	}
	chout <- je
	<-chin
}

func readYaml(yamlBytes []byte) (data map[string]interface{}) {
	err := yaml.Unmarshal(yamlBytes, &data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}
