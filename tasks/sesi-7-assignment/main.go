package main

import (
	"sesi-7-assignment/controllers"
	"sesi-7-assignment/database"

	"github.com/gin-gonic/gin"
)

// assignment ini menggunakan GIN framework dan GORM MySQL

func main() {
	r := gin.Default()
	database.ConnectDatabase()

	// routes
	r.POST("/orders", controllers.CreateOrder)
	r.GET("/orders", controllers.GetAllOrders)
	r.GET("/order/:id", controllers.GetOrderById)
	r.PUT("/orders/:id", controllers.UpdateOrder)
	r.DELETE("/orders/:id", controllers.DeleteOrder)

	// server
	r.Run(":9090")
}
