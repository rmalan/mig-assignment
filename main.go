package main

import (
	"rmalan/go/mig-assignment/config"
	"rmalan/go/mig-assignment/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	// config.DbConnect("root:@tcp(localhost:3306)/mig_assignment?parseTime=true")
	config.DbConnect()
	config.DbMigrate()

	// Initialize Router
	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", controllers.Home)
	router.POST("/auth/register", controllers.RegisterUser)
	router.POST("/auth/login", controllers.Login)

	return router
}
