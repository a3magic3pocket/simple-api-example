package controllers

import (
	"errors"
	"net/http"
	"simple-api-example/models"
	"simple-api-example/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserInput : 유저 인풋 구조체
type UserInput struct {
	models.User
	Password string `json:"Password"`
}

// CreateUser : 유저 생성
func CreateUser(c *gin.Context) {
	userInput := UserInput{}
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "data is invalid"})
		return
	}

	// UserName 중복 여부 체크
	err := userInput.User.Get(userInput.User.Name)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "username is duplicated"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	hash, err := utils.HashAndSalt(userInput.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password too short."})
		return
	}
	userInput.User.SecretKey = hash

	err = userInput.User.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
