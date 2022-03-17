package repositories

import (
	"context"
	"github.com/pkg/errors"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"fmt"
	"strings"
	"github.com/jackc/pgx/v4"
	//"strconv"
)

type TrendRepo struct{	}



func (r *TrendRepo) Insert(ctx context.Context, db database.QueryExecer, td *entities.Trend) ( error) {

		stmt :=
		`
		INSERT INTO trend (collectible_id,in_order, adsvertisement)
			VALUES ($1, $2,$3)
		`

	cmd, err := db.Exec(ctx, stmt, td.CollectibleID, td.In_Order, td.Advertisement)
		if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return  nil
}

func (r *TrendRepo) Delete(ctx context.Context, db database.QueryExecer, id int64) ( error) {
    
    stmt := `
	DELETE FROM %s
	WHERE id = $1`

	stmt = fmt.Sprintf(
		stmt,
		"trend",
	)
	_, err := db.Exec(ctx, stmt, id)
	if err != nil {
  	panic(err)
  	return err
	}

	return  nil
}




func (r *TrendRepo) Get(ctx context.Context, db database.QueryExecer) (entities.Trends, error) {
	stmt :=
		`
		SELECT
			%s
		FROM
			%s
		order by in_order desc
		`
	tableName := (&entities.Trend{}).TableName()
	columnNames, _ := (&entities.Trend{}).FieldsAndValues()

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

	items := make(entities.Trends, 0)
	for row.Next() {
		item := &entities.Trend{}
		_, columnValues := item.FieldsAndValues()
		err := row.Scan(columnValues...)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}