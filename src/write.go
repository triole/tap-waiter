package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/triole/logseal"
)

func mkdir(foldername string) {
	_, err := os.Stat(foldername)
	if os.IsNotExist(err) {
		log.Fatal("folder does, create if", logseal.F{"path": foldername})
		os.MkdirAll(foldername, os.ModePerm)
	}
}

func writeJSON(je tJoinerIndex, outFile string) {
	lg.Info("write joined data to", logseal.F{"path": outFile})
	mkdir(filepath.Dir(outFile))
	jsonData, err := json.Marshal(je)
	if err != nil {
		lg.IfErrError("can not marshal final data", logseal.F{"content": je})
	} else {
		err = os.WriteFile(outFile, jsonData, 0644)
		lg.IfErrError(err, "can not write file", logseal.F{"path": outFile})
	}
}
