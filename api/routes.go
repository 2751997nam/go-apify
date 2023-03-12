package main

import (
	"product-service/api/handlers"
	"product-service/api/handlers/listproducts"
	"product-service/api/handlers/storeproducts"
	"product-service/api/middlewares"

	"github.com/gin-gonic/gin"
)

func routes() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	router.GET("/", handlers.Home)
	productRouter := router.Group("/product")
	{
		productRouter.GET("", listproducts.Find)
		productRouter.POST("", storeproducts.Store)
		productRouter.PATCH(":id", storeproducts.Store)
		productRouter.PUT(":id", storeproducts.Store)
	}
	router.GET("/favicon.ico", handlers.DoNothing)

	return router
}
