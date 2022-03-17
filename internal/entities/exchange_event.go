package entities

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type Block struct {
	Number    uint64
	Hash      string
	Timestamp time.Time
}

type Tx struct {
	Index uint   `json:"-"`
	Hash  string `json:"tx"`
}

type Log struct {
	Index   uint
	Address string
	Data    []byte
	Removed bool
}

type NFTPrice struct {
	Token int16           `json:"price_quote_token"`
	Price decimal.Decimal `json:"price"`
}

const (
	TokenEvent_BlockNumber    = "block_number"
	TokenEvent_BlockHash      = "block_hash"
	TokenEvent_BlockTimestamp = "block_timestamp"
	TokenEvent_TxIndex        = "tx_index"
	TokenEvent_TxHash         = "tx_hash"
	TokenEvent_NFT_TokenId    = "nft_token_id"
	TokenEvent_NFT_Price      = "nft_price"
	TokenEvent_NFT_QuoteToken = "nft_quote_token"
	TokenEvent_LogIndex       = "log_index"
	TokenEvent_LogAddress     = "log_address"
	TokenEvent_LogData        = "log_data"
	TokenEvent_LogRemoved     = "log_removed"
	TokenEvent_Type           = "type"
	TokenEvent_Account        = "account"
	TokenEvent_Sell           = 0
	TokenEvent_Buy            = 1
	TokenEvent_CreatedAt      = "created_at"
)

type TokenEvent struct {
	TokenId    decimal.Decimal
	NFTPrice   NFTPrice `json:"-"`
	QuoteToken *Token   `json:"quote_token"`

	Block     Block `json:"-"`
	Tx        Tx
	Log       Log `json:"-"`
	Type      int
	Price     decimal.Decimal `json:"price"`
	Account   *string         `json:"-"`
	From      *Account        `json:"account"`
	CreatedAt time.Time       `json:"created_at"`
}

func (e *TokenEvent) TableName() string {
	return "exchange_event_token"
}

func (e *TokenEvent) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		TokenEvent_BlockNumber,
		TokenEvent_BlockHash,
		TokenEvent_BlockTimestamp,
		TokenEvent_TxIndex,
		TokenEvent_TxHash,
		TokenEvent_NFT_TokenId,
		TokenEvent_NFT_Price,
		TokenEvent_NFT_QuoteToken,
		TokenEvent_LogIndex,
		TokenEvent_LogAddress,
		TokenEvent_LogData,
		TokenEvent_LogRemoved,
		TokenEvent_Type,
		TokenEvent_Account,
	}
	columnValues := []interface{}{
		&e.Block.Number,
		&e.Block.Hash,
		&e.Block.Timestamp,
		&e.Tx.Index,
		&e.Tx.Hash,
		&e.TokenId,
		&e.NFTPrice.Price,
		&e.NFTPrice.Token,
		&e.Log.Index,
		&e.Log.Address,
		&e.Log.Data,
		&e.Log.Removed,
		&e.Type,
		&e.Account,
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

func (e *TokenEvent) FieldsAndValuesGet(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		TokenEvent_BlockNumber,
		TokenEvent_BlockHash,
		TokenEvent_BlockTimestamp,
		TokenEvent_TxIndex,
		TokenEvent_TxHash,
		TokenEvent_NFT_TokenId,
		TokenEvent_NFT_Price,
		TokenEvent_NFT_QuoteToken,
		TokenEvent_LogIndex,
		TokenEvent_LogAddress,
		TokenEvent_LogData,
		TokenEvent_LogRemoved,
		TokenEvent_Type,
		TokenEvent_Account,
		TokenEvent_CreatedAt,
	}
	columnValues := []interface{}{
		&e.Block.Number,
		&e.Block.Hash,
		&e.Block.Timestamp,
		&e.Tx.Index,
		&e.Tx.Hash,
		&e.TokenId,
		&e.NFTPrice.Price,
		&e.NFTPrice.Token,
		&e.Log.Index,
		&e.Log.Address,
		&e.Log.Data,
		&e.Log.Removed,
		&e.Type,
		&e.Account,
		&e.CreatedAt,
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

func (e *TokenEvent) String() string {
	return fmt.Sprintf("%+v", *e)
}
