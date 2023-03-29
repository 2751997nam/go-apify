package handlers

import (
	"net/http"

	"github.com/2751997nam/go-helpers/utils"
	"github.com/gin-gonic/gin"
)

func DoNothing(c *gin.Context) {

}

func Home(c *gin.Context) {
	utils.ResponseSuccess(c, "ChilleTee Product Service", http.StatusOK)
}
