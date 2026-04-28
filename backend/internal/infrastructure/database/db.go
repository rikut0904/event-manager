package database

import (
	"log"
	"os"
	"time"

	"backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("WARNING: DATABASE_URL is not set. Skipping DB connection.")
		return nil
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// 自動マイグレーション（テーブル作成・更新）
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("error during migration: %v", err)
	}

	// コネクションプールの設定
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db
}
