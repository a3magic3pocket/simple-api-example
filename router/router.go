package router

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"simple-api-example/controllers"
	"simple-api-example/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter : 라우터 세팅
func SetupRouter() *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = []string{
		"Authorization",
		"Content-type",
	}
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{
		"GET", "POST", "DELETE", "PATCH",
	}
	router.Use(cors.New(corsConfig))

	authMiddleware, err := middleware.GetAuthMiddleware()
	if err != nil {
		log.Fatal("JWT Error: " + err.Error())
	}

	// 정적 파일 로드
	var (
		_, b, _, _     = runtime.Caller(0)
		workingDirPath = filepath.Dir(filepath.Dir(b))
	)
	router.LoadHTMLGlob(fmt.Sprintf("%s/views/*", workingDirPath))
	router.StaticFS("/public", http.Dir(fmt.Sprintf("%s/public", workingDirPath)))
	router.GET("favicon.ico", func(c *gin.Context) {
		c.File("public/images/favicon.ico")
	})

	// 인증
	router.POST("/login/basic", authMiddleware.LoginHandler)
	router.POST("/login", authMiddleware.LoginHandler)
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	router.POST("/logout", authMiddleware.LogoutHandler)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page Not Found"})
	})

	// 유저
	user := router.Group("/user")
	user.POST("", controllers.CreateUser)
	user.Use(authMiddleware.MiddlewareFunc())
	{
		user.GET("", controllers.RetrieveUser)
	}

	// 로커
	lockers := router.Group("/lockers")
	lockers.Use(authMiddleware.MiddlewareFunc())
	{
		lockers.GET("", controllers.RetrieveLockers)
		lockers.GET("all", controllers.RetrieveAllLocker)
		lockers.POST("", controllers.CreateLockers)
		lockers.PATCH("", controllers.UpdateLockers)
		lockers.DELETE("", controllers.DeleteLockers)
	}

	locker := router.Group("/locker")
	locker.Use(authMiddleware.MiddlewareFunc())
	{
		locker.GET(":id", controllers.RetrieveLocker)
		locker.PATCH("", controllers.UpdateLocker)
	}

	return router
}
