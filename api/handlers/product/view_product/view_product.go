package viewproduct

import (
	"net/http"
	"product-service/internal/helpers"
	"product-service/internal/models"
	"product-service/internal/types"

	"github.com/gin-gonic/gin"
)

type GetSkuReponse struct {
	Variants        []types.VariantView    `json:"variants"`
	ProductVariants []types.ProductSkuView `json:"productVariants"`
}

type ViewResult struct {
	Product models.Product `json:"product,omitempty"`
	// Category   models.Category `json:"category,omitempty"`
	// FeatureTag models.Tag      `json:"feature_tag,omitempty"`
	Tags     []models.Tag  `json:"tags,omitempty"`
	Variants GetSkuReponse `json:"variants"`
}

func getProduct(id uint64) models.Product {
	db := models.GetDB()
	product := models.Product{
		BaseModel: models.BaseModel{
			ID: id,
		},
	}

	db.Model(&product).Preload("Tags").Preload("Gallery", "type", "PRODUCT").Find(&product)

	return product
}

func View(c *gin.Context) {
	productId := helpers.AnyFloat64ToUint64(c.Param("id"))
	if productId == 0 {
		helpers.ResponseFail(c, "product id is required", http.StatusUnprocessableEntity)
		return
	}

	product := getProduct(productId)

	result := ViewResult{
		Product:  product,
		Variants: GetSkues(product),
	}

	helpers.ResponseSuccess(c, result, http.StatusOK)

}

func ViewVariant(c *gin.Context) {
	productId := helpers.AnyFloat64ToUint64(c.Param("id"))
	if productId == 0 {
		helpers.ResponseFail(c, "product id is required", http.StatusUnprocessableEntity)
		return
	}

	db := models.GetDB()
	product := models.Product{
		BaseModel: models.BaseModel{
			ID: productId,
		},
	}

	db.Model(&product).Select("ID", "ImageUrl", "Status").Find(&product)

	// product := getProduct(productId)

	result := GetSkues(product)

	helpers.ResponseSuccess(c, result, http.StatusOK)
}
