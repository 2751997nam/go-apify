package storeproducts

import (
	"fmt"
	"log"
	"product-service/internal/helpers"
	"product-service/internal/models"

	"github.com/gosimple/slug"
)

func StoreVariants(data map[string]any) map[string]map[string]string {
	retVal := map[string]map[string]string{}
	db := models.GetDB()
	dataVariants := data["variants"].([]any)

	for _, dataVariant := range dataVariants {
		variant := dataVariant.(map[string]any)
		variantValues := variant["values"].([]any)
		if helpers.AnyToString(variant["name"]) == "" || len(variantValues) == 0 {
			continue
		}
		variantSlug := slug.Make(helpers.AnyToString(variant["name"]))
		saveVariant := models.Variant{
			Name: helpers.AnyToString(variant["name"]),
			Slug: variantSlug,
			Type: helpers.AnyToString(variant["type"]),
		}
		if helpers.AnyFloat64ToUint64(variant["id"]) > 0 {
			saveVariant.ID = helpers.AnyFloat64ToUint64(variant["id"])
		} else {
			err := db.Create(&saveVariant).Error
			if err != nil {
				helpers.LogPanic(err)
			}
		}

		for _, val := range variantValues {
			target := map[string]string{}
			value := val.(map[string]any)
			slug := slug.Make(helpers.AnyToString(value["name"]))
			option := models.VariantOption{
				VariantId: saveVariant.ID,
				Name:      helpers.AnyToString(value["name"]),
				Code:      GenerateOptionCode(helpers.AnyToString(value["name"])),
				Slug:      slug,
				ImageUrl:  helpers.AnyToString(value["image_url"]),
			}
			if helpers.AnyFloat64ToUint64(value["id"]) > 0 {
				option.ID = helpers.AnyFloat64ToUint64(value["id"])
				if len(option.ImageUrl) > 0 {
					err := db.Model(&option).Select("ImageUrl").Updates(&option).Error
					log.Println("id", option.ID)
					if err != nil {
						helpers.LogPanic(err)
					}
				}
			} else {
				err := db.Create(&option).Error
				if err != nil {
					helpers.LogPanic(err)
				}
			}
			target["variantOptionId"] = fmt.Sprint(option.ID)
			target["variantId"] = fmt.Sprint(saveVariant.ID)
			retVal[fmt.Sprintf("%s+++%s", variantSlug, slug)] = target
		}
	}

	return retVal
}
