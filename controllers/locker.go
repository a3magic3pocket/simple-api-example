package controllers

import (
	"errors"
	"net/http"
	"simple-api-example/auth"
	"simple-api-example/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LockerOutput : Locker 아웃풋 구조체
type LockerOutput struct {
	ID       int    `json:"ID"`
	Location string `json:"Location"`
}

// UpdateLockersInput : UpdateLockers 인풋 구조체
type UpdateLockersInput struct {
	models.Locker
	UpdateIDs []int `json:"UpdateIDs"`
}

// CreateLockers : Lockers 생성
func CreateLockers(c *gin.Context) {
	lockers := models.Lockers{}
	if err := c.ShouldBindJSON(&lockers); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "data is invalid"})
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err = lockers.Create(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdateLocker : Locker 업데이트
func UpdateLocker(c *gin.Context) {
	locker := models.Locker{}
	if err := c.ShouldBindJSON(&locker); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "data is invalid"})
		return
	}

	// locker ID 유효성 검사
	if locker.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "ID is empty"})
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err = locker.PartialUpdate(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdateLockers : Lockers를 한 값으로 업데이트
func UpdateLockers(c *gin.Context) {
	updateLockersInput := UpdateLockersInput{}
	if err := c.ShouldBindJSON(&updateLockersInput); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "data is invalid"})
		return
	}

	// locker ID 유효성 검사
	nonZeroLockers := models.Lockers{}
	for _, id := range updateLockersInput.UpdateIDs {
		if id != 0 {
			nonZeroLockers = append(nonZeroLockers, models.Locker{ID: id})
		}
	}

	if len(nonZeroLockers) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "all of UpdateIDs are empty"})
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err = nonZeroLockers.PartialUpdate(userID, updateLockersInput.Locker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// DeleteLockers : Lockers 삭제(soft delete)
func DeleteLockers(c *gin.Context) {
	lockers := models.Lockers{}
	if err := c.ShouldBindJSON(&lockers); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "data is invalid"})
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// locker ID 유효성 검사
	for _, locker := range lockers {
		if locker.ID == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "ID is empty"})
			return
		}
	}

	err = lockers.DeleteLockers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// RetrieveLocker : Locker 조회
func RetrieveLocker(c *gin.Context) {
	rawLockerID := c.Param("id")
	lockerID, err := strconv.Atoi(rawLockerID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "locker id is invalid"})
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	locker := models.Locker{}
	err = locker.GetOwned(userID, lockerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "locker is not exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		}
		return
	}

	lockerOuput := LockerOutput{}
	lockerOuput.ID = locker.ID
	lockerOuput.Location = locker.Location

	c.JSON(http.StatusOK, gin.H{"data": lockerOuput})
}

// RetrieveLockers : Lockers 조회
func RetrieveLockers(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	lockers := models.Lockers{}
	err = lockers.GetOwned(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	lockerOuputs := []LockerOutput{}
	for _, locker := range lockers {
		lockerOuput := LockerOutput{}
		lockerOuput.ID = locker.ID
		lockerOuput.Location = locker.Location
		lockerOuputs = append(lockerOuputs, lockerOuput)
	}

	c.JSON(http.StatusOK, gin.H{"data": lockerOuputs})
}

// RetrieveAllLockers : 모든 Lockers 조회
func RetrieveAllLocker(c *gin.Context) {
	lockers := models.Lockers{}
	err := lockers.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	lockerOuputs := []LockerOutput{}
	for _, locker := range lockers {
		lockerOuput := LockerOutput{}
		lockerOuput.ID = locker.ID
		lockerOuput.Location = locker.Location
		lockerOuputs = append(lockerOuputs, lockerOuput)
	}

	c.JSON(http.StatusOK, gin.H{"data": lockerOuputs})
}
