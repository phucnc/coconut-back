package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"strings"
)

type CategoryRepo struct {}

func (r *CategoryRepo) GetCategoryByNames(ctx context.Context, db database.QueryExecer, names []string) ([]*entities.Category, error) {
	for _, name := range names {
		name = strings.ToLower(name)
	}

	stmt :=
		`
		SELECT 
			%s
		FROM
			%s
		WHERE
			lower(name) = ANY ($1);
		`

	tableName := (&entities.Category{}).TableName()
	//columnNames, _ := database.GetFieldNameAndValues(&entities.Category{})
	columnNames, _ := (&entities.Category{}).FieldsAndValues()

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ","),
		tableName,
	)
	//fmt.Println(stmt)

	row, err := db.Query(ctx, stmt, names)
	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, errors.Wrap(err, "Query")
	}

	categories := make([]*entities.Category, 0)
	for row.Next() {
		category := &entities.Category{}
		_, columnValues := category.FieldsAndValues()
		//_, columnValues := database.GetFieldNameAndValues(category)
		err := row.Scan(columnValues...)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoryRepo) GetAll(ctx context.Context, db database.QueryExecer) (entities.Categories, error) {
	stmt :=
		`
		SELECT
			%s
		FROM
			%s
		`

	tableName := (&entities.Category{}).TableName()
	columnNames, _ := (&entities.Category{}).FieldsAndValues()

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ","),
		tableName,
	)

	row, err := db.Query(ctx, stmt)
	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, errors.Wrap(err, "Query")
	}

	categories := make(entities.Categories, 0)
	for row.Next() {
		category := &entities.Category{}
		_, columnValues := category.FieldsAndValues()
		err := row.Scan(columnValues...)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoryRepo) InsertCategory(ctx context.Context, db database.QueryExecer,name string) ( error) {

		stmt :=
		`
		INSERT INTO Category (name)
			VALUES (lower($1))
		`

	cmd, err := db.Exec(ctx, stmt, name)
		if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return  nil
}

func (r *CategoryRepo) UpdateCategory(ctx context.Context, db database.QueryExecer, id int64, name string) ( error) {

		stmt :=
		`  
		UPDATE category set name = lower($1) where id = $2
		`
	cmd, err := db.Exec(ctx, stmt, name, id)
		if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return  nil
}