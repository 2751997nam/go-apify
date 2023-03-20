package viewproduct

import (
	"product-service/internal/config"
	"product-service/internal/helpers"
	"product-service/internal/models"
	"product-service/internal/types"

	"golang.org/x/exp/maps"
)

type OptionInterface interface {
	GetVariantId() uint64
}

func GetSkues(product models.Product) GetSkuReponse {
	db := models.GetDB()
	skues := GetProductSkues(product.ID)
	variantIds := map[uint64]uint64{}
	optionIds := map[uint64]uint64{}
	optionByVariant := map[uint64][]uint64{}
	for _, sku := range skues {
		if len(sku.SkuValues) > 0 {
			sku.SkuValues = sortVariants(sku.SkuValues)

			for _, value := range sku.SkuValues {
				if _, ok := variantIds[value.VariantId]; !ok {
					variantIds[value.VariantId] = value.VariantId
				}
				if _, ok := optionIds[value.VariantOptionId]; !ok {
					optionIds[value.VariantOptionId] = value.VariantOptionId
				}
				if _, ok := optionByVariant[value.VariantId]; !ok {
					optionByVariant[value.VariantId] = []uint64{}
				}
				if _, ok := optionByVariant[value.VariantId]; ok {
					optionByVariant[value.VariantId] = append(optionByVariant[value.VariantId], value.VariantOptionId)
				}
			}
		}
	}
	for VariantId := range optionByVariant {
		optionByVariant[VariantId] = helpers.ArrayUnique(optionByVariant[VariantId])
	}

	variantById := map[uint64]models.Variant{}
	if len(variantIds) > 0 {
		items := []models.Variant{}
		db.Model(models.Variant{}).Where("id In ?", maps.Values(variantIds)).Find(&items)
		for _, item := range items {
			variantById[item.ID] = item
		}
	}
	optionById := map[uint64]models.VariantOption{}
	if len(optionIds) > 0 {
		items := []models.VariantOption{}
		db.Model(models.VariantOption{}).Where("id In ?", maps.Values(optionIds)).Find(&items)
		for _, item := range items {
			optionById[item.ID] = item
		}
	}

	variants := buildVariants(variantById, optionById, optionByVariant)
	// variantOptions := map[uint64]types.OptionView{}
	skuViews := []types.ProductSkuView{}

	for _, sku := range skues {
		productSku := buildSku(sku, product, variantById, optionById)
		skuViews = append(skuViews, productSku)
	}

	variantSorted := sortVariants(maps.Values(variants))

	return GetSkuReponse{
		Variants:        variantSorted,
		ProductVariants: skuViews,
	}
}

func buildSku(sku models.ProductSku, product models.Product, variantById map[uint64]models.Variant, optionById map[uint64]models.VariantOption) types.ProductSkuView {
	skuImg := sku.ImageUrl
	skuStatus := sku.Status
	if len(skuImg) == 0 {
		skuImg = product.ImageUrl
	}
	if product.Status != "ACTIVE" {
		skuStatus = product.Status
	}
	productSku := types.ProductSkuView{
		ID:        sku.ID,
		Price:     sku.Price,
		HighPrice: sku.HighPrice,
		IsDefault: sku.IsDefault,
		ProductId: sku.ProductId,
		Status:    skuStatus,
	}

	gallery := []string{skuImg}
	for _, img := range sku.Gallery {
		gallery = append(gallery, img.ImageUrl)
	}
	productSku.Gallery = gallery
	if len(sku.SkuValues) > 0 {
		for _, item := range sku.SkuValues {
			if option, ok := optionById[item.VariantOptionId]; ok {
				productSku.Variants = append(productSku.Variants, option.ID)
			}
		}
	}

	return productSku
}

func buildVariants(variantById map[uint64]models.Variant, optionById map[uint64]models.VariantOption, optionByVariant map[uint64][]uint64) map[uint64]types.VariantView {
	variants := map[uint64]types.VariantView{}
	for key, optionIds := range optionByVariant {
		if _, ok := variantById[key]; ok {
			variant := types.VariantView{
				ID:     key,
				Slug:   variantById[key].Slug,
				Name:   variantById[key].Name,
				Type:   variantById[key].Type,
				Values: []types.OptionView{},
			}
			for _, optionId := range optionIds {
				if option, ok := optionById[optionId]; ok {
					variant.Values = append(variant.Values, types.OptionView{
						ID:        optionId,
						Name:      option.Name,
						ImageUrl:  option.ImageUrl,
						VariantId: option.VariantId,
					})
				}
			}
			variants[key] = variant
		}
	}

	return variants
}

func GetProductSkues(productId uint64) []models.ProductSku {
	retVal := []models.ProductSku{}
	db := models.GetDB()
	db.Model(models.ProductSku{}).Where("product_id = ?", productId).Preload("Gallery", "type", "VARIANT").Preload("SkuValues").Find(&retVal)

	return retVal
}

func sortVariants[T any](options []T) []T {
	sorderConfig := config.GetSetting()
	configTop := sorderConfig["top"].(map[uint64]string)
	configBottom := sorderConfig["bottom"].(map[uint64]string)
	top := []T{}
	bottom := []T{}
	for id := range configTop {
		for _, op := range options {
			switch any(op).(type) {
			case models.ProductSkuValue:
				option := any(op).(models.ProductSkuValue)
				if option.VariantId == id {
					top = append(top, op)
				}
			case types.VariantView:
				option := any(op).(types.VariantView)
				if option.ID == id {
					top = append(top, op)
				}
			}

		}
	}

	for id := range configBottom {
		for _, op := range options {
			switch any(op).(type) {
			case models.ProductSkuValue:
				option := any(op).(models.ProductSkuValue)
				if option.VariantId == id {
					bottom = append(bottom, op)
				}
			case types.VariantView:
				option := any(op).(types.VariantView)
				if option.ID == id {
					bottom = append(bottom, op)
				}
			}
		}
	}

	return append(top, bottom...)
}
