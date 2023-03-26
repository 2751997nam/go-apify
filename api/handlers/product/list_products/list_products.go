package listproducts

import (
	"math"
	"product-service/internal/helpers"
	"product-service/internal/models"
	timelayout "product-service/internal/pkg/time_layout"
	"product-service/internal/types"
	"regexp"
	"strconv"
	"strings"
	"time"

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

	status := getFilterValue("status", filter, "")
	if len(status) > 0 {
		query.Where("status = ?", status)

	}

	searchQuery := getFilterValue("search", filter, "")
	searchQuery = strings.ToLower(searchQuery)
	if len(searchQuery) > 0 {
		if regexp.MustCompile(`^p(\d+)`).MatchString(searchQuery) {
			searchQuery = searchQuery[0:]
			var productIds []uint64
			db.Model(&models.Product{}).Where("id = ?", searchQuery).Pluck("id", &productIds)

			query.Where("id in ?", productIds)
		}
		if helpers.IsNumeric(searchQuery) {
			var productIds []uint64
			db.Model(&models.Product{}).Where("id = ?", searchQuery).Pluck("id", &productIds)
			if len(productIds) > 0 {
				query.Where("id in ?", productIds)
			} else {
				query.Where("sku = ?", searchQuery)
			}
		} else {
			query.Where("name like ?", "%"+searchQuery+"%")
		}
	}
	query.Select("sb_product.id", "sb_product.name", "Sku", "sb_product.slug", "sb_product.image_url", "Price", "HighPrice", "Status", "ActorId", "UpdaterId", "sb_product.created_at", "sb_product.updated_at")

	categoryId := getFilterValue("category_id", filter, "")
	if len(categoryId) > 0 {
		query.Joins("join sb_product_n_category as spc on sb_product.`id` = spc.`product_id`")
		query.Where("spc.category_id = ?", categoryId)
	}
	createrId := getFilterValue("actor_id", filter, "")
	modifierId := getFilterValue("modifier_id", filter, "")
	if len(createrId) > 0 {
		query.Where("actor_id = ?", createrId)
	}
	if len(modifierId) > 0 {
		query.Where("updater_id = ?", modifierId)
	}
	fromDate := getFilterValue("from", filter, "")
	if len(fromDate) > 0 {
		from, err := time.Parse(timelayout.DateOnly, fromDate)
		if err == nil {
			query.Where("sb_product.created_at >= ?", time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location()))
		}
	}

	toDate := getFilterValue("to", filter, "")
	helpers.Log("toDate", toDate)
	if len(toDate) > 0 {
		to, err := time.Parse(timelayout.DateOnly, toDate)
		helpers.Log("toDate", to)
		if err == nil {
			query.Where("sb_product.created_at <= ?", time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999, to.Location()))
		}
	}

	query.Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})
	query.Preload("Creater").Preload("Modifier")
	query.Order("id desc")

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
