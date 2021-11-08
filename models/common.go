package models

import (
	"gorm.io/gorm"
)

// CreateAllTablesIfNotExists : 테이블이 없을 시 생성
func CreateAllTablesIfNotExists(db *gorm.DB) {
	modelAddrs := []interface{}{
		&Locker{}, &User{},
	}

	if db != nil {
		for _, modelAddr := range modelAddrs {
			if !db.Migrator().HasTable(modelAddr) {
				db.Migrator().CreateTable(modelAddr)
			}
		}
	}
}

// DeleteAllTables : 테이블의 모든 행 삭제(테스트 시에만 사용)
func DeleteAllTables(db *gorm.DB) {
	if db != nil {
		db.Where("1 = 1").Delete(&Locker{})
		db.Where("1 = 1").Delete(&User{})
	}
}
