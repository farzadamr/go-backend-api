package database

import (
	"log"

	"github.com/farzadamr/go-backend-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PreloadEntity struct {
	Entity string
}

var dbClient *gorm.DB

func InitDb(cfg *config.Config) error {
	var err error
	cnn := cfg.Database.DSN()
	dbClient, err = gorm.Open(postgres.Open(cnn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := dbClient.DB()
	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	log.Println("Connected to Postgres database successfully")
	return nil
}

func GetDb() *gorm.DB {
	return dbClient
}

func CloseDb() {
	conn, _ := dbClient.DB()
	conn.Close()
}

// Preload
func Preload(db *gorm.DB, preloads []PreloadEntity) *gorm.DB {
	for _, item := range preloads {
		db = db.Preload(item.Entity)
	}
	return db
}
