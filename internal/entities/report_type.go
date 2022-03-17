package entities

import (
	"fmt"
	"time"
)

const (
	ReportType_Id        = "id"
	ReportType_Description      = "description"
	ReportType_CreatedAt = "created_at"
	ReportType_UpdatedAt = "updated_at"
)

type ReportTypes []*ReportType

func (c ReportTypes) Names() []string {
	names := make([]string, len(c), 0)
	for _, reporttype := range c {
		names = append(names, reporttype.Description)
	}
	return names
}

type ReportType struct {
	Id        int64      `json:"id"`
	Description      string     `json:"description"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

func (e *ReportType) TableName() string {
	return "report_type"
}

func (e *ReportType) FieldMap(names ...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		ReportType_Id:        &e.Id,
		ReportType_Description:      &e.Description,
		ReportType_CreatedAt: &e.CreatedAt,
		ReportType_UpdatedAt: &e.UpdatedAt,
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

func (e *ReportType) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		ReportType_Id, ReportType_Description, ReportType_CreatedAt, ReportType_UpdatedAt,
	}
	columnValues := []interface{}{
		&e.Id, &e.Description, &e.CreatedAt, &e.UpdatedAt,
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

func (e *ReportType) String() string {
	return fmt.Sprintf("%v", *e)
}
