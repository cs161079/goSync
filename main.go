package main

import (
	"fmt"

	oasaSyncApi "github.com/cs161079/goSync/Api"
	logger "github.com/cs161079/godbLib/Utils/goLogger"
)

func main() {
	logger.InitLogger("goSyncApplication")
	busLines, err := oasaSyncApi.GetBusLinesTest()
	if err != nil {
		logger.ERROR(err.Error())
	} else {
		fmt.Print(busLines)
	}
}
