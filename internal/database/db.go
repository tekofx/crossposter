package database

import (
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"gorm.io/gorm"
)

var Database *gorm.DB

func InitializeDb() {
	var err error
	Database, err = gorm.Open(sqlite.Open("./data/database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	Database.AutoMigrate(&model.Post{}, &model.Image{})
}

func DropDatabase() {

	if Database != nil {
		sqlDB, err := Database.DB()
		if err == nil {
			sqlDB.Close()
		}
		Database = nil
	}

	// Remove the database file
	err := os.Remove("./database.db")
	if err != nil && !os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("Error removing database: %d", err))
	}
}
