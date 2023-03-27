package storeproducts

import (
	"product-service/internal/models"

	goHelpers "github.com/2751997nam/go-helpers/pkg/helpers"

	"github.com/samber/lo"
)

func StoreProductNCategory(productId uint64, categoryIds []uint64) {
	db := models.GetDB()
	existedItems := []models.ProductNCategory{}
	db.Where("product_id = ?", productId).Find(&existedItems)
	deleteIds := []uint64{}
	existedIds := []uint64{}
	for _, item := range existedItems {
		if lo.Contains(categoryIds, item.CategoryId) {
			existedIds = append(existedIds, item.CategoryId)
		} else {
			deleteIds = append(deleteIds, item.CategoryId)
		}
	}
	if len(deleteIds) > 0 {
		err := db.Unscoped().Where("category_id IN ?", deleteIds).Where("product_id = ?", productId).Delete(&models.ProductNCategory{}).Error
		if err != nil {
			goHelpers.LogPanic(err)
		}
	}
	storeIds, _ := lo.Difference(categoryIds, existedIds)
	if len(storeIds) > 0 {
		storeData := []models.ProductNCategory{}
		for _, id := range storeIds {
			storeData = append(storeData, models.ProductNCategory{
				ProductId:  productId,
				CategoryId: id,
			})
		}

		db.Create(&storeData)
	}
}
