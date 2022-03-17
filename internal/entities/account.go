package entities

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	Account_Id         = "id"
	Account_Username   = "username"
	Account_Address    = "address"
	Account_Avatar     = "avatar"
	Account_Cover      = "cover"
	Account_Info       = "info"
	Account_Twitter    = "twitter"
	Account_Facebook   = "facebook"
	Account_Tiktok     = "tiktok"
	Account_Instagram  = "instagram"
	Account_Status     = "status"
	Account_KYC_Status = "kyc_status"
	Account_CreatedAt  = "created_at"
	Account_UpdatedAt  = "updated_at"
	Account_DeletedAt  = "deleted_at"
)

type Accounts []*Account

func (c Accounts) Names() []string {
	names := make([]string, len(c), 0)
	/*for _, account := range c {
		names = append(names, account.Username)
	}*/
	return names
}

type Account struct {
	Id         int64          `json:"-"`
	Username   sql.NullString `json:"username"`
	Address    string         `json:"address"`
	Avatar     sql.NullString `json:"avatar"`
	Cover      sql.NullString `json:"cover"`
	Info       sql.NullString `json:"info"`
	Twitter    sql.NullString `json:"twitter"`
	Facebook   sql.NullString `json:"facebook"`
	Tiktok     sql.NullString `json:"tiktok"`
	Instagram  sql.NullString `json:"instagram"`
	Status     int            `json:"status"`
	KYC_status int            `json:"kyc_status"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  *time.Time     `json:"-"`
}

func (e *Account) TableName() string {
	return "account"
}

func (e *Account) FieldMap(names ...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		Account_Id:         &e.Id,
		Account_Username:   &e.Username,
		Account_Address:    &e.Address,
		Account_Avatar:     &e.Avatar,
		Account_Cover:      &e.Cover,
		Account_Info:       &e.Info,
		Account_Twitter:    &e.Twitter,
		Account_Facebook:   &e.Facebook,
		Account_Tiktok:     &e.Tiktok,
		Account_Instagram:  &e.Instagram,
		Account_Status:     &e.Status,
		Account_KYC_Status: &e.KYC_status,
		Account_CreatedAt:  &e.CreatedAt,
		Account_UpdatedAt:  &e.UpdatedAt,
		Account_DeletedAt:  &e.DeletedAt,
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

func (e *Account) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		Account_Id, Account_Username, Account_Address, Account_Info, Account_Avatar, Account_Cover,
		Account_Twitter, Account_Facebook, Account_Instagram, Account_Tiktok, Account_Status, Account_CreatedAt, Account_UpdatedAt, Account_DeletedAt,
	}
	columnValues := []interface{}{
		&e.Id, &e.Username, &e.Address, &e.Info, &e.Avatar, &e.Cover, &e.Twitter, &e.Facebook, &e.Instagram, &e.Tiktok, &e.Status, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
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

func (e *Account) String() string {
	return fmt.Sprintf("%v", *e)
}
