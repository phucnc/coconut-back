package repositories

import (
	"context"
	//"github.com/pkg/errors"
	"github.com/pkg/errors"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	//"time"
	"fmt"
	//"strconv"
	"github.com/jackc/pgx/v4"
	//"github.com/pkg/errors"
	"strings"
)

type NoticeRepo struct {
}

func (r *NoticeRepo) Insert(ctx context.Context, db database.QueryExecer,
	notice *entities.Notice) error {

	names, values := notice.FieldsAndValues(
		entities.Notice_AccountID,
		entities.Notice_Content,
		entities.Notice_FromAccountID,
		entities.Notice_CollectibleID,
	)
	err := database.InsertReturning(
		ctx,
		db,
		notice,
		names,
		values,
		[]string{
			entities.Notice_Id,
			entities.Notice_CreatedAt,
			entities.Notice_UpdatedAt,
			entities.Notice_DeletedAt,
		},
		[]interface{}{
			&notice.Id,
			&notice.CreatedAt,
			&notice.UpdatedAt,
			&notice.DeletedAt,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *NoticeRepo) Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.Notice, error) {
	notice := &entities.Notice{}
	stmt :=
		`
		SELECT
			%s 
		FROM	
			%s
		WHERE
			notice.id = $1
			AND notice.deleted_at IS NULL
		`

	columnNames, columnValues := notice.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = notice.TableName() + "." + columnName
	}

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		notice.TableName(),
	)

	row := db.QueryRow(ctx, stmt, id)
	err := row.Scan(columnValues...)

	switch err {
	case nil:
		return notice, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *NoticeRepo) Delete(ctx context.Context, db database.QueryExecer, id int64) error {
	stmt := `
	DELETE FROM %s
	WHERE id = $1`

	stmt = fmt.Sprintf(
		stmt,
		"notice",
	)
	_, err := db.Exec(ctx, stmt, id)
	if err != nil {
		panic(err)
		return err
	}

	return nil
}

func (r *NoticeRepo) Check(ctx context.Context, db database.QueryExecer, id int64, account_id int64) (bool, error) {
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
			id = $1 AND account_id = $2
		`

	stmt = fmt.Sprintf(
		stmt,
		"notice",
	)

	row := db.QueryRow(ctx, stmt, id, account_id)
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

type NoticePaging struct {
	Status  int
	Account int64
	Limit   int
	Offset  int
}

func (r *NoticeRepo) Paging(ctx context.Context, db database.QueryExecer, paging *NoticePaging) ([]*entities.Notice, error) {

	var args []interface{}

	stmt :=
		`
		SELECT
			%s
		FROM
			%s
			where account_id = $1
			order by id desc
			LIMIT 
			$2
			OFFSET 
		    $3
		`

	notice := &entities.Notice{}
	columnNames, columnValues := notice.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = notice.TableName() + "." + columnName
	}
	columnValues = append(columnValues, &notice.Id)

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		notice.TableName(),
	)

	args = []interface{}{paging.Account, paging.Limit, paging.Offset}

	//fmt.Println(stmt)

	rows, err := db.Query(ctx, stmt, args...)
	switch err {
	case nil:

		notices := make([]*entities.Notice, 0, paging.Limit)
		for rows.Next() {
			notice := &entities.Notice{}
			columnNames, columnValues := notice.FieldsAndValues()
			for i, columnName := range columnNames {
				columnNames[i] = notice.TableName() + "." + columnName
			}

			err := rows.Scan(columnValues...)
			if err != nil {
				return nil, err
			}
			//fmt.Println("%v", notice)
			notices = append(notices, notice)
		}

		//fmt.Printf("%v", notices)

		return notices, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *NoticeRepo) Update(ctx context.Context, db database.QueryExecer,
	id int64, notice *entities.Notice) error {
	//t := time.Now().Truncate(time.Microsecond)

	stmt :=
		`
		UPDATE
			%s
		SET
			%s = $1 
		WHERE
			id = $2
		`
	stmt = fmt.Sprintf(
		stmt,
		notice.TableName(),
		entities.Notice_Content,
	)

	fmt.Sprintf(stmt)
	cmd, err := db.Exec(ctx, stmt, notice.Content, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return nil
}

func (r *NoticeRepo) Read(ctx context.Context, db database.QueryExecer,
	id int) error {
	//t := time.Now().Truncate(time.Microsecond)

	stmt :=
		`
		UPDATE
			%s
		SET
			%s = 1 
		WHERE
			id = $1
		`
	stmt = fmt.Sprintf(
		stmt,
		"notice",
		entities.Notice_Status,
	)

	fmt.Sprintf(stmt)
	cmd, err := db.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return nil
}
