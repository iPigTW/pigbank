package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ipigtw/api/controllers"
	"github.com/ipigtw/api/initialize"
	"github.com/ipigtw/api/middleware"
)
func init() {
	initialize.LoadEnv()
	initialize.ConnectToDb()
	initialize.SyncDB()
}
func main() {
	router := gin.Default()
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/transfer", middleware.RequireAuth, controllers.Transfer)
	router.Run()
}
