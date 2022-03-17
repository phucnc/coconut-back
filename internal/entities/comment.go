package entities

import (
	"time"
)

const (
	Comment_Id        = "id"
	Comment_CollectibleID = "collectible_id"
	Comment_AccountID = "account_id"
	Comment_Content = "content"
	Comment_CreatedAt = "created_at"
	Comment_UpdatedAt = "updated_at"
	Comment_DeletedAt = "deleted_at"
)

type Comment struct {
	Id int64    `json:"id"`
	CollectibleID    int64   `json:"collectible_id"`
	AccountID  int64  `json:"account_id"`
	Account  *Account 	`json:"account"`
	Content string   `json:"content"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     *time.Time  `json:"-"`
}

type Comments []*Comment

func (e *Comment) TableName() string {
	return "comment"
}

func (e *Comment) FieldMap(...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		"id"           : &e.Id,
		"collectible_id": &e.CollectibleID,
		"account_id":    &e.AccountID,
		"content":    	&e.Content,
		"created_at":     &e.CreatedAt,
		"updated_at":     &e.UpdatedAt,
		"deleted_at":     &e.DeletedAt,
	}
	return fieldMap
}

func (e *Comment) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
 		Comment_Id, 
		Comment_CollectibleID,
		Comment_AccountID,
		Comment_Content,
		Comment_CreatedAt,
		Comment_UpdatedAt,
		Comment_DeletedAt,
	}
	columnValues := []interface{}{
		&e.Id,
		&e.CollectibleID,
		&e.AccountID,
		&e.Content,
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
