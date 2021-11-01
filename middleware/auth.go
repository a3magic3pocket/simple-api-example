package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"simple-api-example/auth"
	"simple-api-example/models"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	gin "github.com/gin-gonic/gin"
)

// GetAuthMiddleware : authMiddleware 획득
func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	var identityKey = os.Getenv("IDENTITY_KEY")

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "test zone",
		Key:              []byte(os.Getenv("AUTH_SECRET_KEY")),
		Timeout:          time.Hour * 24,
		MaxRefresh:       time.Hour * 24,
		SendCookie:       true,
		CookieHTTPOnly:   true,
		SecureCookie:     true,
		CookieMaxAge:     time.Hour * 24,
		CookieName:       "simple-locker-auth",
		IdentityKey:      identityKey,
		SigningAlgorithm: "HS256",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			fmt.Println("In PayloadFunc")
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			fmt.Println("In IdentityHandler")
			claims := jwt.ExtractClaims(c)
			return &models.User{
				ID: int(claims[identityKey].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			fmt.Println("In Authenticator")
			user := models.User{}
			var err error

			switch c.FullPath() {
			case "/login/basic":
				user, err = auth.GetUserInfoUsingBasicAuth(c)
			case "/login":
				user, err = auth.GetUserInfoFromBody(c)
			default:
				log.Printf("auth c.FullPath is not allowed, c.FullPath: %s", c.FullPath())
			}
			if err != nil {
				log.Printf("authentication failed. error: %s\n", err.Error())
				return nil, jwt.ErrFailedAuthentication
			}

			return &user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			fmt.Println("In Authorizator")
			result := true

			return result
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Println("In Unauthorized")
			c.JSON(code, gin.H{
				"error": gin.H{
					"code":    code,
					"message": message,
				},
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			statusCode := http.StatusOK
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"code":    statusCode,
					"token":   token,
					"expire":  expire.Format(time.RFC3339),
					"message": "success",
				},
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			statusCode := http.StatusOK
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"code":    statusCode,
					"message": "success",
				},
			})
		},

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		// TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenLookup: "header: Authorization, cookie:simple-locker-auth",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	return authMiddleware, err
}
