package entities

import (
	"fmt"
	"github.com/shopspring/decimal"
)

const (
	SellTokenEvent_BlockNumber    = "block_number"
	SellTokenEvent_BlockHash      = "block_hash"
	SellTokenEvent_BlockTimestamp = "block_timestamp"
	SellTokenEvent_TxIndex        = "tx_index"
	SellTokenEvent_TxHash         = "tx_hash"
	SellTokenEvent_NFT_TokenId    = "nft_token_id"
	SellTokenEvent_NFT_Price      = "nft_price"
	SellTokenEvent_NFT_QuoteToken = "nft_quote_token"
	SellTokenEvent_LogIndex       = "log_index"
	SellTokenEvent_LogAddress     = "log_address"
	SellTokenEvent_LogData        = "log_data"
	SellTokenEvent_LogRemoved     = "log_removed"
	SellTokenEvent_Type           = "type"
	SellTokenEvent_Account        = "account"
)

type SellTokenEvent struct {
	TokenId  decimal.Decimal
	NFTPrice NFTPrice

	Block Block
	Tx    Tx
	Log   Log

	Type    int
	Account string
}

func (e *SellTokenEvent) TableName() string {
	return "exchange_event_token"
}

func (e *SellTokenEvent) FieldsAndValues(names ...string) ([]string, []interface{}) {
	columnNames := []string{
		SellTokenEvent_BlockNumber,
		SellTokenEvent_BlockHash,
		SellTokenEvent_BlockTimestamp,
		SellTokenEvent_TxIndex,
		SellTokenEvent_TxHash,
		SellTokenEvent_NFT_TokenId,
		SellTokenEvent_NFT_Price,
		SellTokenEvent_NFT_QuoteToken,
		SellTokenEvent_LogIndex,
		SellTokenEvent_LogAddress,
		SellTokenEvent_LogData,
		SellTokenEvent_LogRemoved,
		SellTokenEvent_Type,
		SellTokenEvent_Account,
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

func (e *SellTokenEvent) String() string {
	return fmt.Sprintf("%+v", *e)
}
