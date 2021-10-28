package auth

import (
	"encoding/base64"
	"errors"
	"simple-api-example/models"
	"simple-api-example/utils"

	gin "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserInfoUsingBasicAuth : basic auth를 이용하여 유저정보 획득
func GetUserInfoUsingBasicAuth(c *gin.Context) (user models.User, err error) {
	user = models.User{}

	userName, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		return user, errors.New("auth info is invalid")
	}

	err = user.Get(userName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New("your account infomation not exists. you need to sign up")
	}
	if err = utils.ComparePasswords(user.SecretKey, password); err != nil {
		return user, err
	}

	return user, err
}

// GetBasicAuth : username과 password를 받아 base64로 encoding
func GetBasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
