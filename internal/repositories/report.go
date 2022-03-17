package repositories

import (
	"context"
	//"github.com/pkg/errors"
	"github.com/jackc/pgx/v4"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"github.com/pkg/errors"
	//"time"
	"fmt"
	"strings"
	"strconv"

)

type ReportRepo struct {
}

type ReportPaging struct {
	Account int64
	Collectible     int64
	Status     int64
	Offset       int64
	Limit      int64
	From int
	To int

}

func (r *ReportRepo) Insert(ctx context.Context, db database.QueryExecer, 
	report *entities.CollectibleReport) ( error) {

		names, values := report.FieldsAndValues(
		entities.CollectibleReport_CollectibleID,
		entities.CollectibleReport_AccountID,
		entities.CollectibleReport_ReportTypeID,
		entities.CollectibleReport_Content,
	)
	err := database.InsertReturning(
		ctx,
		db,
		report,
		names,
		values,
		[]string{
			entities.CollectibleReport_Id,
			entities.CollectibleReport_CreatedAt,
			entities.CollectibleReport_UpdatedAt,
			entities.CollectibleReport_DeletedAt,
		},
		[]interface{}{
			&report.Id,
			&report.CreatedAt,
			&report.UpdatedAt,
			&report.DeletedAt,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReportRepo) Update(ctx context.Context, db database.QueryExecer, id int64, status int64) ( error) {

		stmt :=
		`  
		UPDATE collectible_report set status = $1 where id = $2
		`
	cmd, err := db.Exec(ctx, stmt, status, id)
		if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return  nil
}

func (r *ReportRepo) Paging(ctx context.Context, db database.QueryExecer, paging *ReportPaging) ([]*entities.Report, error) {


	var stmt string

	var where = " 1> 0";
	if(paging.Collectible >0) {
		where += " AND collectible_id = "+ strconv.FormatInt(paging.Collectible,10)
	}

	if(paging.Account >0) {
		where += " AND account_id = "+ strconv.FormatInt(paging.Account,10)
	}

	if(paging.Status>=0) {
		where += " AND status = "+ strconv.FormatInt(paging.Status,10)
	}

	if( paging.From != 0) {
	
		where += " AND created_at >= now() - INTERVAL '"+ strconv.Itoa(paging.From)+" hours'";
	}

	if( paging.To != 0) {
		where += " AND created_at <= now() - INTERVAL '"+ strconv.Itoa(paging.To) +" hours' ";
	}



	var args []interface{}

		stmt =
			`
		SELECT
			%s
		FROM
			%s
		WHERE %s

		ORDER BY id DESC
		LIMIT 
			$1
		OFFSET 
		    $2
		
		`

		report := &entities.Report{}
		columnNames, columnValues := report.FieldsAndValues()
		for i, columnName := range columnNames {
			columnNames[i] = report.TableName() + "." + columnName
		}
		columnValues = append(columnValues, &report.Id)

		stmt = fmt.Sprintf(
			stmt,
			strings.Join(columnNames, ", "),
			report.TableName(),
			where,
		)
		args = []interface{}{ paging.Limit, paging.Offset}
		//fmt.Println(stmt)
	
	//fmt.Println(stmt)
	rows, err := db.Query(ctx, stmt, args...)
	switch err {
		case nil:

			reports := make([]*entities.Report, 0, paging.Limit)
		for rows.Next() {
			report := &entities.Report{}
			columnNames, columnValues := report.FieldsAndValues()
			for i, columnName := range columnNames {
				columnNames[i] = report.TableName() + "." + columnName
			}

			err := rows.Scan(columnValues...)
			if err != nil {
				return nil, err
			}
			reports = append(reports, report)
		}

		return reports, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

