package table

import (
	"apify-service/api/helpers"
	"apify-service/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/2751997nam/go-helpers/utils"
	"github.com/gin-gonic/gin"
)

func Store(c *gin.Context) {
	data, err := utils.GetRequestData(c)
	if err != nil {
		utils.ResponseFail(c, "An error occured when parsing request body", http.StatusUnprocessableEntity)
		return
	}
	db := models.GetDB()

	saveData := map[string]any{}
	for key, value := range data {
		if key != "table" {
			saveData[key] = value
		}
	}
	saveData["created_at"] = time.Now()
	saveData["updated_at"] = time.Now()

	err = db.Table(helpers.GetTableName(c.Param("table"))).Create(&saveData).Error
	if err != nil {
		utils.QuickLog(saveData, "", helpers.GetTableName(c.Param("table")), "CREATE_ERROR")
		message := fmt.Sprintf("An error occur when store table %s: "+err.Error(), helpers.GetTableName(c.Param("table")))
		utils.ResponseFail(c, message, http.StatusInternalServerError)
		utils.LogPanic(message)
	} else {
		utils.QuickLog(saveData, "", helpers.GetTableName(c.Param("table")), "CREATE")
		utils.ResponseSuccess(c, saveData, http.StatusOK)
	}
}
