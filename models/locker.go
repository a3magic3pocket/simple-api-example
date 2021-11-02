package models

import (
	"simple-api-example/database"
)

type Locker struct {
	ID       int    `json:"ID" gorm:"column:ID;primaryKey;autoIncrement;notNull;"`
	Location string `json:"Location" gorm:"column:Location;"`
	Owner    int    `json:"Owner" gorm:"column:Owner;"`
	Status   string `json:"Status" gorm:"column:Status;default:normal"`
}

type Lockers []Locker

// TableName : Locker 테이블 명
func (Locker) TableName() string {

	return "LOCKER"
}

// CreateLockers : 조회 요청자 소유의 로커들 삽입
func (lockers *Lockers) Create(userID int) error {
	for i := range *lockers {
		(*lockers)[i].Owner = userID
	}
	result := database.DB.Model(lockers).
		Select("Location", "Owner").
		Create(lockers)

	return result.Error
}

// UpdateLockers : 조회 요청자 소유의 로커 수정
func (locker *Locker) PartialUpdate(userID int) error {
	result := database.DB.Model(locker).
		Where(`Owner = ?`, userID).
		Where(`Status = "normal"`).
		Select("Location").
		Updates(locker)

	return result.Error
}

// GetOwnedLocker : 조회 요청자 소유의 로커 조회
func (locker *Locker) GetOwned(userID int, lockerID int) error {
	result := database.DB.
		Where(`Owner = ?`, userID).
		Where(`ID = ?`, lockerID).
		Where(`Status = "normal"`).
		Take(locker)

	return result.Error
}

// GetOwnedLocker : 조회 요청자 소유의 로커들 조회
func (lockers *Lockers) GetOwned(userID int) error {
	result := database.DB.
		Where(`Owner = ?`, userID).
		Where(`Status = "normal"`).
		Find(lockers)

	return result.Error
}

// GetOwnedLocker : 모든 로커들 조회
func (lockers *Lockers) GetAll() error {
	result := database.DB.
		Find(lockers)

	return result.Error
}

// DeleteLocker : 모든 로커들 삭제(soft delete)
func (lockers *Lockers) DeleteLockers(userID int) error {
	result := database.DB.Model(lockers).
		Where(`Owner = ? `, userID).
		Where(`Status = "normal"`).
		Select("Status").
		Updates(Locker{Status: "deleted"})

	return result.Error
}
