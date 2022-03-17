package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"strconv"
	"strings"
	"time"
)

type ExchangeEventRepo struct{}

/*func (r *ExchangeEventRepo) UpsertBuyTokenEvent(ctx context.Context, db database.QueryExecer, event *entities.BuyTokenEvent) error {
	names, values := event.FieldsAndValues()
	_, err := database.InsertIgnoreConflict(ctx, db, event, names, values)
	if err != nil {
		return errors.Wrap(err, "database.Insert()")
	}
	return nil
}*/

func (r *ExchangeEventRepo) UpsertTokenEvent(ctx context.Context, db database.QueryExecer, event *entities.TokenEvent) error {
	names, values := event.FieldsAndValues()
	_, err := database.InsertIgnoreConflict(ctx, db, event, names, values)
	if err != nil {
		return errors.Wrap(err, "database.Insert()")
	}
	return nil
}

func (r *ExchangeEventRepo) GetTotalSoldInDuration(ctx context.Context, db database.QueryExecer, duration time.Duration, quoteToken int16) (decimal.Decimal, error) {
	buyTokenEvent := &entities.TokenEvent{}

	stmt :=
		`
		SELECT
			coalesce(sum(%s), 0)
		FROM
			%s
		WHERE
		    exchange_event_token.type = 1 
			AND %s >= now() - INTERVAL '%s hours'
			AND %s = $1
		`

	stmt = fmt.Sprintf(
		stmt,
		entities.TokenEvent_NFT_Price,
		buyTokenEvent.TableName(),
		entities.TokenEvent_BlockTimestamp,
		strconv.FormatInt(int64(duration.Hours()), 10),
		entities.TokenEvent_NFT_QuoteToken,
	)

	row := db.QueryRow(ctx, stmt, quoteToken)

	var totalSold decimal.Decimal

	switch err := row.Scan(&totalSold); err {
	case nil:
		return totalSold, nil
	case pgx.ErrNoRows:
		return decimal.Zero, nil
	default:
		return decimal.Zero, errors.Wrap(err, "db.Query")
	}

	/*buyTokenEvents := make([]*entities.BuyTokenEvent, 0)
	for rows.Next() {
		buyTokenEvent := &entities.BuyTokenEvent{}
		_, values := buyTokenEvent.FieldsAndValues()
		err := rows.Scan(values...)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		buyTokenEvents = append(buyTokenEvents, buyTokenEvent)
	}
	return buyTokenEvents, nil*/
}

func (r *ExchangeEventRepo) GetHistory(ctx context.Context, db database.QueryExecer, tokenId decimal.NullDecimal) ([]*entities.TokenEvent, error) {

	stmt :=
		`
		SELECT
			%s
		FROM
			%s
	     JOIN
			token ON token.id = exchange_event_token.nft_quote_token
			where nft_token_id = $1
			order by exchange_event_token.created_at desc
		`

	event := &entities.TokenEvent{}
	columnNames, columnValues := event.FieldsAndValuesGet()
	for i, columnName := range columnNames {
		columnNames[i] = event.TableName() + "." + columnName
	}
	//columnValues = append(columnValues, &event.Id)

	quoteToken := &entities.Token{}
	quoteTokenColumnNames, _ := quoteToken.FieldsAndValues()
	for i, quoteTokenColumnName := range quoteTokenColumnNames {
		quoteTokenColumnNames[i] = quoteToken.TableName() + "." + quoteTokenColumnName
	}
	columnNames = append(columnNames, quoteTokenColumnNames...)

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		event.TableName(),
	)

	//fmt.Println(stmt)

	rows, err := db.Query(ctx, stmt, tokenId)
	switch err {
	case nil:

		events := make([]*entities.TokenEvent, 0)
		for rows.Next() {
			event := &entities.TokenEvent{}
			columnNames, columnValues = event.FieldsAndValuesGet()
			for i, columnName := range columnNames {
				columnNames[i] = event.TableName() + "." + columnName
			}

			quoteToken := &entities.Token{}
			_, quoteTokenColumnValues := quoteToken.FieldsAndValues()
			columnValues = append(columnValues, quoteTokenColumnValues...)

			err := rows.Scan(columnValues...)
			if err != nil {
				return nil, err
			}
			event.Price = event.NFTPrice.Price.Div(decimal.NewFromFloat(1000000000000000000))
			event.QuoteToken = quoteToken
			events = append(events, event)
		}

		fmt.Printf("%v", event)

		return events, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *ExchangeEventRepo) CountTrade(ctx context.Context, db database.QueryExecer, token_id decimal.Decimal) int {
	var total int

	stmt :=
		`
		SELECT
			count(*) 
		FROM
			exchange_event_token
		WHERE
			type = 1 AND nft_token_id = $1
		`

		/*	stmt = fmt.Sprintf(
			stmt,
			"",
		)*/

	row := db.QueryRow(ctx, stmt, token_id)
	err := row.Scan(&total)

	switch err {
	case nil:
		return total
	default:
		return 0
	}
}
