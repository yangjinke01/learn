package middleware

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func DBConnection() *gorm.DB {
	if db == nil {
		db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	}
	return db
}
