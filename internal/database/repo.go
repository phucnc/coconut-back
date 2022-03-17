package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"strings"
)

// Insert entity using db Executor
func Insert(
	ctx context.Context,
	db QueryExecer,
	e Entity,
	insertFieldNames []string,
	insertFieldValues []interface{},
) (pgconn.CommandTag, error) {
	placeHolders := GeneratePlaceholder(len(insertFieldNames))
	stmt :=
		`
		INSERT INTO %s (
			%s
		) VALUES (
			%s
		);
		`
	stmt = fmt.Sprintf(
		stmt,
		e.TableName(),
		strings.Join(insertFieldNames, ","),
		strings.Join(placeHolders, ","),
	)
	return db.Exec(ctx, stmt, insertFieldValues...)
}

// InsertIgnoreConflict entity using db Executor
func InsertIgnoreConflict(
	ctx context.Context,
	db QueryExecer,
	e Entity,
	insertFieldNames []string,
	insertFieldValues []interface{},
) (pgconn.CommandTag, error) {
	placeHolders := GeneratePlaceholder(len(insertFieldNames))
	stmt :=
		`
		INSERT INTO %s (
			%s
		) VALUES (
			%s
		) ON CONFLICT DO NOTHING;
		`
	stmt = fmt.Sprintf(
		stmt,
		e.TableName(),
		strings.Join(insertFieldNames, ","),
		strings.Join(placeHolders, ","),
	)
	return db.Exec(ctx, stmt, insertFieldValues...)
}

// InsertReturning entity using db Executor
func InsertReturning(
	ctx context.Context,
	db QueryExecer,
	e Entity,
	insertFieldNames []string,
	insertFieldValues []interface{},
	returnFieldNames []string,
	returnFieldValue []interface{},
) error {
	placeHolders := GeneratePlaceholder(len(insertFieldNames))
	stmt :=
		`
		INSERT INTO %s (
			%s
		) VALUES (
			%s
		) RETURNING
			%s
		;
		`
	stmt = fmt.Sprintf(
		stmt,
		e.TableName(),
		strings.Join(insertFieldNames, ","),
		strings.Join(placeHolders, ","),
		strings.Join(returnFieldNames, ","),
	)
	return db.QueryRow(ctx, stmt, insertFieldValues...).Scan(returnFieldValue...)
}
