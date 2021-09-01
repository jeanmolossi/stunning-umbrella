package core

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDbConnection() (*gorm.DB, error) {
	dsn := "root:root@tcp(up_db:3306)/logger_base"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Log{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
