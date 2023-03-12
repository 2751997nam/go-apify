package handlers

import (
	"net/http"
	"product-service/internal/helpers"

	"github.com/gin-gonic/gin"
)

func DoNothing(c *gin.Context) {

}

func Home(c *gin.Context) {
	helpers.ResponseSuccess(c, nil, http.StatusOK)
}
