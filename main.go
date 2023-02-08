package main

import (
	"github.com/gin-gonic/gin"
	"github.com/okozzie/gogogo/controllers"
	"github.com/okozzie/gogogo/initializers"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()

	//Api route endpoints
	router.GET("/api/ship", controllers.Index)
	router.GET("/api/ship/:id", controllers.Show)
	router.POST("/api/ship", controllers.Store)
	router.PATCH("/api/ship/:id", controllers.Update)
	router.DELETE("/api/ship/:id", controllers.Delete)

	//Not needed but just kept here anyway
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}
