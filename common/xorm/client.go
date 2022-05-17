package xorm

import (
	"github.com/showurl/Zero-IM-Server/common/xorm/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

func GetClient(
	cfg global.MysqlConfig,
) *gorm.DB {
	var db *gorm.DB
	var err error
	var logLevel logger.LogLevel
	level := strings.ToUpper(cfg.LogLevel)
	if level == "" || level == "INFO" {
		logLevel = logger.Info
	} else {
		if level == "SILENT" {
			logLevel = logger.Silent
		}
		if level == "WARNING" {
			logLevel = logger.Warn
		}
		if level == "ERROR" {
			logLevel = logger.Error
		}
	}
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       cfg.Addr,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         191,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		panic(err)
	}
	return db
}
