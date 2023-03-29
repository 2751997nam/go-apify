package storeproducts

import (
	"product-service/internal/models"

	"github.com/2751997nam/go-helpers/utils"
	"github.com/gosimple/slug"
)

func buildProductData(data map[string]any) models.Product {
	retVal := models.Product{
		Name:               utils.AnyToString(data["name"]),
		Sku:                utils.AnyToString(data["sku"]),
		Slug:               slug.Make(utils.AnyToString(data["name"])),
		ImageUrl:           utils.AnyToString(data["image_url"]),
		Price:              utils.AnyToFloat(data["price"]),
		HighPrice:          utils.AnyToFloat(data["high_price"]),
		AddShippingFee:     utils.AnyToFloat(data["add_shipping_fee"]),
		Status:             utils.AnyToString(data["status"]),
		BrandId:            utils.AnyFloat64ToUint64(data["brand_id"]),
		ApproveAdvertising: utils.AnyToInt(data["approve_advertising"]),
		IsTrademark:        utils.AnyToInt(data["is_trademark"]),
		Content:            utils.AnyToString(data["content"]),
		Description:        utils.AnyToString(data["description"]),
		Note:               utils.AnyToString(data["note"]),
	}

	if dataId, ok := data["id"]; ok {
		retVal.ID = utils.AnyFloat64ToUint64(dataId)
	}

	return retVal
}
