package entities

import (
	"fmt"
	"time"
)

const (
	Category_Id        = "id"
	Category_Name      = "name"
	Category_CreatedAt = "created_at"
	Category_UpdatedAt = "updated_at"
	Category_DeletedAt = "deleted_at"
)

type Categories []*Category

func (c Categories) Names() []string {
	names := make([]string, len(c), 0)
	for _, category := range c {
		names = append(names, category.Name)
	}
	return names
}

type Category struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreateAt  time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func (e *Category) TableName() string {
	return "category"
}

func (e *Category) FieldMap(names ...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		Category_Id:        &e.ID,
		Category_Name:      &e.Name,
		Category_CreatedAt: &e.CreateAt,
		Category_UpdatedAt: &e.UpdatedAt,
		Category_DeletedAt: &e.DeletedAt,
	}
	if len(names) == 0 {
		return fieldMap
	} else {
		optionalFields := make(map[string]interface{})
		for _, name := range names {
			optionalFields[name] = fieldMap[name]
		}
		return optionalFields
	}
}

func (e *Category) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		Category_Id, Category_Name, Category_CreatedAt, Category_UpdatedAt, Category_DeletedAt,
	}
	columnValues := []interface{}{
		&e.ID, &e.Name, &e.CreateAt, &e.CreateAt, &e.DeletedAt,
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

func (e *Category) String() string {
	return fmt.Sprintf("%v", *e)
}
