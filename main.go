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
	initializers.ConnectElasticsearch()
	initializers.CheckElasticIndex()
}

func main() {
	r := gin.Default()
	r.Use(middlewares.DatabaseMiddleware(initializers.DB))

	// Auth routes
	a := r.Group("/auth")
	{
		a.POST("/signup/", controllers.Signup)
		a.POST("/login/", controllers.Login)
		a.POST("/validate/", middlewares.Authorization, controllers.Validate)
	}

	// Item routes
	i := r.Group("/item", middlewares.Authorization)
	i.Use(middlewares.ElasticsearchMiddleware(initializers.ES))
	{
		i.POST("/", controllers.ItemCreate)
		i.GET("/", controllers.ItemFindAll)
		i.GET("/:id/", controllers.ItemFindOne)
		i.PUT("/:id/", controllers.ItemUpdate)
		i.DELETE("/:id/", controllers.ItemDelete)
	}
	r.Run()
}
