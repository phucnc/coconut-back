package entities

import (
	//"database/sql"
	"fmt"
	"time"
)

const (
	KYC_Id         = "id"
	KYC_Account_Id = "account_id"
	KYC_Fullname   = "fullname"
	KYC_Email      = "email"
	KYC_Birthday   = "birthday"
	KYC_City       = "city"
	KYC_Country    = "country"
	KYC_FrontId    = "front_id"
	KYC_BackId     = "back_id"
	KYC_Selfienote = "selfienote"
	KYC_Status     = "status"
	KYC_CreatedAt  = "created_at"
	KYC_UpdatedAt  = "updated_at"
	KYC_DeletedAt  = "deleted_at"
)

type KYCs []*Account

func (c KYCs) Names() []string {
	names := make([]string, len(c), 0)
	/*for _, account := range c {
		names = append(names, account.Username)
	}*/
	return names
}

type KYC struct {
	Id         int64  `json:"-"`
	Account_id int64  `json:"accountid"`
	Fullname   string `json:"fullname"`
	Email      string `json:"email"`
	Birthday   string `json:"birthday"`
	City       string `json:"city"`
	Country    string `json:"country"`
	FrontId    string `json:"frontid"`
	BackId     string `json:"backid"`
	Selfienote string `json:"selfienote"`
	Status     int    `json:"status"`

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func (e *KYC) TableName() string {
	return "kyc"
}

func (e *KYC) FieldMap(names ...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		KYC_Id:         &e.Id,
		KYC_Fullname:   &e.Fullname,
		KYC_Email:      &e.Email,
		KYC_Birthday:   &e.Birthday,
		KYC_City:       &e.City,
		KYC_Country:    &e.Country,
		KYC_FrontId:    &e.FrontId,
		KYC_BackId:     &e.BackId,
		KYC_Selfienote: &e.Selfienote,
		KYC_Status:     &e.Status,
		KYC_Account_Id: &e.Account_id,

		KYC_CreatedAt: &e.CreatedAt,
		KYC_UpdatedAt: &e.UpdatedAt,
		KYC_DeletedAt: &e.DeletedAt,
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

func (e *KYC) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		KYC_Id, KYC_Fullname, KYC_Account_Id, KYC_Email, KYC_Birthday, KYC_FrontId,
		KYC_BackId, KYC_Selfienote, KYC_City, KYC_Country, KYC_Status, KYC_CreatedAt, KYC_UpdatedAt, KYC_DeletedAt,
	}
	columnValues := []interface{}{
		&e.Id, &e.Fullname, &e.Account_id, &e.Email, &e.Birthday, &e.FrontId,
		&e.BackId, &e.Selfienote, &e.City, &e.Country, &e.Status, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
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

func (e *KYC) String() string {
	return fmt.Sprintf("%v", *e)
}
