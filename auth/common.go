package auth

import (
	"errors"
	"os"
	"simple-api-example/models"

	"github.com/gin-gonic/gin"
)

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
