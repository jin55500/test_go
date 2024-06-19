package database

import (
	"os"
	Model "test/go/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

var err error

func ConnectDb() {

	dsn := os.Getenv("PATH_DB")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// migrate
	Db.AutoMigrate(&Model.User{})
	Db.AutoMigrate(&Model.Accounting{})
}
