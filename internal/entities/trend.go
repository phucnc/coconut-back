package entities

import (
	"time"
)

const (
	Trend_Id = "Id"
	Trend_CollectibleID = "collectible_id"
	Trend_In_Order = "in_order"
	Trend_Advertisement = "adsvertisement"
	Trend_CreatedAt = "created_at"
)

type Trend struct {
	Id int64 `json:"id"`
	CollectibleID int64 `json:"collectible_id"`
	Collectible      *Collectible     `json:"collectible"`
	In_Order int64 `json:"order"`
	Advertisement bool `json:"advertisement"`
	CreatedAt     time.Time
}



type Trends []*Trend

type TrendResp struct {
	CollectibleID bool 		`json:"Collectible_id"`
}


func (e *Trend) TableName() string {
	return "trend"
}

func (e *Trend) FieldMap(...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		"collectible_id": &e.CollectibleID,
		"order": &e.In_Order,
		"adsvertisement": &e.Advertisement,
		"created_at":     &e.CreatedAt,
	}
	return fieldMap
}

func (e *Trend) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		Trend_Id,
		Trend_CollectibleID,
		Trend_In_Order,
		Trend_Advertisement,
		Trend_CreatedAt,
	}
	columnValues := []interface{}{
		&e.Id,
		&e.CollectibleID,
		&e.In_Order,
		&e.Advertisement,
		&e.CreatedAt,
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
