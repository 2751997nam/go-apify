package main

import (
	"product-service/api/handlers"
	listcategory "product-service/api/handlers/category/list_category"
	storeCategory "product-service/api/handlers/category/store_category"
	listProducts "product-service/api/handlers/product/list_products"
	showProduct "product-service/api/handlers/product/show_product"
	storeProducts "product-service/api/handlers/product/store_products"
	viewProduct "product-service/api/handlers/product/view_product"
	"product-service/api/middlewares"

	"github.com/gin-gonic/gin"
)

func routes() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	productRouter := router.Group("/")
	{
		productRouter.GET("", handlers.Home)
		productRouter.GET("find", listProducts.Find)
		productRouter.GET("view/:id", viewProduct.View)
		productRouter.GET("variant/:id", viewProduct.ViewVariant)
		productRouter.GET("show/:id", showProduct.Show)
		productRouter.POST("", storeProducts.Store)
		productRouter.PATCH(":id", storeProducts.Store)
		productRouter.PUT(":id", storeProducts.Store)
	}

	categoryRouter := router.Group("/category")
	{
		// productRouter.GET("", listProducts.Find)
		categoryRouter.GET("/tree", listcategory.GetTree)
		categoryRouter.POST("", storeCategory.Store)
		categoryRouter.PATCH(":id", storeCategory.Store)
		categoryRouter.PUT(":id", storeCategory.Store)
	}
	router.GET("/favicon.ico", handlers.DoNothing)

	return router
}
