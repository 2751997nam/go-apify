package table

import (
	"apify-service/internal/models"
	"fmt"
	"net/http"

	"github.com/2751997nam/go-helpers/utils"
	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	id := utils.AnyToUint(c.Param("id"))
	db := models.GetDB()
	err := db.Table(c.Param("table")).Where("id = ?", id).Delete(&map[string]any{}).Error
	if err != nil {
		message := fmt.Sprintf("An error occur when delete table %s: "+err.Error(), c.Param("table"))
		utils.ResponseFail(c, message, http.StatusInternalServerError)
		utils.LogPanic(message)
	} else {
		utils.ResponseSuccess(c, nil, http.StatusOK)
	}
}
