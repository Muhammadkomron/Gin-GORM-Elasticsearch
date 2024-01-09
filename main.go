package main

import (
	"gin-gorm-tutorial/controllers"
	"gin-gorm-tutorial/initializers"
	"gin-gorm-tutorial/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
	initializers.SyncDatabase()
	initializers.ConnectElastic()
	initializers.CheckElasticIndex()
}

func main() {
	r := gin.Default()
	u := r.Group("/auth")
	{
		u.POST("/signup/", controllers.Signup)
		u.POST("/login/", controllers.Login)
		u.POST("/validate/", middlewares.RequireAuth, controllers.Validate)
	}
	//i := r.Group("/item", middlewares.RequireAuth)
	i := r.Group("/item")
	{
		i.POST("/", controllers.ItemCreate)
		i.GET("/", controllers.ItemFindAll)
		i.GET("/:id/", controllers.ItemFindOne)
		i.POST("/:id/", controllers.ItemUpdate)
		i.DELETE("/:id/", controllers.ItemDelete)
	}
	r.Run()
}
