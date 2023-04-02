package table

import (
	"net/http"
	"product-service/internal/models"

	"github.com/2751997nam/go-helpers/utils"
	"github.com/gin-gonic/gin"
)

func Show(c *gin.Context) {
	var result []map[string]any
	id := utils.AnyToUint(c.Param("id"))
	db := models.GetDB()
	err := db.Table(c.Param("table")).Where("id = ?", id).Find(&result).Error

	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusNotFound)
		return
	}
	if len(result) > 0 {
		utils.ResponseSuccess(c, result[0], http.StatusOK)
	} else {
		utils.ResponseSuccess(c, nil, http.StatusOK)
	}
}
