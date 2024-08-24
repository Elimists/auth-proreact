package config

import (
	"os"
)

var (
	databaseHost     string
	databasePort     string
	databaseUser     string
	databasePassword string
	databaseName     string
)

func LoadDatabase() {
	databaseHost = os.Getenv("DB_HOST")
	databasePort = os.Getenv("DB_PORT")
	databaseUser = os.Getenv("DB_USER")
	databasePassword = os.Getenv("DB_PASSWORD")
	databaseName = os.Getenv("DB_DATABASE")
}

func GetDatabaseHost() string {
	return databaseHost
}

func GetDatabasePort() string {
	return databasePort
}

func GetDatabaseUser() string {
	return databaseUser
}

func GetDatabasePassword() string {
	return databasePassword
}

func GetDatabaseName() string {
	return databaseName
}
