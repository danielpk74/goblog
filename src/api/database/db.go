package database

import (
	"config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Connect open a DB connection
func Connect() (*gorm.DB, error) {
	config.Load()
	db, err := gorm.Open(config.DB_DRIVE, config.DB_URL)

	if err != nil {
		return nil, err
	}

	return db, nil
}
