package storeproducts

import (
	"fmt"
	"product-service/internal/helpers"
	"product-service/internal/models"

	"github.com/gosimple/slug"
)

func buildProductData(data map[string]any) models.Product {
	fmt.Println("product id : ", data["id"])
	retVal := models.Product{
		BaseModel: models.BaseModel{
			ID: helpers.AnyFloat64ToUint64(data["id"]),
		},
		Name:               helpers.AnyToString(data["name"]),
		Sku:                helpers.AnyToString(data["sku"]),
		Slug:               slug.Make(helpers.AnyToString(data["name"])),
		ImageUrl:           helpers.AnyToString(data["image_url"]),
		Price:              helpers.AnyToFloat(data["price"]),
		HighPrice:          helpers.AnyToFloat(data["high_price"]),
		AddShippingFee:     helpers.AnyToFloat(data["add_shipping_fee"]),
		Status:             helpers.AnyToString(data["status"]),
		BrandId:            helpers.AnyFloat64ToUint64(data["brand_id"]),
		ApproveAdvertising: helpers.AnyToInt(data["approve_advertising"]),
		IsTrademark:        helpers.AnyToInt(data["is_trademark"]),
		Content:            helpers.AnyToString(data["content"]),
		Description:        helpers.AnyToString(data["description"]),
		Note:               helpers.AnyToString(data["note"]),
	}

	return retVal
}
