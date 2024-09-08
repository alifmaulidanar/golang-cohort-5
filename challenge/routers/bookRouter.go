package routers

import (
	"golang-cohort-5/tasks/sesi-5/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/book", controllers.CreateBook)

	router.GET("/books", controllers.GetAllBooks)

	router.GET("/book/:bookID", controllers.GetBook)

	router.PUT("/book/:bookID", controllers.UpdateBook)

	router.DELETE("/book/:bookID", controllers.DeleteBook)

	return router
}
