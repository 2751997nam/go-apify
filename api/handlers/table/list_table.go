package table

import (
	"apify-service/api/helpers"
	"apify-service/internal/models"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/2751997nam/go-helpers/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getFilter(c *gin.Context) map[string]string {
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

func buildQuery(table string, filter map[string]string) *gorm.DB {
	db := models.GetDB()
	query := db.Table(table)

	if filterStr, ok := filter["filters"]; ok && len(filterStr) > 0 {
		buildFilter(query, filterStr)
	}

	if fieldsStr, ok := filter["fields"]; ok {
		buildSelect(query, fieldsStr)
	}

	if sortStr, ok := filter["sorts"]; ok && len(sortStr) > 0 {
		buildSort(query, sortStr)
	}

	return query
}

func buildPaginationQuery(query *gorm.DB, filter map[string]string) *gorm.DB {
	pageSize, _ := strconv.Atoi(getFilterValue("page_size", filter, "50"))
	pageId, _ := strconv.Atoi(getFilterValue("page_id", filter, "0"))
	query.Limit(pageSize).Offset(pageSize * pageId)

	return query
}

func buildFilter(query *gorm.DB, filterStr string) *gorm.DB {
	utils.Log("filterStr", filterStr)
	filters := strings.Split(filterStr, ",")
	for _, str := range filters {
		regex := regexp.MustCompile(`(\w+)(\W{1,3})(.+)`)
		values := regex.FindStringSubmatch(str)
		utils.Log("values", values)
		if len(values) > 3 {
			key, operator, value := values[1], values[2], values[3]
			switch operator {
			case "~":
				query.Where(fmt.Sprintf("`%s` like ?", key), "%"+value+"%")
			case "!~":
				query.Where(fmt.Sprintf("`%s` not like ?", key), "%"+value+"%")
			case "={":
				query.Where(fmt.Sprintf("`%s` In ?", key), strings.Split(value[:len(value)-1], ":"))
			case "!={":
				query.Where(fmt.Sprintf("`%s` Not In ?", key), strings.Split(value[:len(value)-1], ":"))
			default:
				query.Where(fmt.Sprintf("`%s` %s ?", key, operator), value)
			}
		}
	}
	return query
}

func buildSelect(query *gorm.DB, fieldsStr string) *gorm.DB {
	fields := strings.Split(fieldsStr, ",")
	log.Println("fields", fields)
	query.Select(fields)

	return query
}

func buildSort(query *gorm.DB, sortStr string) *gorm.DB {
	values := strings.Split(sortStr, ",")

	for _, value := range values {
		if value[0:1] == "-" {
			query.Order(fmt.Sprintf(`%s desc`, value[1:]))
		} else {
			query.Order(value)
		}
	}

	return query
}

func buildMeta(table string, filter map[string]string) utils.Meta {
	query := buildQuery(table, filter)
	var meta utils.Meta
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
	var result []map[string]any
	filter := getFilter(c)
	query := buildQuery(helpers.GetTableName(c.Param("table")), filter)
	buildPaginationQuery(query, filter).Find(&result)
	meta := buildMeta(helpers.GetTableName(c.Param("table")), filter)

	utils.ResponseWithMeta(c, result, meta)
}
