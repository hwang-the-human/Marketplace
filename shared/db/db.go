package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(dsn string, dst ...interface{}) {
	if dsn == "" {
		logrus.Fatalf("DSN is empty")
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logrus.Fatalf("Failed to connect to DB: %v", err)
	}

	logrus.Info("Successfully connected to DB")

	if err := db.AutoMigrate(dst...); err != nil {
		logrus.Fatalf("Error migrating the database: %v", err)
	}

	logrus.Info("Successfully migrated DB")
}

func GetDB() *gorm.DB {
	if db == nil {
		logrus.Fatal("Database connection has not been initialized.")
	}
	return db
}

func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("Failed to retrieve database instance: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		logrus.Fatalf("Failed to close DB connection: %v", err)
	}

	logrus.Info("DB connection closed.")
}
