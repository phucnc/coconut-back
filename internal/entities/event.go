package entities

import (
	"fmt"
	"time"
	//"database/sql"
)

const (
	Event_Id        = "id"
	Event_Title      = "title"
	Event_Banner      = "banner"
	Event_Content      = "content"
	Event_Status     = "status"
	Event_CreatedAt = "created_at"
	Event_UpdatedAt = "updated_at"
	Event_DeletedAt = "deleted_at"
)

type Events []*Event

func (c Events) Names() []string {
	names := make([]string, len(c), 0)
	/*for _, account := range c {
		names = append(names, account.Username)
	}*/
	return names
}

type Event struct {
	Id        int64      `json:"id"`
	Title      string     `json:"title"`
	Banner      string       `json:"banner"`
	Content      string      `json:"content"`
	Status      int      `json:"status"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func (e *Event) TableName() string {
	return "event"
}

func (e *Event) FieldMap(names ...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		Event_Id:        &e.Id,
		Event_Title:      &e.Title,
		Event_Banner:      &e.Banner,
		Event_Content:      &e.Content,
		Event_Status:      &e.Status,
		Event_CreatedAt: &e.CreatedAt,
		Event_UpdatedAt: &e.UpdatedAt,
		Event_DeletedAt: &e.DeletedAt,
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

func (e *Event) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		Event_Id, Event_Title, Event_Banner, Event_Content , Event_Status, 
		Event_CreatedAt, Event_UpdatedAt, Event_DeletedAt,
	}
	columnValues := []interface{}{
		&e.Id, &e.Title,&e.Banner, &e.Content, &e.Status, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
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

func (e *Event) String() string {
	return fmt.Sprintf("%v", *e)
}
