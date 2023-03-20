package listproducts

import (
	"math"
	"product-service/internal/helpers"
	"product-service/internal/models"
	"product-service/internal/types"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func buildFilter(c *gin.Context) map[string]string {
	params := c.Request.URL.Query()
	filter := map[string]string{}
	for key, values := range params {
		filter[key] = values[0]
	}

	return filter
}

func getFilterValue(key string, filter map[string]string, defaultValue string) string {
	result := defaultValue

	if val, ok := filter[key]; ok {
		result = val
	}

	return result
}

func buildQuery(filter map[string]string) *gorm.DB {
	db := models.GetDB()
	query := db.Model(&models.Product{})

	status := getFilterValue("status", filter, "ACTIVE")

	query.Where("status = ?", status)

	return query
}

func buildPaginationQuery(query *gorm.DB, filter map[string]string) *gorm.DB {
	pageSize, _ := strconv.Atoi(getFilterValue("page_size", filter, "50"))
	pageId, _ := strconv.Atoi(getFilterValue("page_id", filter, "0"))
	query.Limit(pageSize).Offset(pageSize * pageId)

	return query
}

func buildMeta(filter map[string]string) types.Meta {
	query := buildQuery(filter)
	var meta types.Meta
	pageSize, _ := strconv.Atoi(getFilterValue("page_size", filter, "50"))
	pageId, _ := strconv.Atoi(getFilterValue("page_id", filter, "0"))
	meta.PageId = pageId
	meta.PageSize = pageSize
	query.Count(&meta.TotalCount)
	if pageSize != 0 {
		meta.PageCount = int(math.Ceil(float64(meta.TotalCount) / float64(pageSize)))
	}
	if pageId < meta.PageCount-1 {
		meta.HasNext = true
	}

	return meta
}

func Find(c *gin.Context) {
	var products []models.Product
	filter := buildFilter(c)
	query := buildQuery(filter)
	buildPaginationQuery(query, filter).Find(&products)
	meta := buildMeta(filter)
	helpers.ResponseWithMeta(c, products, meta)
}
