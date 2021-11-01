package router

import (
	"log"
	"net/http"
	"simple-api-example/controllers"
	"simple-api-example/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter : 라우터 세팅
func SetupRouter() *gin.Engine {
	router := gin.Default()

	authMiddleware, err := middleware.GetAuthMiddleware()
	if err != nil {
		log.Fatal("JWT Error: " + err.Error())
	}

	// 인증
	router.POST("/login/basic", authMiddleware.LoginHandler)
	router.POST("/logout", authMiddleware.LogoutHandler)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page Not Found"})
	})

	// 유저
	user := router.Group("/user")
	user.POST("", controllers.CreateUser)

	// 로커
	lockers := router.Group("/lockers")
	lockers.Use(authMiddleware.MiddlewareFunc())
	{
		lockers.GET("", controllers.RetreiveLockers)
		lockers.GET("all", controllers.RetreiveAllLocker)
		lockers.POST("", controllers.CreateLockers)
		lockers.DELETE("", controllers.DeleteLockers)
	}

	locker := router.Group("/locker")
	locker.Use(authMiddleware.MiddlewareFunc())
	{
		locker.GET(":id", controllers.RetreiveLocker)
		locker.PATCH("", controllers.UpdateLocker)
	}

	return router
}
