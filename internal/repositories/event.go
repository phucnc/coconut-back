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
	"time"
	"strconv"

)

type EventRepo struct {
}

type Event struct {
}

func (r *EventRepo) Insert(ctx context.Context, db database.QueryExecer, event *entities.Event) error {
	names, values := event.FieldsAndValues(
		entities.Event_Title,
		entities.Event_Banner,
		entities.Event_Content,
	)
	err := database.InsertReturning(
		ctx,
		db,
		event,
		names,
		values,
		[]string{
			entities.Event_Id,
			entities.Event_CreatedAt,
			entities.Event_UpdatedAt,
		},
		[]interface{}{
			&event.Id,
			&event.CreatedAt,
			&event.UpdatedAt,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *EventRepo) Get(ctx context.Context, db database.QueryExecer, id int64) (*entities.Event, error) {
	event := &entities.Event{}
	stmt :=
		`
		SELECT
			%s 
		FROM
			%s
		WHERE
			event.id = $1
			AND event.deleted_at IS NULL
		`

	columnNames, columnValues := event.FieldsAndValues()
	for i, columnName := range columnNames {
		columnNames[i] = event.TableName() + "." + columnName
	}


	stmt = fmt.Sprintf(
		stmt,
		strings.Join(columnNames, ", "),
		event.TableName(),
	)


	row := db.QueryRow(ctx, stmt, id)
	err := row.Scan(columnValues...)

	switch err {
	case nil:
		return event, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}



func (r *EventRepo) Update(ctx context.Context, db database.QueryExecer,
 id int64 ,event *entities.Event) error {
//t := time.Now().Truncate(time.Microsecond)

	stmt :=
		`
		UPDATE
			%s
		SET
			%s = $1 , %s = $2, %s = $3
		WHERE
			id = $4
		`
	stmt = fmt.Sprintf(
		stmt,
		event.TableName(),
		entities.Event_Title,
		entities.Event_Banner,
		entities.Event_Content,
	)

	fmt.Sprintf(stmt);
	cmd, err := db.Exec(ctx, stmt, event.Title, event.Banner, event.Content, id)
		if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return nil
}


func (r *EventRepo) Delete(ctx context.Context, db database.QueryExecer,
 id int64) error {
t := time.Now().Truncate(time.Microsecond)
	stmt :=
		`
		UPDATE
			event
		SET
			%s = $1 
		WHERE
			id = $2
		`
	stmt = fmt.Sprintf(
		stmt,
		entities.Event_DeletedAt,
	)

	fmt.Sprintf(stmt);
	cmd, err := db.Exec(ctx, stmt, t, id)
		if err != nil {
		return err
	}
	if cmd.RowsAffected() < 1 {
		return errors.New("update affected no row")
	}

	return nil
}


type EventPaging struct {
	Status int
	Limit      int
	Offset   int
}



func (r *EventRepo) Paging(ctx context.Context, db database.QueryExecer, paging *EventPaging) ([]*entities.Event, error) {


	var stmt string

	var where = "event.deleted_at IS NULL";

	if(paging.Status>=0) {
		where += " AND status = "+ strconv.Itoa(paging.Status)
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
		ORDER BY event.id DESC
		LIMIT 
			$1
		OFFSET 
		    $2
		`

		event := &entities.Event{}
		columnNames, columnValues := event.FieldsAndValues()
		for i, columnName := range columnNames {
			columnNames[i] = event.TableName() + "." + columnName
		}
		columnValues = append(columnValues, &event.Id)

		stmt = fmt.Sprintf(
			stmt,
			strings.Join(columnNames, ", "),
			event.TableName(),
			where,
		)
		args = []interface{}{ paging.Limit, paging.Offset}
		//fmt.Println(stmt)
	
	//fmt.Println(stmt)
	rows, err := db.Query(ctx, stmt, args...)
	switch err {
		case nil:

			events := make([]*entities.Event, 0, paging.Limit)
		for rows.Next() {
			event := &entities.Event{}
			columnNames, columnValues := event.FieldsAndValues()
			for i, columnName := range columnNames {
				columnNames[i] = event.TableName() + "." + columnName
			}

			err := rows.Scan(columnValues...)
			if err != nil {
				return nil, err
			}
			events = append(events, event)
		}

		return events, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}
