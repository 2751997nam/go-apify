package storeproducts

import (
	"fmt"
	"product-service/internal/helpers"
	"product-service/internal/models"
)

func validateTitle(data map[string]any) (string, bool) {
	var message string
	ok := true

	if name := data["name"]; len(helpers.AnyToString(name)) == 0 {
		message = "Tiêu đề sản phẩm không được bỏ trống"
		ok = false
	}

	return message, ok
}

type SkuExistData struct {
	Sku          string
	NotProductId uint64
	NotId        uint64
}

func CheckSkuExistsInProduct(data SkuExistData) bool {
	db := models.GetDB()
	var item models.Product
	query := db.Model(&models.Product{}).Where("sku = ?", data.Sku)
	if data.NotId != 0 {
		query.Where("id != ?", data.NotId)
	}
	if data.NotProductId != 0 {
		query.Where("product_id != ?", data.NotProductId)
	}
	query.First(&item)
	return item.ID > 0
}

func CheckSkuExistsInProductSku(data SkuExistData) bool {
	db := models.GetDB()
	var item models.ProductSku
	query := db.Model(&models.ProductSku{}).Where("sku = ?", data.Sku)
	if data.NotId != 0 {
		query.Where("id != ?", data.NotId)
	}
	if data.NotProductId != 0 {
		query.Where("product_id != ?", data.NotProductId)
	}
	query.First(&item)
	return item.ID > 0
}

func validateSku(data map[string]any) (string, bool) {
	var message string
	isOk := true

	sku := data["sku"]
	dataId := data["id"].(float64)
	productId := uint64(dataId)
	values, ok := data["productVariants"]
	if ok {
		var productVariants []any = values.([]any)
		for _, value := range productVariants {
			tmp := value.(map[string]any)
			item := models.ProductSku{
				Sku: helpers.AnyToString(tmp["sku"]),
				ModelId: models.ModelId{
					ID: helpers.AnyFloat64ToUint64(tmp["id"]),
				},
			}
			if len(item.Sku) > 0 {
				if CheckSkuExistsInProduct(SkuExistData{Sku: item.Sku, NotId: productId}) {
					message = fmt.Sprintf("Mã %s đã tồn tại trong hệ thống", item.Sku)
					isOk = false
					break
				} else if CheckSkuExistsInProductSku(SkuExistData{Sku: item.Sku, NotId: item.ID}) {
					message = fmt.Sprintf("Mã %s đã tồn tại trong hệ thống", item.Sku)
					isOk = false
					break
				}
			}
		}
	} else if len(helpers.AnyToString(sku)) > 0 {
		if CheckSkuExistsInProduct(SkuExistData{Sku: helpers.AnyToString(sku), NotId: productId}) {
			message = fmt.Sprintf("Mã %s đã tồn tại trong hệ thống", sku)
			isOk = false
		} else if CheckSkuExistsInProductSku(SkuExistData{Sku: helpers.AnyToString(sku), NotProductId: productId}) {
			message = fmt.Sprintf("Mã %s đã tồn tại trong hệ thống", sku)
			isOk = false
		}
	}

	return message, isOk
}

func validate(data map[string]any) (string, bool) {
	message, ok := validateTitle(data)
	if ok {
		message, ok = validateSku(data)
	}

	return message, ok
}
