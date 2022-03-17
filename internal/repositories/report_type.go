package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"strings"
	//"time"
)

type ReportTypeRepo struct {}



func (r *ReportTypeRepo) Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.ReportType, error) {
	reporttype := &entities.ReportType{}
	stmt :=
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			id = $1
		`

	columnNames, columnValues := reporttype.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = reporttype.TableName() + "." + columnName
	}


	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		reporttype.TableName(),
	)


	row := db.QueryRow(ctx, stmt, id)
	err := row.Scan(columnValues...)

	switch err {
	case nil:
		return reporttype, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}



func (r *ReportTypeRepo) GetAllReportTypes(ctx context.Context, db database.QueryExecer) (entities.ReportTypes, error) {
	stmt :=
		`
		SELECT
			%s
		FROM
			%s
		`
	tableName := (&entities.ReportType{}).TableName()
	columnNames, _ := (&entities.ReportType{}).FieldsAndValues()

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

	items := make(entities.ReportTypes, 0)
	for row.Next() {
		item := &entities.ReportType{}
		_, columnValues := item.FieldsAndValues()
		err := row.Scan(columnValues...)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}


func (r *ReportTypeRepo) Insert(ctx context.Context, db database.QueryExecer,description string) ( error) {

		stmt :=
		`
		INSERT INTO report_type (description)
			VALUES ($1)
		`

	cmd, err := db.Exec(ctx, stmt, description)
		if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return  nil
}

func (r *ReportTypeRepo) Update(ctx context.Context, db database.QueryExecer, id int64, description string) ( error) {

		stmt :=
		`  
		UPDATE report_type set description = $1 where id = $2
		`
	cmd, err := db.Exec(ctx, stmt, description, id)
		if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return  nil
}
