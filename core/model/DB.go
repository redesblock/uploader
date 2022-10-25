package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

func New(mode string, dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	var conn gorm.Dialector
	switch mode {
	case "sqlite":
		os.MkdirAll(filepath.Dir(dsn), 0755)
		conn = sqlite.Open(dsn)
	case "mysql":
		conn = mysql.Open(dsn)
	case "postgres":
		conn = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("failed to connect to database: invalid db engine, supported types: sqlite, mysql, postgres")
	}

	db, err := gorm.Open(conn, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)                   //最大空闲连接数
	sqlDB.SetMaxOpenConns(30)                   //最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * 300) //设置连接空闲超时

	db.AutoMigrate(&WatchFile{})
	db.AutoMigrate(&UploadFile{})
	db.AutoMigrate(&Voucher{})
	if err := db.Save(&Voucher{
		ID:      1,
		Host:    "http://58.34.1.130:1633",
		Voucher: "e92110b77f959065768e24a44c5ab04de4f6bc20f0010fbba726ee4b31291797",
		Usable:  true,
	}).Error; err != nil {
		return nil, err
	}
	return db, nil
}
