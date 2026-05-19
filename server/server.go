package server

import (
	"CricTail_Backend/handler"
	"CricTail_Backend/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func StartServer(serverPort string) (router *gin.Engine) {
	fmt.Println("Starting Server")
	router = gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"serverStatus": "Running",
		})
	})

	router.POST("/register", handler.RegisterUser)
	router.POST("/login", handler.LoginUser)

	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/logout", handler.LogoutUser) //either POST(mostly) or DELETE

	return

}
