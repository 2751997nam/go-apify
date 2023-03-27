package storeproducts

import (
	"fmt"
	"product-service/internal/models"

	goHelpers "github.com/2751997nam/go-helpers/pkg/helpers"

	"github.com/gosimple/slug"
)

func buildProductData(data map[string]any) models.Product {
	fmt.Println("product id : ", data["id"])
	retVal := models.Product{
		BaseModel: models.BaseModel{
			ID: goHelpers.AnyFloat64ToUint64(data["id"]),
		},
		Name:               goHelpers.AnyToString(data["name"]),
		Sku:                goHelpers.AnyToString(data["sku"]),
		Slug:               slug.Make(goHelpers.AnyToString(data["name"])),
		ImageUrl:           goHelpers.AnyToString(data["image_url"]),
		Price:              goHelpers.AnyToFloat(data["price"]),
		HighPrice:          goHelpers.AnyToFloat(data["high_price"]),
		AddShippingFee:     goHelpers.AnyToFloat(data["add_shipping_fee"]),
		Status:             goHelpers.AnyToString(data["status"]),
		BrandId:            goHelpers.AnyFloat64ToUint64(data["brand_id"]),
		ApproveAdvertising: goHelpers.AnyToInt(data["approve_advertising"]),
		IsTrademark:        goHelpers.AnyToInt(data["is_trademark"]),
		Content:            goHelpers.AnyToString(data["content"]),
		Description:        goHelpers.AnyToString(data["description"]),
		Note:               goHelpers.AnyToString(data["note"]),
	}

	return retVal
}
