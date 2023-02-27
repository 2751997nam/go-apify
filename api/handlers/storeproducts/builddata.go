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
		ModelId: models.ModelId{
			ID: helpers.AnyToUint(data["id"]),
		},
		Name:               fmt.Sprint(data["name"]),
		Sku:                fmt.Sprint(data["sku"]),
		Slug:               slug.Make(fmt.Sprint(data["name"])),
		ImageUrl:           fmt.Sprint(data["image_url"]),
		Price:              helpers.AnyToFloat(data["price"]),
		HighPrice:          helpers.AnyToFloat(data["high_price"]),
		AddShippingFee:     helpers.AnyToFloat(data["add_shipping_fee"]),
		Status:             fmt.Sprint(data["status"]),
		BrandId:            helpers.AnyToUint(data["brand_id"]),
		ApproveAdvertising: helpers.AnyToInt(data["approve_advertising"]),
		IsTrademark:        helpers.AnyToInt(data["is_trademark"]),
		Content:            fmt.Sprint(data["content"]),
		Description:        fmt.Sprint(data["description"]),
		Note:               fmt.Sprint(data["note"]),
	}

	return retVal
}
