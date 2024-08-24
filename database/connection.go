package database

import (
	"fmt"

	"github.com/Elimists/go-app/config"
	"github.com/Elimists/go-app/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error
	var dbConn *gorm.DB

	environment := config.GetAPIEnvironment()
	if environment == "production" {
		// MySQL connection string
		user := config.GetDatabaseUser()
		password := config.GetDatabasePassword()
		database := config.GetDatabaseName()
		stringConn := fmt.Sprintf("%s:%s@/%s?parseTime=true", user, password, database)

		dbConn, err = gorm.Open(mysql.Open(stringConn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			panic("Could not connect to MySQL database: " + err.Error())
		}
	} else {
		// SQLite connection string
		sqliteFile := "local.db"

		dbConn, err = gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			panic("Could not connect to SQLite database: " + err.Error())
		}
	}

	DB = dbConn

	// Generate tables using the model if they don't exist.
	err = DB.AutoMigrate(
		&models.User{},
		&models.UserVerification{},
	)
	if err != nil {
		panic("Could not auto-migrate database: " + err.Error())
	}
}
