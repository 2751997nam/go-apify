package storeproducts

import (
	"log"
	"net/http"
	"product-service/internal/models"

	goHelpers "github.com/2751997nam/go-helpers/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func Store(c *gin.Context) {
	db := models.GetDB()
	data, err := goHelpers.GetRequestData(c)
	if err != nil {
		log.Println(err)
		goHelpers.ResponseFail(c, "something went wrong", http.StatusUnprocessableEntity)
		return
	}
	message, ok := validate(data)
	if !ok {
		goHelpers.ResponseFail(c, message, http.StatusUnprocessableEntity)
		return
	}

	productData := buildProductData(data)

	if productData.ID > 0 {
		err := db.Omit("CreatedAt").Save(&productData).Error
		if err != nil {
			goHelpers.LogPanic(err)
		}
	} else {
		err := db.Create(&productData).Error
		if err != nil {
			goHelpers.LogPanic(err)
		}
	}
	dataCategoryIds, ok := data["categoryIds"].([]any)
	if ok {
		categoryIds := lo.Map(dataCategoryIds, func(x any, index int) uint64 {
			return goHelpers.AnyFloat64ToUint64(x)
		})
		StoreProductNCategory(productData.ID, categoryIds)
	}

	dataTagIds, ok := data["tagIds"].([]any)
	if ok {
		tagIds := lo.Map(dataTagIds, func(x any, index int) uint64 {
			return goHelpers.AnyFloat64ToUint64(x)
		})
		StoreTag(productData.ID, tagIds)
	}

	metaData := buildProductMeta(data, productData.ID)
	StoreProductMeta(metaData)
	if _, ok := data["productVariants"]; ok {
		if len(data["productVariants"].([]any)) > 0 {
			variantMapping := StoreVariants(data)
			StoreProductSkues(data, variantMapping, productData.ID)
		}
	}

	goHelpers.ResponseSuccess(c, productData, http.StatusOK)
}
