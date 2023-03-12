package main

import (
	"product-service/api/handlers"
	storeCategory "product-service/api/handlers/category/store_category"
	listProducts "product-service/api/handlers/product/list_products"
	storeProducts "product-service/api/handlers/product/store_products"
	"product-service/api/middlewares"

	"github.com/gin-gonic/gin"
)

func routes() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	router.GET("/", handlers.Home)
	productRouter := router.Group("/product")
	{
		productRouter.GET("", listProducts.Find)
		productRouter.POST("", storeProducts.Store)
		productRouter.PATCH(":id", storeProducts.Store)
		productRouter.PUT(":id", storeProducts.Store)
	}

	categoryRouter := router.Group("/category")
	{
		// productRouter.GET("", listProducts.Find)
		categoryRouter.POST("", storeCategory.Store)
		categoryRouter.PATCH(":id", storeCategory.Store)
		categoryRouter.PUT(":id", storeCategory.Store)
	}
	router.GET("/favicon.ico", handlers.DoNothing)

	return router
}
