package repositories

import (
	"context"
	"fmt"
	//"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	//"github.com/pkg/errors"
	//"github.com/shopspring/decimal"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"strings"
	//"database/sql"
	//"time"
)

type KYCRepository struct {
}

type KYC struct {
}

func (r *KYCRepository) Insert(ctx context.Context, db database.QueryExecer, kyc *entities.KYC) error {
	names, values := kyc.FieldsAndValues(
		entities.Account_Address,
	)
	err := database.InsertReturning(
		ctx,
		db,
		kyc,
		names,
		values,
		[]string{
			entities.KYC_Id,
			entities.KYC_CreatedAt,
			entities.KYC_UpdatedAt,
			entities.KYC_DeletedAt,
		},
		[]interface{}{
			&kyc.Id,
			&kyc.CreatedAt,
			&kyc.UpdatedAt,
			&kyc.DeletedAt,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *KYCRepository) Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.KYC, error) {
	kyc := &entities.KYC{}
	stmt :=
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			account.id = $1
			AND account.deleted_at IS NULL
		`

	columnNames, columnValues := kyc.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = kyc.TableName() + "." + columnName
	}

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		kyc.TableName(),
	)

	row := db.QueryRow(ctx, stmt, id)
	err := row.Scan(columnValues...)

	switch err {
	case nil:
		return kyc, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *KYCRepository) GetByAddress(ctx context.Context, db database.QueryExecer, address string) (*entities.Account, error) {
	account := &entities.Account{}
	stmt :=
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			account.address ILIKE $1
			AND account.deleted_at IS NULL
		`

	columnNames, columnValues := account.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = account.TableName() + "." + columnName
	}

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		account.TableName(),
	)

	row := db.QueryRow(ctx, stmt, address)
	err := row.Scan(columnValues...)

	switch err {
	case nil:
		return account, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *KYCRepository) CheckEmail(ctx context.Context, db database.QueryExecer, email string) (bool, error) {
	kyc := &entities.KYC{}
	stmt :=
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			kyc.email = $1
		`
	columnNames, columnValues := kyc.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = kyc.TableName() + "." + columnName
	}

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		kyc.TableName(),
	)

	row := db.QueryRow(ctx, stmt, email)
	err := row.Scan(columnValues...)

	switch err {
	case nil:
		return false, nil
	case pgx.ErrNoRows:
		return true, nil
	default:
		return false, err
	}
}

type KYCPaging struct {
	Limit  int
	Offset int
	Status int
}

func (r *KYCRepository) SearchPaging(ctx context.Context, db database.QueryExecer, paging *AccountPaging) ([]*entities.Account, error) {

	likeKeys := ""
	if paging.Keys != "" {
		likeKeys = fmt.Sprintf(`'%%%s%%'`, paging.Keys)
	}

	var stmt string

	var args []interface{}

	stmt =
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			account.deleted_at IS NULL
			AND (address like %s or lower(username) like lower(%s))
		ORDER BY account.id DESC
		LIMIT 
			$1
		OFFSET 
		    $2
		`

	account := &entities.Account{}
	columnNames, columnValues := account.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = account.TableName() + "." + columnName
	}
	columnValues = append(columnValues, &account.Id)

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		account.TableName(),
		likeKeys,
		likeKeys,
	)
	args = []interface{}{paging.Limit, paging.Offset}
	//fmt.Println(stmt)

	//fmt.Println(stmt)
	rows, err := db.Query(ctx, stmt, args...)
	switch err {
	case nil:

		accounts := make([]*entities.Account, 0, paging.Limit)
		for rows.Next() {
			account := &entities.Account{}
			columnNames, columnValues := account.FieldsAndValues()
			for i, columnName := range columnNames {
				columnNames[i] = account.TableName() + "." + columnName
			}

			err := rows.Scan(columnValues...)
			if err != nil {
				return nil, err
			}
			accounts = append(accounts, account)
		}

		return accounts, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *KYCRepository) Paging(ctx context.Context, db database.QueryExecer, paging *AccountPaging) ([]*entities.Account, error) {

	var stmt string
	where := ""
	if paging.Status == -1 {
		where = "AND account.status = -1"
	} else {
		where = "AND account.status = 0"
	}

	var args []interface{}

	stmt =
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			account.deleted_at IS NULL
			%s
		ORDER BY account.id DESC
		LIMIT 
			$1
		OFFSET 
		    $2
		
		`

	account := &entities.Account{}
	columnNames, columnValues := account.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = account.TableName() + "." + columnName
	}
	columnValues = append(columnValues, &account.Id)

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		account.TableName(),
		where,
	)
	args = []interface{}{paging.Limit, paging.Offset}
	//fmt.Println(stmt)

	//fmt.Println(stmt)
	rows, err := db.Query(ctx, stmt, args...)
	switch err {
	case nil:

		accounts := make([]*entities.Account, 0, paging.Limit)
		for rows.Next() {
			account := &entities.Account{}
			columnNames, columnValues := account.FieldsAndValues()
			for i, columnName := range columnNames {
				columnNames[i] = account.TableName() + "." + columnName
			}

			err := rows.Scan(columnValues...)
			if err != nil {
				return nil, err
			}
			accounts = append(accounts, account)
		}

		return accounts, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}
