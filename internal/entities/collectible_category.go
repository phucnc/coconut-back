package entities

import (
	"time"
)

const (
	CollectibleCategory_CollectibleID = "collectible_id"
	CollectibleCategory_CategoryID = "category_id"
	CollectibleCategory_CreatedAt = "created_at"
	CollectibleCategory_UpdatedAt = "updated_at"
	CollectibleCategory_DeletedAt = "deleted_at"
)

type CollectibleCategory struct {
	CollectibleID int64
	CategoryId    int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

func (e *CollectibleCategory) TableName() string {
	return "collectible_category"
}

func (e *CollectibleCategory) FieldMap(...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		"collectible_id": &e.CollectibleID,
		"category_id":    &e.CategoryId,
		"created_at":     &e.CreatedAt,
		"updated_at":     &e.UpdatedAt,
		"deleted_at":     &e.DeletedAt,
	}
	return fieldMap
}

func (e *CollectibleCategory) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		CollectibleCategory_CollectibleID,
		CollectibleCategory_CategoryID,
		CollectibleCategory_CreatedAt,
		CollectibleCategory_UpdatedAt,
		CollectibleCategory_DeletedAt,
	}
	columnValues := []interface{}{
		&e.CollectibleID,
		&e.CategoryId,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.DeletedAt,
	}
	if len(names) == 0 {
		return columnNames, columnValues
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
