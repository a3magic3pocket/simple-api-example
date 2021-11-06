package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB: DB
var DB *gorm.DB

// SetDB : DB Connection 획득
func SetDB() {
	var err error

	var (
		_, b, _, _     = runtime.Caller(0)
		workingDirPath = filepath.Dir(filepath.Dir(b))
	)
	pvDirPath := fmt.Sprintf("%s/%s", workingDirPath, "sqlite3")
	os.MkdirAll(pvDirPath, os.ModePerm)

	ginMode := os.Getenv("GIN_MODE")

	switch strings.ToLower(ginMode) {
	case "debug":
		DB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s/debug.db", pvDirPath)), &gorm.Config{})
	case "release":
		DB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s/operation.db", pvDirPath)), &gorm.Config{})
	default:
		log.Fatalf("gin mode is not allowed. gin mode : %s\n", ginMode)
	}
	if err != nil {
		log.Fatalf("failed to get DB connection. error : %s\n", err.Error())
	}
}
