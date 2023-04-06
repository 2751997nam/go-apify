package main

import (
	"apify-service/api/handlers"
	"apify-service/api/handlers/table"
	"apify-service/api/middlewares"

	"github.com/gin-gonic/gin"
)

func routes() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	tableRouter := router.Group(":table")
	{
		tableRouter.GET("", table.Find)
		tableRouter.GET(":id", table.Show)
		tableRouter.POST("", table.Store)
		tableRouter.PATCH(":id", table.Update)
		tableRouter.PUT(":id", table.Update)
		tableRouter.DELETE(":id", table.Delete)
	}
	router.GET("/favicon.ico", handlers.DoNothing)

	return router
}
