package storecategory

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	listcategory "product-service/api/handlers/category/list_category"
	"product-service/internal/models"

	goHelpers "github.com/2751997nam/go-helpers/pkg/helpers"

	"github.com/gin-gonic/gin"
	sluggify "github.com/gosimple/slug"
	nestedset "github.com/longbridgeapp/nested-set"
)

func buildData(c *gin.Context) (models.Category, error) {
	retVal := models.Category{}
	data, err := goHelpers.GetRequestData(c)
	if err != nil {
		return retVal, err
	}
	retVal.ID = int(goHelpers.GetInput("id", data, float64(0)))
	Name := goHelpers.GetInput("name", data, "")
	retVal.Name = Name
	retVal.Slug = goHelpers.GetInput("slug", data, sluggify.Make(Name))
	retVal.Type = goHelpers.GetInput("type", data, "PRODUCT")
	retVal.Description = goHelpers.GetInput("description", data, "")
	retVal.ImageUrl = goHelpers.GetInput("image_url", data, "")
	retVal.BigImageUrl = goHelpers.GetInput("big_image_url", data, "")
	parentId := uint(goHelpers.GetInput("parent_id", data, float64(0)))
	retVal.ParentId = sql.NullInt64{
		Int64: int64(parentId),
		Valid: true,
	}
	retVal.IsHidden = int(goHelpers.GetInput("is_hidden", data, float64(0)))

	return retVal, nil
}

func Store(c *gin.Context) {
	category, err := buildData(c)
	if err != nil {
		log.Println(err)
		goHelpers.ResponseFail(c, "Something went wrong", http.StatusUnprocessableEntity)
	}
	db := models.GetDB()
	if (category.ID) > 0 {
		err = db.Omit("CreatedAt").Save(&category).Error
	} else {
		err = db.Save(&category).Error
	}
	if err != nil {
		goHelpers.ResponseFail(c, "Something went wrong", http.StatusInternalServerError)
		goHelpers.LogPanic(err)
	}
	if category.ParentId.Valid && category.ParentId.Int64 > 0 {
		parent := models.Category{}
		db.Where("id = ?", category.ParentId.Int64).First(&parent)
		goHelpers.LogJson("parent", parent)
		if parent.ID > 0 {
			err := nestedset.MoveTo(db, &category, &parent, nestedset.MoveDirectionLeft)
			if err != nil {
				goHelpers.ResponseFail(c, "Something went wrong", http.StatusInternalServerError)
				goHelpers.LogPanic(err)
			}
		}
	} else {
		nestedset.Create(db, &category, nil)
	}
	breadcrumbs := []models.Category{}
	listcategory.GetChildPath(category.ID, &breadcrumbs)
	value, err := json.Marshal(breadcrumbs)
	if err != nil {
		category.Breadcrumb = string(value)
		err := db.Select("Breadcrumb").Save(&category).Error
		if err != nil {
			goHelpers.ResponseFail(c, "Something went wrong", http.StatusInternalServerError)
			goHelpers.LogPanic(err)
		}
	}

	goHelpers.ResponseSuccess(c, category, http.StatusOK)
}
