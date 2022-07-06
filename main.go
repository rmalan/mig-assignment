package main

import (
	"rmalan/go/mig-assignment/config"
	"rmalan/go/mig-assignment/controllers"
	"rmalan/go/mig-assignment/middlewares"

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

	secured := router.Use(middlewares.Auth())
	{
		secured.GET("/attendance", controllers.AttendanceByUser)
		secured.PUT("/check-in", controllers.CheckIn)
		secured.PUT("/check-out", controllers.CheckOut)

		secured.GET("/activity", controllers.ActivityByUser)
		secured.GET("/activity/:date", controllers.ActivityByUserByDate)
		secured.POST("/activity", controllers.Store)
		secured.PATCH("/activity/:id", controllers.Update)
		secured.DELETE("/activity/:id", controllers.Destroy)
	}

	return router
}
