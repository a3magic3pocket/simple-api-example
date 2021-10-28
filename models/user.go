package models

import (
	"simple-api-example/database"
)

// User : 유저
type User struct {
	ID        int    `json:"ID" gorm:"column:ID;primaryKey;autoIncrement;notNull;"`
	Name      string `json:"Name" gorm:"column:Name;unique;"`
	SecretKey string `json:"SecretKey" gorm:"column:SecretKey;"`
	Group     string `json:"Group" gorm:"column:Group;"`
}

type Users []User

// TableName : 유저 테이블명
func (user User) TableName() string {

	return "USER"
}

// Create : 유저 생성
func (user *User) Create() error {
	result := database.DB.Create(user)

	return result.Error
}

// Create : 유저들 생성(테스트 시에만 사용)
func (users Users) Create() error {
	result := database.DB.Create(users)

	return result.Error
}

// Get : userName으로 유저 조회
func (user *User) Get(userName string) error {
	result := database.DB.Where(`Name = ?`, userName).Take(user)

	return result.Error
}

// GetLoginInfo : 로그인 정보 조회
func (user *User) GetLoginInfo(userName string, secretKey string) error {
	result := database.DB.
		Where(`Name = ? AND SecretKey = ?`, userName, secretKey).
		Take(user)

	return result.Error
}
