package logger

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Log struct {
	Db      *gorm.DB `gorm:"-"`
	Type    string   `gorm:"type:varchar(255)"`
	Message string   `gorm:"type:varchar(255)"`
	RefID   string   `gorm:"type:bigint;primary_key;unique"`
}

func NewLogger(db *gorm.DB) *Log {
	return &Log{
		Db: db,
	}
}

func (l *Log) AddLog() error {
	logCreated := l.Db.Create(l)
	if logCreated.Error != nil {
		errorMsg := fmt.Sprintf("Erro ao criar log: %s", logCreated.Error)
		return errors.New(errorMsg)
	}

	return nil
}
