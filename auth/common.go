package auth

import (
	"errors"
	"os"
	"simple-api-example/models"
	"simple-api-example/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserInfo struct {
	UserName string `json:"UserName" form:"UserName"`
	Password string `json:"Password" form:"Password"`
}

// @Summary 요청 body에서 유저정보 획득
// @Description 요청 body에서 유저정보 획득
// @Tags auth
// @Param UserInfo body UserInfo true "UserInfo"
// @Accept  json
// @Produce  json
// @Router /login [post]
// @Success 200 {object} controllers.AuthTokenResponse
// @Failure 400 {object} controllers.SwagFail
// GetUserInfoFromBody : 요청 body에서 유저정보 획득
func GetUserInfoFromBody(c *gin.Context) (user models.User, err error) {
	user = models.User{}

	userInfo := UserInfo{}

	contentType := c.Request.Header["Content-Type"]
	if len(contentType) > 0 && strings.ToLower(contentType[0]) == "application/json" {
		err = c.ShouldBindJSON(&userInfo)
	} else {
		err = c.ShouldBind(&userInfo)
	}
	if err != nil {
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

// @Summary 로그아웃
// @Description 로그아웃
// @Tags auth
// @Accept  json
// @Produce  json
// @Router /logout [post]
// @Success 200 {object} controllers.SwagSucc
// @Failure 400 {object} controllers.SwagFail
// Logout : 로그아웃
func Logout(c *gin.Context) {

}
