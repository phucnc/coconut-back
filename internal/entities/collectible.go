package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

const (
	Collectible_Id                  = "id"
	Collectible_GUID                = "guid"
	Collectible_Title               = "title"
	Collectitle_Description         = "description"
	Collectible_UploadFile          = "upload_file"
	Collectible_Royalties           = "royalties"
	Collectible_UnlockOncePurchased = "unlock_once_purchased"
	Collectible_InstantSalePrice    = "instant_sale_price"
	Collectible_Properties          = "properties"
	Collectible_TokenId             = "token_id"
	Collectible_TokenOwner          = "token_owner"
	Collectible_Creator             = "creator"
	Collectible_View                = "view"
	Collectible_TotalLike           = "total_like"
	Collectible_Status              = "status"
	Collectible_Token               = "token"
	Collectible_QuoteTokenId        = "quote_token_id"
	Collectible_CreatedAt           = "created_at"
	Collectible_UpdatedAt           = "updated_at"
	Collectible_DeletedAt           = "deleted_at"
)

type Collectible struct {
	Id                  int64               `json:"-"`
	GUID                uuid.UUID           `json:"id"`
	Title               string              `json:"title"`
	Description         string              `json:"description"`
	UploadFile          string              `json:"upload_file"`
	Royalties           decimal.Decimal     `json:"royalties"`
	InstantSalePrice    decimal.Decimal     `json:"instant_sale_price"`
	UnlockOncePurchased bool                `json:"unlock_once_purchased"`
	Properties          map[string]string   `json:"properties"`
	View                int                 `json:"view"`
	TotoalLike          int                 `json:"total_like"`
	Token               *string             `json:"token"`
	TokenId             decimal.NullDecimal `json:"token_id"`
	TokenOwner          *string             `json:"token_owner"`
	Status              int                 `json:"status"`
	Creator             string              `json:"creator"`
	CreatedAt           time.Time           `json:"-"`
	UpdatedAt           time.Time           `json:"-"`
	DeletedAt           *time.Time          `json:"-"`

	QuoteToken  *Token               `json:"quote_token"`
	Categories  []*Category          `json:"categories"`
	Creator_acc *Account             `json:"creator_acc"`
	Owner       *Account             `json:"owner"`
	Like        *CollectibleLikeResp `json:"like"`
}

func (e *Collectible) TableName() string {
	return "collectible"
}

func (e *Collectible) FieldMap(names ...string) map[string]interface{} {
	fieldMap := map[string]interface{}{
		Collectible_Id:                  &e.Id,
		Collectible_Title:               &e.Title,
		Collectitle_Description:         &e.Description,
		Collectible_UploadFile:          &e.UploadFile,
		Collectible_Royalties:           &e.Royalties,
		Collectible_UnlockOncePurchased: &e.UnlockOncePurchased,
		Collectible_InstantSalePrice:    &e.InstantSalePrice,
		Collectible_Properties:          &e.Properties,
		Collectible_Creator:             &e.Creator,
		Collectible_View:                &e.View,
		Collectible_TotalLike:           &e.TotoalLike,
		Collectible_Status:              &e.Status,
		Collectible_CreatedAt:           &e.CreatedAt,
		Collectible_UpdatedAt:           &e.UpdatedAt,
		Collectible_DeletedAt:           &e.DeletedAt,
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

func (e *Collectible) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		Collectible_Id,
		Collectible_GUID,
		Collectible_Title,
		Collectitle_Description,
		Collectible_UploadFile,
		Collectible_Royalties,
		Collectible_UnlockOncePurchased,
		Collectible_InstantSalePrice,
		Collectible_Properties,
		Collectible_TokenId,
		Collectible_TokenOwner,
		Collectible_Creator,
		Collectible_View,
		Collectible_TotalLike,
		Collectible_Status,
		Collectible_Token,
		Collectible_CreatedAt,
		Collectible_UpdatedAt,
		Collectible_DeletedAt,
	}
	columnValues := []interface{}{
		&e.Id,
		&e.GUID,
		&e.Title,
		&e.Description,
		&e.UploadFile,
		&e.Royalties,
		&e.UnlockOncePurchased,
		&e.InstantSalePrice,
		&e.Properties,
		&e.TokenId,
		&e.TokenOwner,
		&e.Creator,
		&e.View,
		&e.TotoalLike,
		&e.Status,
		&e.Token,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.DeletedAt,
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
