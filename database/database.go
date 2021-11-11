package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	"github.com/sithumonline/quick-note/config"
)

func Database() *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.GetEnv("database.URL")), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database")
	}
	log.Info("Database connected")
	return db
}
