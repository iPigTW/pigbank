package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ipigtw/pigbank/controllers"
	"github.com/ipigtw/pigbank/initialize"
	"github.com/ipigtw/pigbank/middleware"
)
func init() {
	initialize.LoadEnv()
	initialize.ConnectToDb()
	initialize.SyncDB()
}
func main() {
	router := gin.Default()
	router.Any("/signup", controllers.Signup)
	router.Any("/login", controllers.Login)
	router.Any("/transfer", middleware.RequireAuth, controllers.Transfer)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)
	

	router.Run()
}
