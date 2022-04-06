package repositories

import (
	"context"
	"fmt"
	//"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	//"github.com/shopspring/decimal"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"strings"
	//"database/sql"
	//"time"
)

type AccountRepository struct {
}

type Account struct {
}

func (r *AccountRepository) Insert(ctx context.Context, db database.QueryExecer, account *entities.Account) error {
	names, values := account.FieldsAndValues(
		entities.Account_Address,
	)
	err := database.InsertReturning(
		ctx,
		db,
		account,
		names,
		values,
		[]string{
			entities.Account_Id,
			entities.Account_CreatedAt,
			entities.Account_UpdatedAt,
			entities.Account_DeletedAt,
		},
		[]interface{}{
			&account.Id,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.DeletedAt,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.Account, error) {
	account := &entities.Account{}
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

	columnNames, columnValues := account.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = account.TableName() + "." + columnName
	}

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		account.TableName(),
	)

	row := db.QueryRow(ctx, stmt, id)
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

func (r *AccountRepository) GetByAddress(ctx context.Context, db database.QueryExecer, address string) (*entities.Account, error) {
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

func (r *AccountRepository) CheckUserName(ctx context.Context, db database.QueryExecer, address string, username string) (bool, error) {
	account := &entities.Account{}
	stmt :=
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			account.username = $1
			AND account.address != $2
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

	row := db.QueryRow(ctx, stmt, username, address)
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

func (r *AccountRepository) Update(ctx context.Context, db database.QueryExecer,
	address string, username string, info string) error {
	account := &entities.Account{}

	stmt :=
		`
		UPDATE
			%s
		SET
			%s = $1 , %s = $2
		WHERE
			lower(%s) = lower($3)
		`
	stmt = fmt.Sprintf(
		stmt,
		account.TableName(),
		entities.Account_Username,
		entities.Account_Info,
		entities.Account_Address,
	)

	fmt.Sprintf(stmt)
	if username != "" {
		cmd, err := db.Exec(ctx, stmt, username, info, address)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() < 1 {
			return errors.New("update affected no row")
		}
	} else {
		cmd, err := db.Exec(ctx, stmt, nil, info, address)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() < 1 {
			return errors.New("update affected no row")
		}
	}

	return nil
}

func (r *AccountRepository) UpdateAvatar(ctx context.Context, db database.QueryExecer,
	address string, avatar string) error {
	account := &entities.Account{}
	stmt :=
		`
		UPDATE
			%s
		SET
			%s = $1 
		WHERE
			lower(%s) = lower($2)
		`
	stmt = fmt.Sprintf(
		stmt,
		account.TableName(),
		entities.Account_Avatar,
		entities.Account_Address,
	)

	cmd, err := db.Exec(ctx, stmt, avatar, address)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *AccountRepository) UpdateCover(ctx context.Context, db database.QueryExecer,
	address string, cover string) error {
	account := &entities.Account{}
	stmt :=
		`
		UPDATE
			%s
		SET
			%s = $1
		WHERE
			lower(%s) = lower($2)
		`
	stmt = fmt.Sprintf(
		stmt,
		account.TableName(),
		entities.Account_Cover,
		entities.Account_Address,
	)

	cmd, err := db.Exec(ctx, stmt, cover, address)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *AccountRepository) Block(ctx context.Context, db database.QueryExecer,
	address string) error {
	account := &entities.Account{}
	stmt :=
		`
		UPDATE
			%s
		SET
			status = -1
		WHERE
			%s = $1
		`
	stmt = fmt.Sprintf(
		stmt,
		account.TableName(),
		entities.Account_Address,
	)

	cmd, err := db.Exec(ctx, stmt, address)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *AccountRepository) ResetKYC(ctx context.Context, db database.QueryExecer,
	address string, status int) error {
	account := &entities.Account{}
	stmt :=
		`
		UPDATE
			%s
		SET
			kyc_status = 0
		WHERE
			%s = $1
		`
	stmt = fmt.Sprintf(
		stmt,
		account.TableName(),
		entities.Account_Address,
	)

	cmd, err := db.Exec(ctx, stmt, address)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}
	return nil
}

func (r *AccountRepository) GetByUsername(ctx context.Context, db database.QueryExecer, username string) (*entities.Account, error) {
	account := &entities.Account{}
	stmt :=
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			account.username = $1
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

	row := db.QueryRow(ctx, stmt, username)
	err := row.Scan(columnValues...)

	switch err {
	case nil:
		return account, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
	return nil, nil
}

type AccountPaging struct {
	Keys   string
	Limit  int
	Offset int
	Status int
}

func (r *AccountRepository) SearchPaging(ctx context.Context, db database.QueryExecer, paging *AccountPaging) ([]*entities.Account, error) {

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

func (r *AccountRepository) Paging(ctx context.Context, db database.QueryExecer, paging *AccountPaging) ([]*entities.Account, error) {

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
