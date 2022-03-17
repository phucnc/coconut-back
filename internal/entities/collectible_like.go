package entities

import (
	"time"
)

const (
	CollectibleLike_CollectibleID = "collectible_id"
	CollectibleLike_AccountID = "account_id"
	CollectibleLike_CreatedAt = "created_at"
	CollectibleLike_UpdatedAt = "updated_at"
)

type CollectibleLike struct {
	CollectibleID int64
	AccountID    int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CollectibleLikeResp struct {
	Liked bool 		`json:"liked"`
	Total    int 	`json:"total"`
}


func (e *CollectibleLike) TableName() string {
	return "collectible_like"
}

func (e *CollectibleLike) FieldMap(...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		"collectible_id": &e.CollectibleID,
		"account_id":    &e.AccountID,
		"created_at":     &e.CreatedAt,
		"updated_at":     &e.UpdatedAt,
	}
	return fieldMap
}

func (e *CollectibleLike) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		CollectibleLike_CollectibleID,
		CollectibleLike_AccountID,
		CollectibleLike_CreatedAt,
		CollectibleLike_UpdatedAt,
	}
	columnValues := []interface{}{
		&e.CollectibleID,
		&e.AccountID,
		&e.CreatedAt,
		&e.UpdatedAt,
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
