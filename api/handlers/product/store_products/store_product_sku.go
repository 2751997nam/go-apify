package storeproducts

import (
	"product-service/internal/helpers"
	"product-service/internal/models"
	"strconv"

	"github.com/samber/lo"

	"github.com/gosimple/slug"
)

func StoreProductSkues(data map[string]any, variantMapping map[string]map[string]string, productId uint64) ([]uint64, []uint64, []uint64) {
	db := models.GetDB()
	createIds := []uint64{}
	updateIds := []uint64{}
	deleteIds := []uint64{}

	values := data["productVariants"].([]any)
	if len(values) == 0 {
		return createIds, updateIds, deleteIds
	}
	existedSkues := []models.ProductSku{}
	db.Model([]models.ProductSku{}).Preload("SkuValues").Where("product_id = ?", productId).Find(&existedSkues)
	existedSkuesById := map[uint64]models.ProductSku{}

	existedIds := []uint64{}
	variantExistLen := 0
	valueBySkus := map[uint64][]uint64{}
	if len(existedSkues) > 0 {
		variantExistLen = len(existedSkues[0].SkuValues)
		for _, item := range existedSkues {
			existedSkuesById[item.ID] = item
			existedIds = append(existedIds, item.ID)
			for _, skuValue := range item.SkuValues {
				valueBySkus[item.ID] = append(valueBySkus[item.ID], skuValue.ID)
			}
		}
	}
	skuIds := []uint64{}
	firstProductSku := values[0].(map[string]any)
	firstProductSkuOptions := firstProductSku["variants"].([]any)
	optionLen := len(firstProductSkuOptions)
	if variantExistLen != optionLen {
		removeAllSku(productId)
	}

	storeGalleryData := []models.ProductGalleryData{}
	deleteGallerySkuIds := []uint64{}
	storeProductSkuValueData := []models.ProductSkuValue{}
	isUpdate := false

	for _, v := range values {
		value := v.(map[string]any)
		if variantExistLen != optionLen && helpers.AnyFloat64ToUint64(value["id"]) > 0 {
			value["id"] = 0
		}
		productSku := buildProductSku(value, productId)
		skuOptions := value["variants"].([]any)
		if productSku.ID == 0 {
			isUpdate = true
			if len(productSku.Sku) == 0 {
				productSku.Sku = GenerateSkuCode(productId, skuOptions, 30)
			}
			err := db.Create(&productSku).Error
			if err != nil {
				helpers.LogPanic(err)
			}
			createIds = append(createIds, productSku.ID)
			for _, so := range skuOptions {
				skuOption := so.(map[string]any)
				variantSlug := ""
				if len(helpers.AnyToString(skuOption["variant_slug"])) == 0 {
					variantSlug = helpers.AnyToString(skuOption["variant_slug"])
				} else {
					variantSlug = slug.Make(helpers.AnyToString(skuOption["variant"]))
				}
				optionSlug := ""
				if len(helpers.AnyToString(skuOption["slug"])) > 0 {
					optionSlug = helpers.AnyToString(skuOption["slug"])
				} else {
					optionSlug = slug.Make(helpers.AnyToString(skuOption["name"]))
				}
				key := variantSlug + "+++" + optionSlug
				mappingValue, ok := variantMapping[key]
				if ok {
					variantId, _ := strconv.ParseUint(mappingValue["variantId"], 10, 32)
					variantOptionId, _ := strconv.ParseUint(mappingValue["variantOptionId"], 10, 32)
					if variantId > 0 && variantOptionId > 0 {
						storeProductSkuValueData = append(storeProductSkuValueData, models.ProductSkuValue{
							SkuId:           productSku.ID,
							ProductId:       productId,
							VariantId:       variantId,
							VariantOptionId: variantOptionId,
						})
					}
				}
			}
		} else if _, ok := existedSkuesById[productSku.ID]; ok {
			if len(productSku.Sku) == 0 {
				productSku.Sku = GenerateSkuCode(productId, skuOptions, 30)
			}

			if hasChange(existedSkuesById[productSku.ID], productSku) {
				err := db.Where("id = ?", productSku.ID).Omit("CreatedAt").Updates(&productSku).Error
				if err != nil {
					helpers.LogPanic(err)
				}
			}
			skuIds = append(skuIds, productSku.ID)
		}
		if gallery, ok := value["gallery"].([]string); ok {
			if len(gallery) > 0 {
				storeGalleryData = append(storeGalleryData, models.ProductGalleryData{
					ProductId: productId,
					Gallery:   gallery,
					Type:      "VARIANT",
				})
			} else {
				deleteGallerySkuIds = append(deleteGallerySkuIds, productSku.ID)
			}
		} else if isUpdate {
			deleteGallerySkuIds = append(deleteGallerySkuIds, productSku.ID)
		}
	}
	if len(storeProductSkuValueData) > 0 {
		err := db.Model(&models.ProductSkuValue{}).Omit("ID").Create(storeProductSkuValueData).Error
		if err != nil {
			helpers.LogPanic(err)
		}
	}

	if len(storeGalleryData) > 0 {
		bulkStoreGallery(storeGalleryData)
	}
	if len(deleteGallerySkuIds) > 0 {
		for _, chunk := range lo.Chunk(deleteGallerySkuIds, 100) {
			err := db.Unscoped().Where("product_id IN ?", chunk).Where("type = ?", "VARIANT").Delete(&models.ProductGallery{}).Error
			if err != nil {
				helpers.LogPanic(err)
			}
		}
	}
	if isUpdate {
		if variantExistLen == optionLen && len(existedIds) > 0 && len(skuIds) > 0 {
			for _, id := range existedIds {
				if _, ok := valueBySkus[id]; ok && !lo.Contains(skuIds, id) {
					deleteIds = append(deleteIds, id)
				}
			}
		}

		if len(deleteIds) > 0 {
			for _, chunk := range lo.Chunk(deleteIds, 100) {
				err := db.Where("sku_id IN (?)", chunk).Delete(&models.ProductSkuValue{}).Error
				if err != nil {
					helpers.LogPanic(err)
				}
				err = db.Where("id IN (?)", chunk).Delete(&models.ProductSku{}).Error
				if err != nil {
					helpers.LogPanic(err)
				}
			}
		}
	}

	return createIds, updateIds, deleteIds
}

func hasChange(old models.ProductSku, new models.ProductSku) bool {
	if old.Sku != new.Sku {
		return true
	}
	if old.Price != new.Price {
		return true
	}
	if old.HighPrice != new.HighPrice {
		return true
	}
	if old.ImageUrl != new.ImageUrl {
		return true
	}
	if old.Status != new.Status {
		return true
	}

	if old.IsDefault != new.IsDefault {
		return true
	}

	return false
}

func buildProductSku(input map[string]any, productId uint64) models.ProductSku {
	retVal := models.ProductSku{
		BaseModel: models.BaseModel{
			ID: helpers.AnyFloat64ToUint64(input["id"]),
		},
		ProductId: productId,
		ImageUrl:  helpers.AnyToString(input["image_url"]),
		Price:     0,
		HighPrice: 0,
		Sku:       helpers.AnyToString(input["sku"]),
		IsDefault: helpers.AnyToInt(input["is_default"]),
	}

	if helpers.AnyToFloat(input["price"]) > 0 {
		retVal.Price = helpers.AnyToFloat(input["price"])
	}

	if helpers.AnyToFloat(input["high_price"]) > 0 {
		retVal.HighPrice = helpers.AnyToFloat(input["high_price"])
	}

	if len(helpers.AnyToString(input["status"])) > 0 {
		retVal.Status = helpers.AnyToString(input["status"])
	} else {
		retVal.Status = "ACTIVE"
	}

	return retVal
}

func removeAllSku(productId uint64) {
	db := models.GetDB()
	err := db.Model(&models.ProductSkuValue{}).Where("product_id = ?", productId).Delete(&models.ProductSkuValue{}).Error
	if err != nil {
		helpers.LogPanic(err)
	}
	err = db.Model(&models.ProductSku{}).Where("product_id = ?", productId).Delete(&models.ProductSku{}).Error
	if err != nil {
		helpers.LogPanic(err)
	}
}
