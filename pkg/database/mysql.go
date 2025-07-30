package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQLClient() (*gorm.DB, error) {
	var (
		DB_User   = os.Getenv("DB_USER")
		DB_Pass   = os.Getenv("DB_PASS")
		DB_Host   = os.Getenv("DB_HOST")
		DB_Port   = os.Getenv("DB_PORT")
		DB_DbName = os.Getenv("DB_NAME")
	)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DB_User, DB_Pass, DB_Host, DB_Port, DB_DbName,
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(fmt.Sprintf("failed to connect database : %s", err))
	}

	sqlDb, errorDb := db.DB()

	if errorDb != nil {
		panic("Issue on db connection")
	}

	sqlDb.SetMaxIdleConns(100)
	sqlDb.SetMaxOpenConns(500)
	sqlDb.SetConnMaxIdleTime(time.Minute)

	sqlDb.SetConnMaxLifetime(time.Minute * 60)

	err = sqlDb.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
