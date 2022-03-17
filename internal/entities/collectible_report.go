package entities

import (
	"time"
)

const (
	CollectibleReport_Id        = "id"
	CollectibleReport_CollectibleID = "collectible_id"
	CollectibleReport_AccountID = "account_id"
	CollectibleReport_ReportTypeID = "report_type_id"
	CollectibleReport_Content = "content"
	CollectibleReport_CreatedAt = "created_at"
	CollectibleReport_UpdatedAt = "updated_at"
	CollectibleReport_DeletedAt = "deleted_at"
)

type CollectibleReport struct {
	Id int64
	CollectibleID    int64
	AccountID  int64
	ReportTypeID int
	Content string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

func (e *CollectibleReport) TableName() string {
	return "collectible_report"
}

func (e *CollectibleReport) FieldMap(...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		"collectible_id": &e.CollectibleID,
		"account_id":    &e.AccountID,
		"report_type_id":    &e.ReportTypeID,
		"content":    &e.Content,
		"created_at":     &e.CreatedAt,
		"updated_at":     &e.UpdatedAt,
		"deleted_at":     &e.DeletedAt,
	}
	return fieldMap
}

func (e *CollectibleReport) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		CollectibleReport_Id,
		CollectibleReport_CollectibleID,
		CollectibleReport_AccountID,
		CollectibleReport_ReportTypeID,
		CollectibleReport_Content,
		CollectibleReport_CreatedAt,
		CollectibleReport_UpdatedAt,
		CollectibleReport_DeletedAt,
	}
	columnValues := []interface{}{
		&e.Id,
		&e.CollectibleID,
		&e.AccountID,
		&e.ReportTypeID,
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
