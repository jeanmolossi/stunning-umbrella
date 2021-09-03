package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"strconv"
	"up-planilhas-go/core/logger"
)

func NewDbConnection() (*gorm.DB, error) {
	dsn := generateDsn()
	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	shouldEnableLogs, err := strconv.ParseBool(os.Getenv("ENABLE_DB_LOGS"))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&logger.Log{})
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

func generateDsn() string {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
}
