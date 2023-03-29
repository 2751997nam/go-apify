package storecategory

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	listcategory "product-service/api/handlers/category/list_category"
	"product-service/internal/models"

	"github.com/2751997nam/go-helpers/utils"
	"github.com/gin-gonic/gin"
	sluggify "github.com/gosimple/slug"
	nestedset "github.com/longbridgeapp/nested-set"
)

func buildData(c *gin.Context) (models.Category, error) {
	retVal := models.Category{}
	data, err := utils.GetRequestData(c)
	if err != nil {
		return retVal, err
	}
	retVal.ID = int(utils.GetInput("id", data, float64(0)))
	Name := utils.GetInput("name", data, "")
	retVal.Name = Name
	retVal.Slug = utils.GetInput("slug", data, sluggify.Make(Name))
	retVal.Type = utils.GetInput("type", data, "PRODUCT")
	retVal.Description = utils.GetInput("description", data, "")
	retVal.ImageUrl = utils.GetInput("image_url", data, "")
	retVal.BigImageUrl = utils.GetInput("big_image_url", data, "")
	parentId := uint(utils.GetInput("parent_id", data, float64(0)))
	retVal.ParentId = sql.NullInt64{
		Int64: int64(parentId),
		Valid: true,
	}
	retVal.IsHidden = int(utils.GetInput("is_hidden", data, float64(0)))

	return retVal, nil
}

func Store(c *gin.Context) {
	category, err := buildData(c)
	if err != nil {
		log.Println(err)
		utils.ResponseFail(c, "Something went wrong", http.StatusUnprocessableEntity)
	}
	db := models.GetDB()
	if (category.ID) > 0 {
		err = db.Omit("CreatedAt").Save(&category).Error
	} else {
		err = db.Save(&category).Error
	}
	if err != nil {
		utils.ResponseFail(c, "Something went wrong", http.StatusInternalServerError)
		utils.LogPanic(err)
	}
	if category.ParentId.Valid && category.ParentId.Int64 > 0 {
		parent := models.Category{}
		db.Where("id = ?", category.ParentId.Int64).First(&parent)
		utils.LogJson("parent", parent)
		if parent.ID > 0 {
			err := nestedset.MoveTo(db, &category, &parent, nestedset.MoveDirectionLeft)
			if err != nil {
				utils.ResponseFail(c, "Something went wrong", http.StatusInternalServerError)
				utils.LogPanic(err)
			}
		}
	} else {
		nestedset.Create(db, &category, nil)
	}
	breadcrumbs := []models.Category{}
	listcategory.GetRootPath(category.ID, &breadcrumbs)
	value, err := json.Marshal(breadcrumbs)
	if err != nil {
		category.Breadcrumb = string(value)
		err := db.Select("Breadcrumb").Save(&category).Error
		if err != nil {
			utils.ResponseFail(c, "Something went wrong", http.StatusInternalServerError)
			utils.LogPanic(err)
		}
	}

	utils.ResponseSuccess(c, category, http.StatusOK)
}
