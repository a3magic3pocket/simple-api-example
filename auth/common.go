package auth

import (
	"errors"
	"os"
	"simple-api-example/models"
	"simple-api-example/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserInfo struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

// GetUserInfoFromBody : 요청 body에서 이용하여 유저정보 획득
func GetUserInfoFromBody(c *gin.Context) (user models.User, err error) {
	userInfo := UserInfo{}
	if err = c.ShouldBindJSON(&userInfo); err != nil {
		return user, errors.New("auth info is invalid")
	}

	err = user.Get(userInfo.UserName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New("your account infomation not exists. you need to sign up")
	}
	if err = utils.ComparePasswords(user.SecretKey, userInfo.Password); err != nil {
		return user, err
	}

	return user, err
}

// GetUserID : 유저 ID 리턴
func GetUserID(c *gin.Context) (userID int, err error) {
	identityKey := os.Getenv("IDENTITY_KEY")
	err = errors.New("auth info is not valid")
	rawUser, exists := c.Get(identityKey)
	if !exists {
		return -1, err
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		return -1, err
	}

	return user.ID, nil
}
