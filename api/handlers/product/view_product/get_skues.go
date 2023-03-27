package viewproduct

import (
	templateSku "product-service/api/handlers/product/template_sku"
	"product-service/internal/config"
	"product-service/internal/models"
	"product-service/internal/types"
	"sort"

	goHelpers "github.com/2751997nam/go-helpers/pkg/helpers"

	"golang.org/x/exp/maps"
)

type OptionInterface interface {
	GetVariantId() uint64
}

func GetSkues(product models.Product) GetSkuReponse {
	db := models.GetDB()
	pnt := models.ProductNTemplate{}
	db.Where("product_id = ?", product.ID).First(&pnt)
	skues := []models.ProductSku{}
	if pnt.ID > 0 {
		template := models.Template{
			BaseModel: models.BaseModel{
				ID: pnt.TemplateId,
			},
		}
		db.First(&template)
		skues = templateSku.Parse(product.ID, template)
		product.ID = template.ProductIdFake
	} else {
		skues = GetProductSkues(product.ID)
	}
	variantIds := map[uint64]uint64{}
	optionIds := map[uint64]uint64{}
	optionByVariant := map[uint64][]uint64{}
	for key, sku := range skues {
		if len(sku.SkuValues) > 0 {
			skues[key].SkuValues = sortVariants(sku.SkuValues)
			sku.SkuValues = skues[key].SkuValues
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
		optionByVariant[VariantId] = goHelpers.ArrayUnique(optionByVariant[VariantId])
	}

	variantById := map[uint64]models.Variant{}
	if len(variantIds) > 0 {
		items := []models.Variant{}
		db.Model(models.Variant{}).Where("id In ?", maps.Values(variantIds)).Order("id asc").Find(&items)
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
	skuViews := []types.ProductSkuView{}

	for _, sku := range skues {
		productSku := buildSku(sku, product, variantById, optionById)
		skuViews = append(skuViews, productSku)
	}

	variantSorted := sortVariants(maps.Values(variants))
	for key, item := range variantSorted {
		if item.Slug == "size" {
			variantSorted[key].Values = sortSize(variantSorted[key].Values)
			break
		} else {
			variantSorted[key].Values = sortValues(variantSorted[key].Values)
		}
	}

	return GetSkuReponse{
		Variants:        variantSorted,
		ProductVariants: skuViews,
	}
}

func buildSku(sku models.ProductSku, product models.Product, variantById map[uint64]models.Variant, optionById map[uint64]models.VariantOption) types.ProductSkuView {
	skuStatus := sku.Status
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

	gallery := []string{}
	if len(sku.ImageUrl) > 0 {
		gallery = append(gallery, sku.ImageUrl)
	}
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
					optionView := types.OptionView{
						ID:   optionId,
						Name: option.Name,
						Slug: option.Slug,
					}
					if variant.Type == "IMAGE" {
						optionView.ImageUrl = option.ImageUrl
					}
					variant.Values = append(variant.Values, optionView)
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

func sortSize(options []types.OptionView) []types.OptionView {
	sizeSorder := config.SizeSorder
	retVal := []types.OptionView{}
	others := []types.OptionView{}
	for _, slug := range sizeSorder {
		for _, option := range options {
			if option.Slug == slug {
				retVal = append(retVal, option)
				break
			}
		}
	}

	for _, option := range options {
		check := false
		for _, slug := range sizeSorder {
			if option.Slug == slug {
				check = true
				break
			}
		}
		if !check {
			others = append(others, option)
		}
	}

	return append(retVal, others...)
}
func sortValues(options []types.OptionView) []types.OptionView {
	sort.Slice(options, func(i, j int) bool {
		return options[i].ID < options[j].ID
	})

	return options
}

func sortVariants[T any](options []T) []T {
	sorderConfig := config.VariantSorder
	configTop := sorderConfig["top"].([]uint64)
	configBottom := sorderConfig["bottom"].([]uint64)
	check := map[uint64]bool{}
	top := []T{}
	middle := []T{}
	bottom := []T{}
	for _, id := range configTop {
		for _, op := range options {
			switch any(op).(type) {
			case models.ProductSkuValue:
				option := any(op).(models.ProductSkuValue)
				if option.VariantId == id {
					top = append(top, op)
					check[option.VariantId] = true
				}
			case types.VariantView:
				option := any(op).(types.VariantView)
				if option.ID == id {
					top = append(top, op)
					check[option.ID] = true
				}
			}

		}
	}

	for _, id := range configBottom {
		for _, op := range options {
			switch any(op).(type) {
			case models.ProductSkuValue:
				option := any(op).(models.ProductSkuValue)
				if option.VariantId == id {
					bottom = append(bottom, op)
					check[option.VariantId] = true
				}
			case types.VariantView:
				option := any(op).(types.VariantView)
				if option.ID == id {
					bottom = append(bottom, op)
					check[option.ID] = true
				}
			}
		}
	}

	for _, op := range options {
		switch any(op).(type) {
		case models.ProductSkuValue:
			option := any(op).(models.ProductSkuValue)
			if _, ok := check[option.VariantId]; !ok {
				middle = append(middle, op)
			}
		case types.VariantView:
			option := any(op).(types.VariantView)
			if _, ok := check[option.ID]; !ok {
				middle = append(middle, op)
			}
		}
	}

	return append(append(top, middle...), bottom...)
}
