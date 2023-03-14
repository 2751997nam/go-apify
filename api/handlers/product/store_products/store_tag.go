package storeproducts

import (
	"product-service/internal/helpers"
	"product-service/internal/models"

	"github.com/samber/lo"
)

func StoreTag(productId uint64, tagIds []uint64) {
	db := models.GetDB()
	existedItems := []models.TagRefer{}
	db.Where("product_id = ?", productId).Where("refer_type = ?", "PRODUCT").Find(&existedItems)
	deleteIds := []uint64{}
	existedIds := []uint64{}
	for _, item := range existedItems {
		if lo.Contains(tagIds, item.TagId) {
			existedIds = append(existedIds, item.TagId)
		} else {
			deleteIds = append(deleteIds, item.TagId)
		}
	}
	if len(deleteIds) > 0 {
		err := db.Unscoped().Where("tag_id IN ?", deleteIds).Where("refer_id = ?", productId).Where("refer_type = ?", "PRODUCT").Delete(&models.TagRefer{}).Error
		if err != nil {
			helpers.LogPanic(err)
		}
	}
	storeIds, _ := lo.Difference(tagIds, existedIds)
	if len(storeIds) > 0 {
		storeData := []models.TagRefer{}
		for _, id := range storeIds {
			storeData = append(storeData, models.TagRefer{
				ReferType: "PRODUCT",
				ReferId:   productId,
				TagId:     id,
			})
		}

		err := db.Create(&storeData).Error
		if err != nil {
			helpers.LogPanic(err)
		}
	}
}
