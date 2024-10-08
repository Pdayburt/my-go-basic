package dao

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	if db == nil {

		zap.L().Error(" InitTable database connection is nil", zap.Error(errors.New("传入的db *gorm.DB为空")))
		return errors.New("database connection is nil")
	}
	return db.AutoMigrate(&Interactive{})
}
