package controller

import (
	"os"
)

var (
	APIEnvironment string
	APIName        string
	APIPort        string
	APIVersion     string
)

func LoadAPIConfig() {
	APIEnvironment = os.Getenv("ENVIRONMENT")
	APIName = os.Getenv("API_NAME")
	APIPort = os.Getenv("API_PORT")
	APIVersion = os.Getenv("API_VERSION")
}
