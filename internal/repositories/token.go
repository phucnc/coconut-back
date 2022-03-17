package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"strings"
)

type TokenRepository struct {}

func (r *TokenRepository) GetByName(ctx context.Context, db database.QueryExecer, name string) (*entities.Token, error) {
	token := new(entities.Token)

	stmt :=
		`
		SELECT 
			%s
		FROM
			%s
		WHERE
			lower(name) = lower($1);
		`

	columnNames, columnValues := token.FieldsAndValues()

	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		token.TableName(),
	)

	row := db.QueryRow(ctx, stmt, name)
	err := row.Scan(columnValues...)
	switch err {
	case nil:
		return token, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}