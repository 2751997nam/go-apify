package storeproducts

import (
	"encoding/json"
	"product-service/internal/models"

	goHelpers "github.com/2751997nam/go-helpers/pkg/helpers"
)

func buildProductMeta(data map[string]any, productId uint64) []models.ProductMeta {
	retVal := []models.ProductMeta{}
	seoMeta := map[string]string{}
	if data["meta_title"] != nil {
		seoMeta["meta_title"] = goHelpers.AnyToString(data["meta_title"])
	}
	if data["meta_description"] != nil {
		seoMeta["meta_description"] = goHelpers.AnyToString(data["meta_description"])
	}
	if data["meta_keywords"] != nil {
		seoMeta["meta_keywords"] = goHelpers.AnyToString(data["meta_keywords"])
	}

	if len(seoMeta) > 0 {
		jsonStr, _ := json.Marshal(seoMeta)
		retVal = append(retVal, models.ProductMeta{
			ProductId: productId,
			Key:       "seo",
			Value:     string(jsonStr),
		})
	}

	if goHelpers.AnyToInt(data["is_custom_design"]) == 0 {
		retVal = append(retVal, models.ProductMeta{
			ProductId: productId,
			Key:       "is_custom_design",
			Value:     "1",
		})
	}

	return retVal
}

func StoreProductMeta(data []models.ProductMeta) {
	db := models.GetDB()
	err := db.Model(&models.ProductMeta{}).Create(data).Error
	if err != nil {
		goHelpers.LogPanic(err)
	}
}
