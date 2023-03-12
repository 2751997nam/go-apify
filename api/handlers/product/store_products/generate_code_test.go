package storeproducts

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func Test_GenerateOptionCode(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	code := GenerateOptionCode("white")
	fmt.Println(code)
}
