package main

import (
	"encoding/json"
	"os"

	"github.com/triole/logseal"
)

func writeJSON(je tJoinerIndex, outFile string) {
	lg.Info("write joined data to", logseal.F{"path": outFile})
	jsonData, err := json.Marshal(je)
	if err != nil {
		lg.IfErrError("can not marshal final data", logseal.F{"content": je})
	} else {
		err = os.WriteFile(outFile, jsonData, 0644)
		lg.IfErrError(err, "can not write file", logseal.F{"path": outFile})
	}
}
