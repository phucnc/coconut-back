package database

import (
	"fmt"
)

type Entity interface {
	TableName() string
	FieldsAndValues(...string) ([]string, []interface{})
}

func GeneratePlaceholder(n int) []string {
	placeholders := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%v", i))
	}
	return placeholders
}
