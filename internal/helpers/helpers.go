package helpers

import (
	"fmt"
	"net/http"
	"product-service/internal/types"
	"reflect"
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
