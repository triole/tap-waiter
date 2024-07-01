package util

import "github.com/triole/logseal"

type Util struct {
	Lg logseal.Logseal
}

func Init(lg logseal.Logseal) (ut Util) {
	return Util{
		Lg: lg,
	}
}
