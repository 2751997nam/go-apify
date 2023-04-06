package helpers

import "fmt"

func GetTableName(table string) string {
	return fmt.Sprintf(`sb_%s`, table)
}
