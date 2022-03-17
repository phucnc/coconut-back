package entities

import (
	"fmt"
	"time"
)

const (
	Token_Id        = "id"
	Token_Name      = "name"
	Token_CreatedAt = "created_at"
	Token_UpdatedAt = "updated_at"
	Token_DeletedAt = "deleted_at"
)

type Token struct {
	ID        int32      `json:"-"`
	Name      string     `json:"name"`
	CreateAt  time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func (e *Token) TableName() string {
	return "token"
}

func (e *Token) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		Token_Id, Token_Name, Token_CreatedAt, Token_UpdatedAt, Token_DeletedAt,
	}
	columnValues := []interface{}{
		&e.ID, &e.Name, &e.CreateAt, &e.UpdatedAt, &e.DeletedAt,
	}
	if len(names) == 0 {
		return columnNames, columnValues
	}
	search := make(map[string]bool)
	for _, name := range names {
		search[name] = true
	}
	filteredNames := make([]string, 0, len(names))
	filteredValues := make([]interface{}, 0, len(names))
	for i, columnName := range columnNames {
		_, ok := search[columnName]
		if ok {
			filteredNames = append(filteredNames, columnNames[i])
			filteredValues = append(filteredValues, columnValues[i])
		}
	}
	return filteredNames, filteredValues
}

func (e *Token) String() string {
	return fmt.Sprintf("%+v", *e)
}