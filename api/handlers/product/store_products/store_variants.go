package storeproducts

import (
	"fmt"
	"log"
	"product-service/internal/models"

	goHelpers "github.com/2751997nam/go-helpers/pkg/helpers"

	"github.com/gosimple/slug"
)

func StoreVariants(data map[string]any) map[string]map[string]string {
	retVal := map[string]map[string]string{}
	db := models.GetDB()
	dataVariants := data["variants"].([]any)

	for _, dataVariant := range dataVariants {
		variant := dataVariant.(map[string]any)
		variantValues := variant["values"].([]any)
		if goHelpers.AnyToString(variant["name"]) == "" || len(variantValues) == 0 {
			continue
		}
		variantSlug := slug.Make(goHelpers.AnyToString(variant["name"]))
		saveVariant := models.Variant{
			Name: goHelpers.AnyToString(variant["name"]),
			Slug: variantSlug,
			Type: goHelpers.AnyToString(variant["type"]),
		}
		if goHelpers.AnyFloat64ToUint64(variant["id"]) > 0 {
			saveVariant.ID = goHelpers.AnyFloat64ToUint64(variant["id"])
		} else {
			err := db.Create(&saveVariant).Error
			if err != nil {
				goHelpers.LogPanic(err)
			}
		}

		for _, val := range variantValues {
			target := map[string]string{}
			value := val.(map[string]any)
			slug := slug.Make(goHelpers.AnyToString(value["name"]))
			option := models.VariantOption{
				VariantId: saveVariant.ID,
				Name:      goHelpers.AnyToString(value["name"]),
				Code:      GenerateOptionCode(goHelpers.AnyToString(value["name"])),
				Slug:      slug,
				ImageUrl:  goHelpers.AnyToString(value["image_url"]),
			}
			if goHelpers.AnyFloat64ToUint64(value["id"]) > 0 {
				option.ID = goHelpers.AnyFloat64ToUint64(value["id"])
				if len(option.ImageUrl) > 0 {
					err := db.Model(&option).Select("ImageUrl").Updates(&option).Error
					log.Println("id", option.ID)
					if err != nil {
						goHelpers.LogPanic(err)
					}
				}
			} else {
				err := db.Create(&option).Error
				if err != nil {
					goHelpers.LogPanic(err)
				}
			}
			target["variantOptionId"] = fmt.Sprint(option.ID)
			target["variantId"] = fmt.Sprint(saveVariant.ID)
			retVal[fmt.Sprintf("%s+++%s", variantSlug, slug)] = target
		}
	}

	return retVal
}
