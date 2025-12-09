package routers

import (
	"app/appMessaggistica/controllers"
	"app/appMessaggistica/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.StaticFile("/", "./index.html")

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		api.POST("/messages", controllers.SendMessage)
		api.GET("/messages", controllers.GetMessages)
	}
	return r
}
