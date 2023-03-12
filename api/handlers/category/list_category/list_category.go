package listcategory

import (
	"product-service/internal/models"

	"github.com/samber/lo"
)

func GetTree() {

}

func GetRootNode(cateId int) models.Category {
	db := models.GetDB()
	retVal := models.Category{}

	db.Where("id = ?", cateId).Select("ID", "Name", "Slug", "Lft", "Rgt").First(&retVal)
	if retVal.ParentId.Valid && retVal.ParentId.Int64 > 0 {
		return GetRootNode(int(retVal.ParentId.Int64))
	}

	return retVal
}

func GetChildPath(cateId int, retVal *[]models.Category) {
	db := models.GetDB()
	category := models.Category{}
	db.Where("id = ?", cateId).Select("ID", "Name", "Slug", "Lft", "Rgt").First(&category)
	if category.ParentId.Valid && category.ParentId.Int64 > 0 {
		*retVal = append(*retVal, category)
		GetChildPath(int(category.ParentId.Int64), retVal)
	}

	lo.Reverse(*retVal)
}
