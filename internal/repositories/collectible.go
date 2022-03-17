package repositories

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"strconv"
	"strings"
	"time"
)

type CollectibleRepository struct {
}

type collectible struct {
}

func (r *CollectibleRepository) Insert(ctx context.Context, db database.QueryExecer, collectible *entities.Collectible) error {
	names, values := collectible.FieldsAndValues(
		entities.Collectible_Title,
		entities.Collectitle_Description,
		entities.Collectible_UploadFile,
		entities.Collectible_Royalties,
		entities.Collectible_UnlockOncePurchased,
		entities.Collectible_InstantSalePrice,
		entities.Collectible_Properties,
		entities.Collectible_Creator,
	)
	names = append(names, entities.Collectible_QuoteTokenId)
	values = append(values, collectible.QuoteToken.ID)
	err := database.InsertReturning(
		ctx,
		db,
		collectible,
		names,
		values,
		[]string{
			entities.Collectible_Id,
			entities.Collectible_GUID,
			entities.Collectible_CreatedAt,
			entities.Collectible_UpdatedAt,
			entities.Collectible_DeletedAt,
		},
		[]interface{}{
			&collectible.Id,
			&collectible.GUID,
			&collectible.CreatedAt,
			&collectible.UpdatedAt,
			&collectible.DeletedAt,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *CollectibleRepository) GetById(ctx context.Context, db database.QueryExecer, id int64) (*entities.Collectible, error) {
	collectible := &entities.Collectible{}

	stmt :=
		`
		SELECT
			%s ,
			array_agg(category.id) ,
			array_agg(category.name) ,
			array_agg(category.created_at) ,
			array_agg(category.updated_at) ,
			array_agg(category.deleted_at)
		FROM
			%s
		JOIN
			collectible_category ON collectible.id = collectible_category.collectible_id
		JOIN
			category ON category.id = collectible_category.category_id
		JOIN
			token ON token.id = collectible.quote_token_id
		WHERE
			collectible.id = $1
			AND collectible.deleted_at IS NULL
			AND collectible.token_id IS NOT NULL
			AND EXISTS ((
				SELECT
					1
				FROM
					exchange_event_token
				WHERE
						exchange_event_token.nft_token_id = collectible.token_id
						AND exchange_event_token.type = 0
				ORDER BY
					exchange_event_token.block_number DESC
				LIMIT
					1
			))


		GROUP BY
			%s
		`

		/*
					AND (
				SELECT
					nft_price
				FROM
					exchange_event_buy_token
				WHERE
					exchange_event_buy_token.nft_token_id = collectible.token_id
				ORDER BY
					exchange_event_buy_token.block_number DESC
				LIMIT
					1
			) IS NULL
		*/

	columnNames, columnValues := collectible.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = collectible.TableName() + "." + columnName
	}

	quoteToken := &entities.Token{}
	quoteTokenColumnNames, quoteTokenColumnValues := quoteToken.FieldsAndValues()
	for i, quoteTokenColumnName := range quoteTokenColumnNames {
		quoteTokenColumnNames[i] = quoteToken.TableName() + "." + quoteTokenColumnName
	}
	columnNames = append(columnNames, quoteTokenColumnNames...)
	columnValues = append(columnValues, quoteTokenColumnValues...)

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		collectible.TableName(),
		strings.Join(columnNames, ", "),
	)

	categoryIds := make([]int64, 0)
	categoryNames := make([]string, 0)
	categoryCreatedAtList := make([]time.Time, 0)
	categoryUpdatedAtList := make([]time.Time, 0)
	categoryDeletedAtList := make([]*time.Time, 0)
	columnValues = append(columnValues, &categoryIds)
	columnValues = append(columnValues, &categoryNames)
	columnValues = append(columnValues, &categoryCreatedAtList)
	columnValues = append(columnValues, &categoryUpdatedAtList)
	columnValues = append(columnValues, &categoryDeletedAtList)

	row := db.QueryRow(ctx, stmt, id)
	err := row.Scan(columnValues...)

	//fmt.Println("%v ", columnValues)

	switch err {
	case nil:
		categories := make([]*entities.Category, 0, len(categoryIds))
		for i, categoryId := range categoryIds {
			category := &entities.Category{
				ID:        categoryId,
				Name:      categoryNames[i],
				CreateAt:  categoryCreatedAtList[i],
				UpdatedAt: categoryUpdatedAtList[i],
				DeletedAt: categoryDeletedAtList[i],
			}
			categories = append(categories, category)
		}
		collectible.Categories = categories
		collectible.QuoteToken = quoteToken
		return collectible, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *CollectibleRepository) Get(ctx context.Context, db database.QueryExecer, guid uuid.UUID) (*entities.Collectible, error) {
	collectible := &entities.Collectible{}

	stmt :=
		`
		SELECT
			%s ,
			array_agg(category.id) ,
			array_agg(category.name) ,
			array_agg(category.created_at) ,
			array_agg(category.updated_at) ,
			array_agg(category.deleted_at)
		FROM
			%s
		JOIN
			collectible_category ON collectible.id = collectible_category.collectible_id
		JOIN
			category ON category.id = collectible_category.category_id
		JOIN
			token ON token.id = collectible.quote_token_id
		WHERE
			collectible.guid = $1
			AND collectible.deleted_at IS NULL
			AND collectible.token_id IS NOT NULL
			AND EXISTS ((
				SELECT
					1
				FROM
					exchange_event_token
				WHERE
						exchange_event_token.nft_token_id = collectible.token_id
						AND  exchange_event_token.type = 0 
				ORDER BY
					exchange_event_token.block_number DESC
				LIMIT
					1
			))

		GROUP BY
			%s
		`

		/*
					AND (
				SELECT
					nft_price
				FROM
					exchange_event_buy_token
				WHERE
					exchange_event_buy_token.nft_token_id = collectible.token_id
				ORDER BY
					exchange_event_buy_token.block_number DESC
				LIMIT
					1
			) IS NULL
		*/

	columnNames, columnValues := collectible.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = collectible.TableName() + "." + columnName
	}

	quoteToken := &entities.Token{}
	quoteTokenColumnNames, quoteTokenColumnValues := quoteToken.FieldsAndValues()
	for i, quoteTokenColumnName := range quoteTokenColumnNames {
		quoteTokenColumnNames[i] = quoteToken.TableName() + "." + quoteTokenColumnName
	}
	columnNames = append(columnNames, quoteTokenColumnNames...)
	columnValues = append(columnValues, quoteTokenColumnValues...)

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		collectible.TableName(),
		strings.Join(columnNames, ", "),
	)

	categoryIds := make([]int64, 0)
	categoryNames := make([]string, 0)
	categoryCreatedAtList := make([]time.Time, 0)
	categoryUpdatedAtList := make([]time.Time, 0)
	categoryDeletedAtList := make([]*time.Time, 0)
	columnValues = append(columnValues, &categoryIds)
	columnValues = append(columnValues, &categoryNames)
	columnValues = append(columnValues, &categoryCreatedAtList)
	columnValues = append(columnValues, &categoryUpdatedAtList)
	columnValues = append(columnValues, &categoryDeletedAtList)

	row := db.QueryRow(ctx, stmt, guid.String())
	err := row.Scan(columnValues...)

	//fmt.Println("%v ", columnValues)

	switch err {
	case nil:
		categories := make([]*entities.Category, 0, len(categoryIds))
		for i, categoryId := range categoryIds {
			category := &entities.Category{
				ID:        categoryId,
				Name:      categoryNames[i],
				CreateAt:  categoryCreatedAtList[i],
				UpdatedAt: categoryUpdatedAtList[i],
				DeletedAt: categoryDeletedAtList[i],
			}
			categories = append(categories, category)
		}
		collectible.Categories = categories
		collectible.QuoteToken = quoteToken
		return collectible, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *CollectibleRepository) Check(ctx context.Context, db database.QueryExecer, id uuid.UUID) (bool, error) {
	collectible := &entities.Collectible{}
	var total int
	var check bool
	check = false

	stmt :=
		`
		SELECT
			count(*) 
		FROM
			%s
		WHERE
			collectible.guid = $1
			AND collectible.deleted_at IS NULL
		`

	/*columnNames, columnValues := account.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = collectible.TableName() + "." + columnName
	}*/

	stmt = fmt.Sprintf(
		stmt,
		collectible.TableName(),
	)

	row := db.QueryRow(ctx, stmt, id)
	err := row.Scan(&total)

	if total > 0 {
		check = true
	}

	switch err {
	case nil:
		return check, nil
	default:
		return false, err
	}
}

type CollectiblePaging struct {
	Cursor     uuid.UUID
	Filter     string
	Sort       string
	Title      string
	Limit      int
	Categories []string
	Options    string
	Address    string
	Status     int
}

func (r *CollectibleRepository) Paging(ctx context.Context, db database.QueryExecer, paging *CollectiblePaging) ([]*entities.Collectible, error) {

	filterFields := map[string]string{
		"created-date":       "created_at",
		"created_date":       "created_at",
		"instant-sale-price": "instant_sale_price",
		"instant_sale_price": "instant_sale_price",
		"price":              "instant_sale_price",
		"view":               "view",
		"like":               "total_like",
		"trade":              "total_trade",
		//"like": ""
	}
	filterField, ok := filterFields[paging.Filter]
	if !ok {
		return nil, errors.New("filter field  is not valid")
	}

	sortOperators := map[string]string{
		"asc":        ">=",
		"desc":       "<=",
		"ascending":  ">=",
		"descending": "<=",
	}

	sortOperator, ok := sortOperators[paging.Sort]
	if !ok {
		return nil, errors.New("sort field  is not valid")
	}

	sortFields := map[string]string{
		"asc":        "ASC",
		"desc":       "DESC",
		"ascending":  "ASC",
		"descending": "DESC",
	}
	sortField, ok := sortFields[paging.Sort]
	if !ok {
		return nil, errors.New("sort field  is not valid")
	}

	categoriesComparator := "!="
	categories := make([]string, 0)
	if len(paging.Categories) > 0 {
		categories = paging.Categories
		categoriesComparator = "&&"
		for _, category := range categories {
			category = strings.ToLower(category)
		}
	}

	likeTitle := entities.Collectible_Title
	if paging.Title != "" {
		likeTitle = fmt.Sprintf(`'%%%s%%'`, paging.Title)
	}

	var stmt string

	/*lock :=  ` AND (
		SELECT
			nft_price
		FROM
			exchange_event_buy_token
		WHERE
			exchange_event_buy_token.nft_token_id = collectible.token_id
		ORDER BY
			exchange_event_buy_token.block_number DESC
		LIMIT
			1
	) IS NULL `*/

	lock := "AND collectible.unlock_once_purchased = FALSE"

	var args []interface{}

	where := ""
	if paging.Status == -1 {
		where = "AND collectible.status = -1"
	} else {
		where = "AND collectible.status = 0"
	}

	if paging.Address != "" {

		switch paging.Options {
		case "sold":
			where += "AND lower(collectible.creator) = lower('" + paging.Address + "') AND lower(collectible.creator) <> lower(collectible.token_owner)"
			lock = ""
		case "bought":
			where += "AND lower(collectible.token_owner) = lower('" + paging.Address + "')   AND collectible.unlock_once_purchased = TRUE"
			lock = ""
		case "owner":
			where += "AND lower(collectible.token_owner) = lower('" + paging.Address + "') AND collectible.unlock_once_purchased = FALSE"
			lock = ""
		case "creator":
			where += "AND lower(collectible.creator) = lower('" + paging.Address + "')"
			lock = ""
		default:
			where = ""
		}
	}

	if paging.Cursor == uuid.Nil {
		stmt =
			`
		SELECT
			%s ,
			array_agg(category.id) ,
			array_agg(category.name) ,
			array_agg(category.created_at) ,
			array_agg(category.updated_at) ,
			array_agg(category.deleted_at)
		FROM
			%s
		JOIN
			collectible_category ON collectible.id = collectible_category.collectible_id
		JOIN
			category ON category.id = collectible_category.category_id
		JOIN
			token ON token.id = collectible.quote_token_id
		WHERE
			collectible.deleted_at IS NULL
			%s
			AND collectible.token_id IS NOT NULL
			AND EXISTS ((
				SELECT
					1
				FROM
					exchange_event_token
				WHERE
						exchange_event_token.nft_token_id = collectible.token_id
						AND exchange_event_token.type = 0 
				ORDER BY
					exchange_event_token.block_number DESC
				LIMIT
					1
			))
            %s
			AND lower(collectible.title) LIKE lower(%s)
		GROUP BY
			%s
		HAVING
    		array_agg(category.name) %s ($1)
		ORDER BY
			%s %s, collectible.id %s
		LIMIT 
			$2
		`

		collectible := &entities.Collectible{}
		columnNames, columnValues := collectible.FieldsAndValues()
		for i, columnName := range columnNames {
			columnNames[i] = collectible.TableName() + "." + columnName
		}
		columnValues = append(columnValues, &collectible.Id)

		quoteToken := &entities.Token{}
		quoteTokenColumnNames, _ := quoteToken.FieldsAndValues()
		for i, quoteTokenColumnName := range quoteTokenColumnNames {
			quoteTokenColumnNames[i] = quoteToken.TableName() + "." + quoteTokenColumnName
		}
		columnNames = append(columnNames, quoteTokenColumnNames...)

		stmt = fmt.Sprintf(
			stmt,
			strings.Join(columnNames, ", "),
			collectible.TableName(),
			where,
			lock,
			likeTitle,
			strings.Join(columnNames, ", "),
			categoriesComparator,
			collectible.TableName()+"."+filterField,
			sortField,
			sortField,
		)
		args = []interface{}{categories, paging.Limit}
		//fmt.Println(stmt)
	} else {
		stmt =
			`
		WITH cte AS (
			SELECT
				collectible.id,
				collectible.created_at,
				collectible.instant_sale_price
			FROM
				collectible
			WHERE
				collectible.guid = $1
			LIMIT 1
		)
		SELECT
			%s ,
			array_agg(category.id) ,
			array_agg(category.name) ,
			array_agg(category.created_at) ,
			array_agg(category.updated_at) ,
			array_agg(category.deleted_at)
		FROM
			%s
		JOIN
			collectible_category ON collectible.id = collectible_category.collectible_id
		JOIN
			category ON category.id = collectible_category.category_id
		JOIN
			token ON token.id = collectible.quote_token_id
		WHERE
			collectible.deleted_at IS NULL
			%s
			AND collectible.token_id IS NOT NULL
			AND collectible.%s %s (SELECT cte.%s FROM cte)
			AND collectible.id %s (SELECT cte.id FROM cte)
			AND EXISTS ((
				SELECT
					1
				FROM
					exchange_event_token
				WHERE
						exchange_event_token.nft_token_id = collectible.token_id
						AND exchange_event_token.type = 0 
				ORDER BY
					exchange_event_token.block_number DESC
				LIMIT
					1
			))

			%s
		
		AND lower(collectible.title) LIKE lower(%s)
		GROUP BY
			%s
		HAVING
    		array_agg(category.name) %s ($2)
		ORDER BY
			%s %s, collectible.id %s
		LIMIT 
			$3
		`

		collectible := &entities.Collectible{}
		columnNames, columnValues := collectible.FieldsAndValues()
		for i, columnName := range columnNames {
			columnNames[i] = collectible.TableName() + "." + columnName
		}
		columnValues = append(columnValues, &collectible.Id)

		quoteToken := &entities.Token{}
		quoteTokenColumnNames, _ := quoteToken.FieldsAndValues()
		for i, quoteTokenColumnName := range quoteTokenColumnNames {
			quoteTokenColumnNames[i] = quoteToken.TableName() + "." + quoteTokenColumnName
		}
		columnNames = append(columnNames, quoteTokenColumnNames...)

		stmt = fmt.Sprintf(
			stmt,
			strings.Join(columnNames, ", "),
			collectible.TableName(),
			where,
			filterField,
			sortOperator,
			filterField,
			strings.ReplaceAll(sortOperator, "=", ""),
			lock,
			likeTitle,
			strings.Join(columnNames, ", "),
			categoriesComparator,
			collectible.TableName()+"."+filterField,
			sortField,
			sortField,
		)
		args = []interface{}{paging.Cursor, categories, paging.Limit}
		//fmt.Println(stmt)
	}

	//fmt.Println(stmt)
	//fmt.Println("lock %s", lock)
	//fmt.Println("where %s", where)

	rows, err := db.Query(ctx, stmt, args...)
	switch err {
	case nil:
		collectibles := make([]*entities.Collectible, 0, paging.Limit)
		for rows.Next() {
			collectible := &entities.Collectible{}
			columnNames, columnValues := collectible.FieldsAndValues()
			for i, columnName := range columnNames {
				columnNames[i] = collectible.TableName() + "." + columnName
			}

			quoteToken := &entities.Token{}
			_, quoteTokenColumnValues := quoteToken.FieldsAndValues()
			/*for i, quoteTokenColumnName := range quoteTokenColumnNames {
				quoteTokenColumnNames[i] = quoteToken.TableName() + "." + quoteTokenColumnName
			}
			columnNames = append(columnNames, quoteTokenColumnNames...)*/
			columnValues = append(columnValues, quoteTokenColumnValues...)

			categoryIds := make([]int64, 0)
			categoryNames := make([]string, 0)
			categoryCreatedAtList := make([]time.Time, 0)
			categoryUpdatedAtList := make([]time.Time, 0)
			categoryDeletedAtList := make([]*time.Time, 0)
			columnValues = append(columnValues, &categoryIds)
			columnValues = append(columnValues, &categoryNames)
			columnValues = append(columnValues, &categoryCreatedAtList)
			columnValues = append(columnValues, &categoryUpdatedAtList)
			columnValues = append(columnValues, &categoryDeletedAtList)
			err := rows.Scan(columnValues...)
			if err != nil {
				return nil, err
			}
			categories := make([]*entities.Category, 0, len(categoryIds))
			for i, categoryId := range categoryIds {
				category := &entities.Category{
					ID:        categoryId,
					Name:      categoryNames[i],
					CreateAt:  categoryCreatedAtList[i],
					UpdatedAt: categoryUpdatedAtList[i],
					DeletedAt: categoryDeletedAtList[i],
				}
				categories = append(categories, category)
			}
			// like
			//creator
			//owner

			collectible.Categories = categories
			collectible.QuoteToken = quoteToken

			collectibles = append(collectibles, collectible)
		}
		return collectibles, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *CollectibleRepository) UpdateTokenInfoByGUID(ctx context.Context, db database.QueryExecer, guid string, tokenId decimal.Decimal, token string, tokenOwner string) error {
	collectible := &entities.Collectible{}

	stmt :=
		`
		UPDATE
			%s
		SET
			%s = $1 , %s = $2, %s = $3
		WHERE
			%s = $4
		`

	stmt = fmt.Sprintf(
		stmt,
		collectible.TableName(),
		entities.Collectible_Token,
		entities.Collectible_TokenId,
		entities.Collectible_TokenOwner,
		entities.Collectible_GUID,
	)

	cmd, err := db.Exec(ctx, stmt, token, tokenId, tokenOwner, guid)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *CollectibleRepository) UpdateLock(ctx context.Context, db database.QueryExecer, tokenId decimal.Decimal, lock bool) error {
	collectible := &entities.Collectible{}

	stmt :=
		`
		UPDATE
			%s
		SET
			%s = $1 
		WHERE
			%s = $2
		`

	stmt = fmt.Sprintf(
		stmt,
		collectible.TableName(),
		entities.Collectible_UnlockOncePurchased,
		entities.Collectible_TokenId,
	)

	cmd, err := db.Exec(ctx, stmt, lock, tokenId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *CollectibleRepository) UpdateView(ctx context.Context, db database.QueryExecer, id int64) error {
	collectible := &entities.Collectible{}

	stmt :=
		`
		UPDATE
			%s
		SET
			view = view+1 
		WHERE
			id = $1
		`

	stmt = fmt.Sprintf(
		stmt,
		collectible.TableName(),
	)

	cmd, err := db.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *CollectibleRepository) UpdateTotalLike(ctx context.Context, db database.QueryExecer, id int64, total_like int) error {
	collectible := &entities.Collectible{}

	stmt :=
		`
		UPDATE
			%s
		SET
			total_like = $1 
		WHERE
			id = $2
		`

	stmt = fmt.Sprintf(
		stmt,
		collectible.TableName(),
	)

	cmd, err := db.Exec(ctx, stmt, total_like, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *CollectibleRepository) UpdateTotalTrade(ctx context.Context, db database.QueryExecer, token_id decimal.Decimal, total_trade int) error {
	collectible := &entities.Collectible{}

	stmt :=
		`
		UPDATE
			%s
		SET
			total_trade = $1 
		WHERE
			token_id = $2
		`

	stmt = fmt.Sprintf(
		stmt,
		collectible.TableName(),
	)

	cmd, err := db.Exec(ctx, stmt, total_trade, token_id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *CollectibleRepository) UpdateResell(ctx context.Context, db database.QueryExecer, id int64,
	price decimal.Decimal, quote_token int32) error {
	collectible := &entities.Collectible{}

	stmt :=
		`
		UPDATE
			%s
		SET
			%s = $2 , %s = $3

		WHERE
			id = $1
		`

	stmt = fmt.Sprintf(
		stmt,
		collectible.TableName(),
		entities.Collectible_InstantSalePrice,
		entities.Collectible_QuoteTokenId,
	)

	//fmt.Println(stmt)

	cmd, err := db.Exec(ctx, stmt, id, price, quote_token)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *CollectibleRepository) UpdateStatus(ctx context.Context, db database.QueryExecer, id int64, status int) error {
	collectible := &entities.Collectible{}

	stmt :=
		`
		UPDATE
			%s
		SET
			status = $1 
		WHERE
			id = $2
		`

	stmt = fmt.Sprintf(
		stmt,
		collectible.TableName(),
	)

	cmd, err := db.Exec(ctx, stmt, status, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *CollectibleRepository) GetByName(ctx context.Context, db database.QueryExecer, name string, limit int) ([]*entities.Collectible, error) {
	collectible := &entities.Collectible{}

	stmt :=
		`
		SELECT
			%s ,
			array_agg(category.id) ,
			array_agg(category.name) ,
			array_agg(category.created_at) ,
			array_agg(category.updated_at) ,
			array_agg(category.deleted_at)
		FROM
			%s
		JOIN
			collectible_category ON collectible.id = collectible_category.collectible_id
		JOIN
			category ON category.id = collectible_category.category_id
		JOIN
			token ON token.id = collectible.quote_token_id
		WHERE
			collectible.title LIKE '%%%s%%'
			AND collectible.deleted_at IS NULL
			AND collectible.token_id IS NOT NULL
			AND EXISTS ((
				SELECT
					1
				FROM
					exchange_event_token
				WHERE
						exchange_event_token.nft_token_id = collectible.token_id
						AND exchange_event_token.type = 0 
				ORDER BY
					exchange_event_token.block_number DESC
				LIMIT
					1
			))

		GROUP BY
			%s
		ORDER BY
			collectible.title
		LIMIT
			$1
		`

		/*
					AND (
				SELECT
					nft_price
				FROM
					exchange_event_buy_token
				WHERE
					exchange_event_buy_token.nft_token_id = collectible.token_id
				ORDER BY
					exchange_event_buy_token.block_number DESC
				LIMIT
					1
			) IS NULL
		*/

	columnNames, columnValues := collectible.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = collectible.TableName() + "." + columnName
	}

	quoteToken := &entities.Token{}
	quoteTokenColumnNames, quoteTokenColumnValues := quoteToken.FieldsAndValues()
	for i, quoteTokenColumnName := range quoteTokenColumnNames {
		quoteTokenColumnNames[i] = quoteToken.TableName() + "." + quoteTokenColumnName
	}
	columnNames = append(columnNames, quoteTokenColumnNames...)
	columnValues = append(columnValues, quoteTokenColumnValues...)

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		collectible.TableName(),
		name,
		strings.Join(columnNames, ", "),
	)
	//fmt.Println(stmt)

	categoryIds := make([]int64, 0)
	categoryNames := make([]string, 0)
	categoryCreatedAtList := make([]time.Time, 0)
	categoryUpdatedAtList := make([]time.Time, 0)
	categoryDeletedAtList := make([]*time.Time, 0)
	columnValues = append(columnValues, &categoryIds)
	columnValues = append(columnValues, &categoryNames)
	columnValues = append(columnValues, &categoryCreatedAtList)
	columnValues = append(columnValues, &categoryUpdatedAtList)
	columnValues = append(columnValues, &categoryDeletedAtList)

	rows, err := db.Query(ctx, stmt, limit)
	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}

	collectibles := make([]*entities.Collectible, 0)
	for rows.Next() {
		collectible := &entities.Collectible{}
		_, columnValues := collectible.FieldsAndValues()

		quoteToken := &entities.Token{}
		_, quoteTokenColumnValues := quoteToken.FieldsAndValues()
		columnValues = append(columnValues, quoteTokenColumnValues...)

		categoryIds := make([]int64, 0)
		categoryNames := make([]string, 0)
		categoryCreatedAtList := make([]time.Time, 0)
		categoryUpdatedAtList := make([]time.Time, 0)
		categoryDeletedAtList := make([]*time.Time, 0)
		columnValues = append(columnValues, &categoryIds)
		columnValues = append(columnValues, &categoryNames)
		columnValues = append(columnValues, &categoryCreatedAtList)
		columnValues = append(columnValues, &categoryUpdatedAtList)
		columnValues = append(columnValues, &categoryDeletedAtList)

		err := rows.Scan(columnValues...)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan()")
		}

		categories := make([]*entities.Category, 0, len(categoryIds))
		for i, categoryId := range categoryIds {
			category := &entities.Category{
				ID:        categoryId,
				Name:      categoryNames[i],
				CreateAt:  categoryCreatedAtList[i],
				UpdatedAt: categoryUpdatedAtList[i],
				DeletedAt: categoryDeletedAtList[i],
			}
			categories = append(categories, category)
		}
		collectible.Categories = categories
		collectible.QuoteToken = quoteToken

		collectibles = append(collectibles, collectible)
	}

	return collectibles, nil
}

func (r *CollectibleRepository) GetTotalInDuration(ctx context.Context, db database.QueryExecer, duration time.Duration) (decimal.Decimal, error) {
	collectible := &entities.Collectible{}

	var stmt string

	if duration > 0 {

		stmt =
			`
		SELECT
			count(*)
		FROM
			%s
		WHERE
			%s >= now() - INTERVAL '%s hours'
			AND collectible.deleted_at IS NULL
		`

		stmt = fmt.Sprintf(
			stmt,
			collectible.TableName(),
			entities.Collectible_CreatedAt,
			strconv.FormatInt(int64(duration.Hours()), 10),
		)
	} else {

		stmt =
			`
				SELECT
					count(*)
				FROM
					%s
				WHERE collectible.deleted_at IS NULL
				`
		stmt = fmt.Sprintf(
			stmt,
			collectible.TableName(),
		)
	}

	row := db.QueryRow(ctx, stmt)

	var total decimal.Decimal

	switch err := row.Scan(&total); err {
	case nil:
		return total, nil
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

func (r *CollectibleRepository) Delete(ctx context.Context, db database.QueryExecer, id int64) error {
	collectible := &entities.Collectible{}

	stmt :=
		`
		UPDATE
			%s
		SET
			deleted_at = $1 
		WHERE
			id = $2
		`

	stmt = fmt.Sprintf(
		stmt,
		collectible.TableName(),
	)

	cmd, err := db.Exec(ctx, stmt, time.Now(), id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *CollectibleRepository) Block(ctx context.Context, db database.QueryExecer,
	id int64) error {
	collectible := &entities.Collectible{}
	stmt :=
		`
		UPDATE
			%s
		SET
			status = -1
		WHERE
			id = $1
		`
	stmt = fmt.Sprintf(
		stmt,
		collectible.TableName(),
	)

	cmd, err := db.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("block affected no row")
	}
	return nil
}

func (r *CollectibleRepository) GetOwnerAddress(ctx context.Context, db database.QueryExecer, tokenId decimal.Decimal) *string {
	collectible := &entities.Collectible{}

	var stmt string

	stmt =
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			%s = $1 
			limit 1
		`

	stmt = fmt.Sprintf(
		stmt,
		entities.Collectible_TokenOwner,
		collectible.TableName(),
		entities.Collectible_TokenId,
	)

	row := db.QueryRow(ctx, stmt, tokenId)

	var owner string

	switch err := row.Scan(&owner); err {
	case nil:
		return &owner
	case pgx.ErrNoRows:
		return nil
	default:
		return nil
	}

}

func (r *CollectibleRepository) GetCreatorAddress(ctx context.Context, db database.QueryExecer, tokenId decimal.Decimal) *string {
	collectible := &entities.Collectible{}

	var stmt string

	stmt =
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			%s = $1 
			limit 1
		`

	stmt = fmt.Sprintf(
		stmt,
		entities.Collectible_Creator,
		collectible.TableName(),
		entities.Collectible_TokenId,
	)

	row := db.QueryRow(ctx, stmt, tokenId)

	var creator string

	switch err := row.Scan(&creator); err {
	case nil:
		return &creator
	case pgx.ErrNoRows:
		return nil
	default:
		return nil
	}

}

func (r *CollectibleRepository) GetIdByTokenId(ctx context.Context, db database.QueryExecer, tokenId decimal.Decimal) int64 {
	collectible := &entities.Collectible{}

	var stmt string

	stmt =
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			%s = $1 
			limit 1
		`

	stmt = fmt.Sprintf(
		stmt,
		entities.Collectible_Id,
		collectible.TableName(),
		entities.Collectible_TokenId,
	)

	row := db.QueryRow(ctx, stmt, tokenId)

	var id int64

	switch err := row.Scan(&id); err {
	case nil:
		return id
	case pgx.ErrNoRows:
		return 0
	default:
		return 0
	}

}
