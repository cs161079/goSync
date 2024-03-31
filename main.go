package main

import (
	"fmt"
	"os"

	oasaSyncWeb "github.com/cs161079/goSync/Web"
	logger "github.com/cs161079/godbLib/Utils/goLogger"
	"github.com/joho/godotenv"
)

func initEnviroment() {
	// loads values from .env into the system
	if err := godotenv.Load("enviroment.env"); err != nil {
		logger.ERROR("No .env file found")
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func main() {
	logger.InitLogger("goSyncApplication")
	initEnviroment()
	/*
		Test with function that makes request and test it in for loop
	*/
	//****************************************************************
	// for i := 0; i < 1000; i++ {
	// 	_, err := oasaSyncApi.GetBusLinesTest()
	// 	if err != nil {
	// 		logger.ERROR(err.Error())
	// 	} else {
	// 		logger.INFO("Response Succesfully.")
	// 	}
	// }
	//****************************************************************
	logger.INFO(getEnv("LOGS_PATH", "/oasaLogs"))
	responseStr, err := oasaSyncWeb.MakeRequest("getLines")
	if err != nil {
		logger.ERROR(err.Error())
	} else {
		fmt.Print(responseStr)
	}

}
