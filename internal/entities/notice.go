package entities

import (
	"database/sql"
	"time"
)

const (
	Notice_Id            = "id"
	Notice_AccountID     = "account_id"
	Notice_FromAccountID = "from_account"
	Notice_CollectibleID = "collectible_id"
	Notice_Content       = "content"
	Notice_Status        = "status"
	Notice_CreatedAt     = "created_at"
	Notice_UpdatedAt     = "updated_at"
	Notice_DeletedAt     = "deleted_at"
)

const (
	Notice_someone_bought_nft = "Someone bought your NFT"
	Notice_bought_nft         = "You have successfully purchased 1 NFT"
	Notice_someone_comment    = "Someone commented on your NFT"
	Notice_like_nft           = "Someone liked your NFT"
	Notice_create_nft         = "You have successfully created 1 NFT"
	Notice_report_send        = "You have successfully reported 1 NFT, we will review the correctness of this report and notify you of the report results"
	Notice_report_receive     = "1 NFT created by you has been reported, we will review the correctness of this report and notify you of the report's results (The reporting reasons are not related to copyright)"
	Notice_report_report      = "1 NFT created by you has been reported for copyright reasons, if you do not want your NFT to be removed from conteNFT, please contact conteNFT's support at support@contenft.com to provide information and verify ownership of that content"
	Notice_report_kyc         = "You have successfully verified the KYC process, now you have a 'KYC Verified' mark"
	Notice_resell_nft         = "Someone has posted 1 NFT created by you for re-sell"
	Notice_resell_nft_done    = "1 NFT created by you has been successfully resold"
	//Notice_resell_nft         = "You have posted 1 NFT for re-sell, when someone buys 1 NFT owned by you we will send you a notification"
)

type Notice struct {
	Id            int64         `json:"id"`
	AccountID     int64         `json:"account_id"`
	Account       *Account      `json:"account"`
	FromAccountID sql.NullInt64 `json:"-"`
	FromAccount   *Account      `json:"from_account"`
	CollectibleID sql.NullInt64 `json:"-"`
	Collectible   *Collectible  `json:"collectibe"`
	Content       string        `json:"content"`
	Status        int           `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeletedAt     *time.Time    `json:"-"`
}

type Notices []*Notice

func (e *Notice) TableName() string {
	return "notice"
}

func (e *Notice) FieldMap(...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		"id":         &e.Id,
		"account_id": &e.AccountID,
		"content":    &e.Content,
		"status":     &e.Status,
		"created_at": &e.CreatedAt,
		"updated_at": &e.UpdatedAt,
		"deleted_at": &e.DeletedAt,
	}
	return fieldMap
}

func (e *Notice) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		Notice_Id,
		Notice_AccountID,
		Notice_FromAccountID,
		Notice_CollectibleID,
		Notice_Content,
		Notice_Status,
		Notice_CreatedAt,
		Notice_UpdatedAt,
		Notice_DeletedAt,
	}
	columnValues := []interface{}{
		&e.Id,
		&e.AccountID,
		&e.FromAccountID,
		&e.CollectibleID,
		&e.Content,
		&e.Status,
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
