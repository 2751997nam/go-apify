package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"product-service/internal/types"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess(c *gin.Context, result any, status int) {
	c.IndentedJSON(status, types.Response{
		Status: "successful",
		Result: result,
	})
}

func ResponseSuccessWithMessage(c *gin.Context, result any, message string) {
	c.IndentedJSON(http.StatusOK, types.Response{
		Status:  "successful",
		Result:  result,
		Message: message,
	})
}

func ResponseWithMeta(c *gin.Context, result any, meta types.Meta) {
	c.IndentedJSON(http.StatusOK, types.Response{
		Status: "successful",
		Result: result,
		Meta:   meta,
	})
}

func ResponseFail(c *gin.Context, message string, status int) {
	c.IndentedJSON(status, types.Response{
		Status:  "fail",
		Message: message,
	})
}

func AnyToString(value any) string {
	if value != nil {
		return strings.Trim(fmt.Sprint(value), " ")
	}

	return ""
}

func AnyToInt(value any) int {
	result, _ := strconv.Atoi(AnyToString(value))
	return result
}

func AnyToUint(value any) uint64 {
	result, _ := strconv.ParseInt(AnyToString(value), 10, 64)
	fmt.Println("AnyToUint", value, uint64(result))
	return uint64(result)
}

func AnyFloat64ToUint64(value any) uint64 {
	var result float64 = 0
	if reflect.TypeOf(value).Name() == "string" {
		result, _ = strconv.ParseFloat(AnyToString(value), 64)
	} else {
		result = value.(float64)
	}

	return uint64(result)
}

func AnyToFloat(value any) float32 {
	result, _ := strconv.ParseFloat(AnyToString(value), 64)
	return float32(result)
}

func ArrayChunk[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func Join[T any](items []T, sep string) string {
	strs := []string{}
	for _, item := range items {
		strs = append(strs, fmt.Sprint(item))
	}

	return strings.Join(strs, sep)
}

func IsNumeric(str string) bool {
	return regexp.MustCompile(`\d+`).MatchString(str)
}

func GetUrlFields(url string) []string {
	retVal := []string{}
	regex := *regexp.MustCompile(`\/:(\w+)($|\/)`)
	res := regex.FindAllStringSubmatch(url, -1)
	if len(res) > 0 && len(res[0]) > 1 {
		retVal = append(retVal, res[0][1])
	}
	return retVal
}

func GetRequestBody(c *gin.Context) (map[string]any, error) {
	var data map[string]any
	bodyAsByteArray, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(bodyAsByteArray), &data); err != nil {
		return nil, err
	}

	fields := GetUrlFields(c.FullPath())

	for _, field := range fields {
		if IsNumeric(c.Param(field)) {
			value, _ := strconv.ParseFloat(c.Param(field), 64)
			data[field] = value
		} else {
			data[field] = c.Param(field)
		}
	}

	query := c.Request.URL.Query()
	for field, value := range query {
		data[field] = value
	}

	return data, nil
}

func GetInput[T any](key string, data map[string]any, defaultValue T) T {
	value, ok := data[key]
	if ok {
		return value.(T)
	}

	return defaultValue
}

func LogJson(prefix string, data any) {
	str, _ := json.MarshalIndent(data, "", "\t")
	log.Println(prefix, string(str))
}

func LogPanic(data any) {
	log.Panic(data)
}
