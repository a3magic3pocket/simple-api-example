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

	err = locker.PartialUpdate(userID, locker.ID)
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
	lockerIDs := []int{}
	for _, locker := range lockers {
		if locker.ID == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "ID is empty"})
			return
		}
		lockerIDs = append(lockerIDs, locker.ID)
	}

	err = lockers.DeleteLockers(userID, lockerIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// RetreiveLocker : Locker 조회
func RetreiveLocker(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"data": locker})
}

// RetreiveLockers : Lockers 조회
func RetreiveLockers(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"data": lockers})
}

// RetrieveAllLockers : 모든 Lockers 조회
func RetreiveAllLocker(c *gin.Context) {
	lockers := models.Lockers{}
	err := lockers.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": lockers})
}
