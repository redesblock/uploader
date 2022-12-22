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
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxLifetime(time.Second * 300)

	db.AutoMigrate(&WatchFile{})
	db.AutoMigrate(&UploadFile{})
	db.AutoMigrate(&Voucher{})
	if err := db.Save(&Voucher{
		ID:      1,
		Node:    "183.131.181.164",
		Voucher: "dd812517f2ecfe75d7b08e908a857c8703477949770fbb76f2244d0d90cb4a12",
		Usable:  true,
	}).Error; err != nil {
		return nil, err
	}
	return db, nil
}
