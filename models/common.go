package models

import (
	"gorm.io/gorm"
)

// CreateAllTablesIfNotExists : 테이블이 없을 시 생성(테스트 시에만 사용)
func CreateAllTablesIfNotExists(db *gorm.DB) {
	if db != nil {
		db.AutoMigrate(&Locker{})
		db.AutoMigrate(&User{})
	}
}

// DeleteAllTables : 테이블의 모든 행 삭제(테스트 시에만 사용)
func DeleteAllTables(db *gorm.DB) {
	if db != nil {
		db.Where("1 = 1").Delete(&Locker{})
		db.Where("1 = 1").Delete(&User{})
	}
}
