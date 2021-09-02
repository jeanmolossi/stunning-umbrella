package core

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"strconv"
)

func NewDbConnection() (*gorm.DB, error) {
	dsn := "root:root@tcp(up_db:3306)/logger_base"
	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	shouldEnableLogs, err := strconv.ParseBool(os.Getenv("ENABLE_DB_LOGS"))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Log{})
	db.LogMode(shouldEnableLogs)
	if err != nil {
		return nil, err
	}

	shouldCleanDb, err := strconv.ParseBool(os.Getenv("SHOULD_CLEAN_DB"))
	if shouldCleanDb && err == nil {
		db.Exec("DELETE FROM logs")
	}

	return db, nil
}
