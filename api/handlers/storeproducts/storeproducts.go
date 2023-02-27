package storeproducts

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"product-service/internal/helpers"

	"github.com/gin-gonic/gin"
)

func Store(c *gin.Context) {
	var data map[string]any
	bodyAsByteArray, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
		helpers.ResponseFail(c, "something went wrong", http.StatusUnprocessableEntity)
		return
	}
	if err = json.Unmarshal([]byte(bodyAsByteArray), &data); err != nil {
		log.Fatal(err)
		helpers.ResponseFail(c, "something went wrong", http.StatusUnprocessableEntity)
		return
	}
	message, ok := validate(data)
	if !ok {
		helpers.ResponseFail(c, message, http.StatusUnprocessableEntity)
		return
	}

	productData := buildProductData(data)

	helpers.ResponseSuccess(c, productData, http.StatusAccepted)
}
