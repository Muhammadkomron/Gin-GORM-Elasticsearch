package main

import (
	"gin-gorm-tutorial/controllers"
	"gin-gorm-tutorial/initializers"
	"gin-gorm-tutorial/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/signup/", controllers.Signup)
	r.POST("/login/", controllers.Login)
	r.POST("/validate/", middlewares.RequireAuth, controllers.Validate)
	r.Run()
}
