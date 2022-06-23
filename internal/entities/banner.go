package entities

import (
	//"database/sql"
	"fmt"
	"time"
)

const (
	Banner_Id        = "id"
	Banner_Name      = "name"
	Banner_Status    = "status"
	Banner_Order     = "in_order"
	Banner_Picture   = "picture"
	Banner_Link      = "link"
	Banner_CreatedAt = "created_at"
	Banner_UpdatedAt = "updated_at"
	Banner_DeletedAt = "deleted_at"
)

type Banners []*Banner

/*func (c Accounts) Names() []string {
names := make([]string, len(c), 0)
/*for _, account := range c {
	names = append(names, account.Username)
}*/
//return names
//}

type Banner struct {
	Id        int64      `json:"-"`
	Name      string     `json:"name"`
	Link      string     `json:"link"`
	Status    int        `json:"status"`
	Picture   string     `json:"picture"`
	Order     int        `json:"order"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func (e *Banner) TableName() string {
	return "banner"
}

func (e *Banner) FieldMap(names ...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		Banner_Id:         &e.Id,
		Banner_Name:       &e.Name,
		Banner_Link:       &e.Link,
		Banner_Picture:    &e.Picture,
		Banner_Status:     &e.Status,
		Banner_Order:      &e.Order,
		Account_CreatedAt: &e.CreatedAt,
		Account_UpdatedAt: &e.UpdatedAt,
		Account_DeletedAt: &e.DeletedAt,
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

func (e *Banner) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		Banner_Id, Banner_Name, Banner_Picture, Banner_Link, Banner_Status, Banner_Order,
		Banner_CreatedAt, Banner_UpdatedAt, Banner_DeletedAt,
	}
	columnValues := []interface{}{
		&e.Id, &e.Name, &e.Picture, &e.Link, &e.Status, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
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

func (e *Banner) String() string {
	return fmt.Sprintf("%v", *e)
}
