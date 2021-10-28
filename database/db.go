package database

import (
	"log"
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB: DB
var DB *gorm.DB

// SetDB : DB Connection 획득
func SetDB() {
	var err error
	ginMode := os.Getenv("GIN_MODE")
	switch strings.ToLower(ginMode) {
	case "debug":
		DB, err = gorm.Open(sqlite.Open("debug.db"), &gorm.Config{})
	case "operation":
		DB, err = gorm.Open(sqlite.Open("operation.db"), &gorm.Config{})
	default:
		log.Fatalf("gin mode is not allowed. gin mode : %s\n", ginMode)
	}
	if err != nil {
		log.Fatalf("failed to get DB connection. error : %s\n", err.Error())
	}
}
