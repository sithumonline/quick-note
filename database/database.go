package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

func Database() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database")
	}
	log.Info("Database connected")
	return db
}
