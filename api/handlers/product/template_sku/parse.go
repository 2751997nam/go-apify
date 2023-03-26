package templatesku

import (
	"fmt"
	"product-service/internal/models"
	"sort"
	"strings"
)

func Parse(productId uint64, template models.Template) []models.ProductSku {
	retVal := []models.ProductSku{}
	db := models.GetDB()

	skues := []models.TemplateSku{}
	db.Model(&models.TemplateSku{}).Where("template_id = ?", template.ID).Preload("Gallery", "type", "VARIANT").Preload("SkuValues").Find(&skues)
	DecoreItem(productId, &skues)
	for _, item := range skues {
		sku := models.ProductSku{
			BaseModel: models.BaseModel{
				ID: item.ID,
			},
			ProductId: template.ProductIdFake,
			Sku:       item.Sku,
			Price:     item.Price,
			HighPrice: item.HighPrice,
			ImageUrl:  item.ImageUrl,
			IsDefault: item.IsDefault,
			Status:    item.Status,
			Gallery:   []models.ProductGallery{},
			SkuValues: []models.ProductSkuValue{},
		}
		for _, img := range item.Gallery {
			gal := models.ProductGallery{
				ProductId: item.ID,
				ImageUrl:  img.ImageUrl,
				Type:      img.Type,
			}
			sku.Gallery = append(sku.Gallery, gal)
		}
		for _, value := range item.SkuValues {
			val := models.ProductSkuValue{
				ProductId:       productId,
				SkuId:           value.SkuId,
				VariantId:       value.VariantId,
				VariantOptionId: value.VariantOptionId,
			}
			sku.SkuValues = append(sku.SkuValues, val)
		}
		retVal = append(retVal, sku)
	}

	return retVal
}

func DecoreItem(productId uint64, items *[]models.TemplateSku) {
	for _, item := range *items {
		DecoreSkuCode(&item, productId)
	}
	overwriteImgs := GetOverwriteImageUrls(productId)
	if len(overwriteImgs) > 0 {
		AppendOverwriteImageUrl(items, overwriteImgs)
	} else {
		AppendDesignImageUrl(productId, items)
	}
}

func GetOverwriteImageUrls(productId uint64) map[uint64][]models.TemplateGalleryOverWrite {
	db := models.GetDB()
	items := []models.TemplateGalleryOverWrite{}
	db.Where("product_id = ?", productId).Find(&items)

	retVal := map[uint64][]models.TemplateGalleryOverWrite{}

	for _, item := range items {
		if _, ok := retVal[item.ProductSkuId]; !ok {
			retVal[item.ProductSkuId] = []models.TemplateGalleryOverWrite{}
		}
		retVal[item.ProductSkuId] = append(retVal[item.ProductSkuId], item)
	}
	for key := range retVal {
		sort.Slice(retVal[key], func(a, b int) bool {
			if retVal[key][a].IsPrimary > retVal[key][b].IsPrimary {
				return true
			}
			return retVal[key][a].ID < retVal[key][b].ID
		})
	}

	return retVal
}

func getTemplateDesign(productId uint64) string {
	code := ""
	designCode := models.TemplateDesignCode{}
	db := models.GetDB()
	db.Where("product_id = ?", productId).First(&designCode)
	if designCode.ID > 0 {
		code = designCode.DesignCode
	}

	return code
}

func AppendOverwriteImageUrl(items *[]models.TemplateSku, overwriteImgs map[uint64][]models.TemplateGalleryOverWrite) {
	for key, item := range *items {
		if values, ok := overwriteImgs[item.ID]; ok && len(values) > 0 {
			(*items)[key].ImageUrl = values[0].ImageUrl
			if len(values) > 1 {
				(*items)[key].Gallery = []models.TemplateGallery{}
				for _, img := range values[1:] {
					(*items)[key].Gallery = append((*items)[key].Gallery, models.TemplateGallery{
						ImageUrl:  img.ImageUrl,
						ProductId: item.ID,
						Type:      "VARIANT",
					})
				}
			}
		} else {
			(*items)[key].ImageUrl = ""
			(*items)[key].Gallery = nil
		}
	}
}

func AppendDesignImageUrl(productId uint64, items *[]models.TemplateSku) {
	code := getTemplateDesign(productId)
	if len(code) > 0 {
		for key := range *items {
			(*items)[key].ImageUrl = strings.Replace((*items)[key].ImageUrl, "[DESIGN_ID]", code, -1)
		}
	}
}

func DecoreSkuCode(item *models.TemplateSku, productId uint64) {
	item.Sku = fmt.Sprintf("P%d%s", productId, item.Sku)
}
