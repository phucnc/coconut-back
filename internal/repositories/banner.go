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
	"strconv"
	"time"
)

type BannerRepository struct {
}

type Banner struct {
}

func (r *BannerRepository) Insert(ctx context.Context, db database.QueryExecer, banner *entities.Banner) error {
	names, values := banner.FieldsAndValues(
		entities.Banner_Name,
		entities.Banner_Picture,
		entities.Banner_Link,
		entities.Banner_Status,
		entities.Banner_Order,
	)
	err := database.InsertReturning(
		ctx,
		db,
		banner,
		names,
		values,
		[]string{
			entities.Banner_Id,
			entities.Banner_CreatedAt,
			entities.Banner_UpdatedAt,
		},
		[]interface{}{
			&banner.Id,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *BannerRepository) Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.Banner, error) {
	banner := &entities.Banner{}
	stmt :=
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			banner.id = $1
			AND banner.deleted_at IS NULL
		`

	columnNames, columnValues := banner.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = banner.TableName() + "." + columnName
	}

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		banner.TableName(),
	)

	row := db.QueryRow(ctx, stmt, id)
	err := row.Scan(columnValues...)

	switch err {
	case nil:
		return banner, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *BannerRepository) Update(ctx context.Context, db database.QueryExecer,
	id int64, banner *entities.Banner) error {
	//t := time.Now().Truncate(time.Microsecond)
	stmt :=
		`
		UPDATE
			%s
		SET
			%s = $1 , %s = $2 ,  %s = $3 , %s = $4 , %s = $5 
		WHERE
			id = $6
		`
	stmt = fmt.Sprintf(
		stmt,
		banner.TableName(),
		entities.Banner_Name,
		entities.Banner_Link,
		entities.Banner_Picture,
		entities.Banner_Status,
		entities.Banner_Order,
	)

	fmt.Sprintf(stmt)
	//fmt.Println(stmt)
	cmd, err := db.Exec(ctx, stmt, banner.Name, banner.Link,
		banner.Picture, banner.Status, banner.Order, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return nil
}

func (r *BannerRepository) Delete(ctx context.Context, db database.QueryExecer,
	id int64) error {
	t := time.Now().Truncate(time.Microsecond)
	stmt :=
		`
		UPDATE
			banner
		SET
			%s = $1 
		WHERE
			id = $2
		`
	stmt = fmt.Sprintf(
		stmt,
		entities.Banner_DeletedAt,
	)

	fmt.Sprintf(stmt)
	cmd, err := db.Exec(ctx, stmt, t, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return nil
}

type BannerPaging struct {
	Status int
	Limit  int
	Offset int
}

func (r *BannerRepository) Paging(ctx context.Context, db database.QueryExecer, paging *BannerPaging) ([]*entities.Banner, error) {

	var stmt string

	var where = "banner.deleted_at IS NULL"

	if paging.Status >= 0 {
		where += " AND status = " + strconv.Itoa(paging.Status)
	}

	var args []interface{}

	stmt =
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			%s
		ORDER BY banner.in_order ASC
		LIMIT 
			$1
		OFFSET 
		    $2
		`

	banner := &entities.Banner{}
	columnNames, columnValues := banner.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = banner.TableName() + "." + columnName
	}
	columnValues = append(columnValues, &banner.Id)

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		banner.TableName(),
		where,
	)
	args = []interface{}{paging.Limit, paging.Offset}
	//fmt.Println(stmt)

	//fmt.Println(stmt)
	rows, err := db.Query(ctx, stmt, args...)
	switch err {
	case nil:

		banners := make([]*entities.Banner, 0, paging.Limit)
		for rows.Next() {
			banner := &entities.Banner{}
			columnNames, columnValues := banner.FieldsAndValues()
			for i, columnName := range columnNames {
				columnNames[i] = banner.TableName() + "." + columnName
			}

			err := rows.Scan(columnValues...)
			if err != nil {
				return nil, err
			}
			banners = append(banners, banner)
		}

		return banners, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}
