package storeproducts

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"product-service/internal/helpers"
	"product-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func Store(c *gin.Context) {
	db := models.GetDB()
	var data map[string]any
	bodyAsByteArray, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Panic(err)
		helpers.ResponseFail(c, "something went wrong", http.StatusUnprocessableEntity)
		return
	}
	if err = json.Unmarshal([]byte(bodyAsByteArray), &data); err != nil {
		log.Panic(err)
		helpers.ResponseFail(c, "something went wrong", http.StatusUnprocessableEntity)
		return
	}
	message, ok := validate(data)
	if !ok {
		helpers.ResponseFail(c, message, http.StatusUnprocessableEntity)
		return
	}

	productData := buildProductData(data)

	if productData.ID > 0 {
		err := db.Omit("CreatedAt").Save(&productData).Error
		if err != nil {
			log.Panic(err)
		}
	} else {
		err := db.Create(&productData).Error
		if err != nil {
			log.Panic(err)
		}
	}
	dataCategoryIds, ok := data["categoryIds"].([]any)
	if ok {
		categoryIds := lo.Map(dataCategoryIds, func(x any, index int) uint64 {
			return helpers.AnyFloat64ToUint64(x)
		})
		StoreProductNCategory(productData.ID, categoryIds)
	}

	dataTagIds, ok := data["tagIds"].([]any)
	if ok {
		tagIds := lo.Map(dataTagIds, func(x any, index int) uint64 {
			return helpers.AnyFloat64ToUint64(x)
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

	helpers.ResponseSuccess(c, productData, http.StatusOK)
}
