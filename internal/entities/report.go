package entities

import (
	"fmt"
	"time"
)

const (
	Report_Id        = "id"
	Report_Collectible_Id      = "collectible_id"
	Report_Account_Id      = "account_id"
	Report_Type_Id     = "report_type_id"
	Report_Description     = "description"
	Report_Status    = "status"
	Report_CreatedAt = "created_at"
	Report_UpdatedAt = "updated_at"
)

type Reports []*Report

/*func (c Reports) Names() []string {
	names := make([]string, len(c), 0)
	for _, report := range c {
		names = append(names, report.Report_Type_Id)
	}
	return names
}*/

type Report struct {
	Id        int64      `json:"id"`
	Account_Id      int64     `json:"account_id"`
	Account      *Account     `json:"account"`
	Collectible_Id      int64     `json:"collectible_id"`
	Collectible      *Collectible     `json:"collectible"`
	Report_Type_Id int64 `json:"type"`
	Description string `json:"description"`
	Status int `json:"status"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

func (e *Report) TableName() string {
	return "collectible_report"
}

func (e *Report) FieldMap(names ...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		Report_Id:        &e.Id,
		Report_Collectible_Id:      &e.Collectible_Id,
		Report_Account_Id:      &e.Account_Id,
		Report_Type_Id:      &e.Report_Type_Id,
		Report_Status:      &e.Status,
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

func (e *Report) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		Report_Id, Report_Collectible_Id,Report_Account_Id, Report_Type_Id, Report_Status,  Report_CreatedAt, Report_UpdatedAt,
	}
	columnValues := []interface{}{
		&e.Id, &e.Collectible_Id,&e.Account_Id,  &e.Report_Type_Id ,&e.Status, &e.CreatedAt, &e.UpdatedAt,
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

func (e *Report) String() string {
	return fmt.Sprintf("%v", *e)
}
