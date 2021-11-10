package controllers

import (
	"errors"
	"net/http"
	"simple-api-example/auth"
	"simple-api-example/models"
	"simple-api-example/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserInput : 유저 인풋 구조체
type UserInput struct {
	Name     string `json:"UserName"`
	Group    string `json:"Group"`
	Password string `json:"Password"`
}

// UserOutput : 유저 아웃풋 구조체
type UserOutput struct {
	UserName string `json:"UserName"`
}

// @Summary 유저 생성
// @Description 유저 생성
// @Tags user
// @Param brand body UserInput true "UserInput"
// @Accept  json
// @Produce  json
// @Router /user [post]
// @Success 200 {object} SwagSucc
// @Failure 400 {object} SwagFail
// CreateUser : 유저 생성
func CreateUser(c *gin.Context) {
	userInput := UserInput{}
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "data is invalid"})
		return
	}

	// UserName 중복 여부 체크
	user := models.User{}
	user.Name = userInput.Name
	user.Group = userInput.Group
	err := user.Get(userInput.Name)
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
	user.SecretKey = hash

	err = user.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// @Summary 유저 조회
// @Description 유저 조회
// @Tags user
// @Accept  json
// @Produce  json
// @Router /user [get]
// @Success 200 {object} SwagSuccRetrieveUser
// @Failure 400 {object} SwagFail
// CreateUser : 유저 생성
//  RetrieveUser : 유저 조회
func RetrieveUser(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}
	err = user.GetOwned(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error occured"})
		return
	}

	userOutput := UserOutput{}
	userOutput.UserName = user.Name

	c.JSON(http.StatusOK, gin.H{"data": userOutput})
}
