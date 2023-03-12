package storeproducts

import (
	"fmt"
	"math/rand"
	"product-service/internal/models"
	"regexp"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

func GenerateOptionCode(optionName string) string {
	db := models.GetDB()
	/*
	   Đếm số từ: $countWord
	   $countWord >= 3 => Lấy các chữ cái đầu của tên
	   $countWord = 2 => Lấy 2 chữ cái từ đầu + 1 chữ cái từ 2
	   $countWord = 1 => Lấy 3 chữ cái từ đầu

	   Nếu đã trùng + thêm ASCII cho tới khi không trùng
	   Nếu hết bảng dùng gencode auto incre redis
	*/
	var rex = regexp.MustCompile("[^A-Za-z0-9 ]")
	var rexSpace = regexp.MustCompile("\\s+")
	name := rex.ReplaceAllString(optionName, "")
	name = rexSpace.ReplaceAllString(name, " ")
	name = strings.Trim(name, " ")
	words := strings.Split(name, " ")
	acronym := ""
	if len(words) >= 3 {
		for k, v := range words {
			acronym += v[0:1]
			if k == 2 {
				break
			}
		}
	} else if len(words) == 2 {
		acronym = words[0]
		if len(words[0]) >= 2 {
			acronym = words[0][0:2]
		}
		acronym += words[1][0:1]
	} else if len(words) == 1 {
		acronym = name
		if len(words[0]) >= 3 {
			acronym = words[0][0:3]
		}
	}
	acronym = strings.ToUpper(acronym)
	attrCode := acronym
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	try := 0
	for {
		option := models.VariantOption{}
		db.Where("code = ?", attrCode).First(&option)
		if option.ID == 0 {
			break
		}
		attrCode = acronym + "_" + alphabet[try:try+1]
		try++
		if try >= len(alphabet) {
			acronym += "A"
			try = 0
		}
		if len(acronym) > 6 {
			break
		}
	}

	return attrCode
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateSkuCode(productId uint64, options []any, limit int) string {
	var retVal []string = []string{fmt.Sprint("P", productId)}
	db := models.GetDB()
	for _, v := range options {
		item := v.(map[string]any)
		code, ok := item["code"].(string)
		if !ok {
			optionSlug, ok := item["slug"].(string)
			if !ok {
				optionSlug = slug.Make(item["name"].(string))
			}
			option := models.VariantOption{
				Slug: optionSlug,
			}
			db.Order("id desc").First(&option)
			if len(option.Code) > 0 {
				retVal = append(retVal, option.Code)
			}
		} else {
			retVal = append(retVal, code)
		}
	}

	microtime := fmt.Sprint(time.Now().UnixMicro())
	str := RandStringRunes(3) + microtime[:len(microtime)-3] + RandStringRunes(3)
	sublen := limit - len(str) - 1

	return strings.Join(retVal, "-")[0:sublen] + str
}
