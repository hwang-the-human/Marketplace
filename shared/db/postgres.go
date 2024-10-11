package db

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(dsn string) (Database, error) {
	if dsn == "" {
		return nil, errors.New("DSN is empty")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	logrus.Info("Successfully connected to PostgreSQL DB")

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) GetDB() *gorm.DB {
	return p.db
}

func (p *PostgresDB) CloseDB() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgresDB) Migrate(dst ...interface{}) error {
	if err := p.db.AutoMigrate(dst...); err != nil {
		return err
	}
	logrus.Info("Successfully migrated PostgreSQL DB")
	return nil
}
