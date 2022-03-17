package entities

import (
	"time"
)

const (
	CommentLike_commentID = "comment_id"
	CommentLike_AccountID = "account_id"
	CommentLike_CreatedAt = "created_at"
	CommentLike_UpdatedAt = "updated_at"
)

type CommentLike struct {
	CommentID int64
	AccountID    int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (e *CommentLike) TableName() string {
	return "collectible_like"
}

func (e *CommentLike) FieldMap(...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		"comment_id": &e.CommentID,
		"account_id":    &e.AccountID,
		"created_at":     &e.CreatedAt,
		"updated_at":     &e.UpdatedAt,
	}
	return fieldMap
}

func (e *CommentLike) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		CommentLike_commentID,
		CommentLike_AccountID,
		CommentLike_CreatedAt,
		CommentLike_UpdatedAt,
	}
	columnValues := []interface{}{
		&e.CommentID,
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
