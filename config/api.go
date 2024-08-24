package config

import (
	"os"
)

var (
	apiEnvironment string
	apiName        string
	apiPort        string
	apiVersion     string
)

func LoadAPI() {
	apiEnvironment = os.Getenv("ENVIRONMENT")
	apiName = os.Getenv("API_NAME")
	apiPort = os.Getenv("API_PORT")
	apiVersion = os.Getenv("API_VERSION")
}

func GetAPIEnvironment() string {
	return apiEnvironment
}

func GetAPIName() string {
	return apiName
}

func GetAPIPort() string {
	return apiPort
}

func GetAPIVersion() string {
	return apiVersion
}
