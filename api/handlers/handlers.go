package handlers

import (
	"net/http"

	goHelpers "github.com/2751997nam/go-helpers/pkg/helpers"

	"github.com/gin-gonic/gin"
)

func DoNothing(c *gin.Context) {

}

func Home(c *gin.Context) {
	goHelpers.ResponseSuccess(c, "ChilleTee Product Service", http.StatusOK)
}
