package listcategory

import (
	"database/sql"
	"net/http"
	"product-service/internal/models"
	"strings"

	"github.com/2751997nam/go-helpers/utils"
	"github.com/gin-gonic/gin"
)

type CategoryTreeItem struct {
	ID       int           `json:"id"`
	Name     string        `json:"name,omitempty"`
	Slug     string        `json:"slug,omitempty"`
	Lft      int           `json:"_lft,omitempty" gorm:"column:_lft"`
	Rgt      int           `json:"_rgt,omitempty" gorm:"column:_rgt"`
	ParentId sql.NullInt64 `json:"parent_id,omitempty"`
	// Title    string        `json:"title"`
}

func GetTree(c *gin.Context) {
	db := models.GetDB()
	categories := []models.Category{}
	db.Where("is_hidden = ?", 0).Select("ID", "Name", "Slug", "Lft", "Rgt", "ParentId").Find(&categories)
	retVal := []CategoryTreeItem{}
	categoryById := map[int]models.Category{}
	for _, category := range categories {
		categoryById[category.ID] = category
	}
	for _, category := range categories {
		item := CategoryTreeItem{
			ID:       category.ID,
			Name:     category.Name,
			Slug:     category.Slug,
			Lft:      category.Lft,
			Rgt:      category.Rgt,
			ParentId: category.ParentId,
		}
		names := []string{category.Name}
		if category.ParentId.Valid && category.ParentId.Int64 > 0 {
			paths := []models.Category{}
			GetRootPathWithCache(int(category.ParentId.Int64), &paths, &categoryById)
			for _, path := range paths {
				names = append([]string{path.Name}, names...)
			}
			item.Name = strings.Join(names, " > ")
		}
		retVal = append(retVal, item)
	}

	utils.ResponseSuccess(c, retVal, http.StatusOK)
}

func GetRootNode(cateId int) models.Category {
	db := models.GetDB()
	retVal := models.Category{}

	db.Where("id = ?", cateId).Select("ID", "Name", "Slug", "Lft", "Rgt", "ParentId").First(&retVal)
	if retVal.ParentId.Valid && retVal.ParentId.Int64 > 0 {
		return GetRootNode(int(retVal.ParentId.Int64))
	}

	return retVal
}

func GetRootPath(cateId int, retVal *[]models.Category) {
	db := models.GetDB()
	category := models.Category{}
	db.Where("id = ?", cateId).Select("ID", "Name", "Slug", "Lft", "Rgt", "ParentId").First(&category)
	*retVal = append([]models.Category{category}, (*retVal)...)
	if category.ParentId.Valid && category.ParentId.Int64 > 0 {
		GetRootPath(int(category.ParentId.Int64), retVal)
	}
}

func GetRootPathWithCache(cateId int, retVal *[]models.Category, cache *map[int]models.Category) {
	if category, ok := (*cache)[cateId]; ok {
		*retVal = append([]models.Category{category}, (*retVal)...)
		if category.ParentId.Valid && category.ParentId.Int64 > 0 {
			GetRootPathWithCache(int(category.ParentId.Int64), retVal, cache)
		}
	}
}
