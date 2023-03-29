package storeproducts

import (
	"fmt"
	"product-service/internal/models"

	"github.com/2751997nam/go-helpers/utils"
	"github.com/gosimple/slug"
)

func StoreVariants(data map[string]any) map[string]map[string]string {
	retVal := map[string]map[string]string{}
	db := models.GetDB()
	dataVariants := data["variants"].([]any)

	for _, dataVariant := range dataVariants {
		variant := dataVariant.(map[string]any)
		variantValues := variant["values"].([]any)
		if utils.AnyToString(variant["name"]) == "" || len(variantValues) == 0 {
			continue
		}
		variantSlug := slug.Make(utils.AnyToString(variant["name"]))
		saveVariant := models.Variant{
			Name: utils.AnyToString(variant["name"]),
			Slug: variantSlug,
			Type: utils.AnyToString(variant["type"]),
		}
		if utils.AnyFloat64ToUint64(variant["id"]) > 0 {
			saveVariant.ID = utils.AnyFloat64ToUint64(variant["id"])
		} else {
			err := db.Create(&saveVariant).Error
			if err != nil {
				utils.LogPanic(err)
			} else {
				utils.QuickLog(saveVariant, saveVariant.ID, "VARIANT", "CREATE")
			}
		}

		for _, val := range variantValues {
			target := map[string]string{}
			value := val.(map[string]any)
			slug := slug.Make(utils.AnyToString(value["name"]))
			option := models.VariantOption{
				VariantId: saveVariant.ID,
				Name:      utils.AnyToString(value["name"]),
				Code:      GenerateOptionCode(utils.AnyToString(value["name"])),
				Slug:      slug,
				ImageUrl:  utils.AnyToString(value["image_url"]),
			}
			if utils.AnyFloat64ToUint64(value["id"]) > 0 {
				option.ID = utils.AnyFloat64ToUint64(value["id"])
				if len(option.ImageUrl) > 0 {
					err := db.Model(&option).Select("ImageUrl").Updates(&option).Error
					if err != nil {
						utils.LogPanic(err)
					} else {
						utils.QuickLog(option, option.ID, "VARIANT_OPTION", "UPDATE")
					}
				}
			} else {
				err := db.Create(&option).Error
				if err != nil {
					utils.LogPanic(err)
				} else {
					utils.QuickLog(option, option.ID, "VARIANT_OPTION", "CREATE")
				}
			}
			target["variantOptionId"] = fmt.Sprint(option.ID)
			target["variantId"] = fmt.Sprint(saveVariant.ID)
			retVal[fmt.Sprintf("%s+++%s", variantSlug, slug)] = target
		}
	}

	return retVal
}
