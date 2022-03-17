package repositories

import (
	"context"
	//"github.com/pkg/errors"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	//"time"
	"fmt"
	//"strconv"
	"github.com/jackc/pgx/v4"
	//"github.com/pkg/errors"
	"strings"

)

type CommentRepo struct {
}


func (r *CommentRepo) Insert(ctx context.Context, db database.QueryExecer, 
	comment *entities.Comment) ( error) {

		names, values := comment.FieldsAndValues(
		entities.Comment_CollectibleID,
		entities.Comment_AccountID,
		entities.Comment_Content,
	)
	err := database.InsertReturning(
		ctx,
		db,
		comment,
		names,
		values,
		[]string{
			entities.Comment_Id,
			entities.Comment_CreatedAt,
			entities.Comment_UpdatedAt,
			entities.Comment_DeletedAt,
		},
		[]interface{}{
			&comment.Id,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.DeletedAt,
		},
	)
	if err != nil {
		return err
	}
	return nil
}



func (r *CommentRepo) Delete(ctx context.Context, db database.QueryExecer, id int64) ( error)  {
	stmt := `
	DELETE FROM %s
	WHERE id = $1`

	stmt = fmt.Sprintf(
		stmt,
		"comment",
	)
	_, err := db.Exec(ctx, stmt, id)
	if err != nil {
  	panic(err)
  	return err
	}

	return  nil
}

func (r *CommentRepo) Check(ctx context.Context, db database.QueryExecer, id int64, account_id int64) (bool, error) {
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
		"comment",
	)

	row := db.QueryRow(ctx, stmt, id, account_id)
	err := row.Scan(&total)

	if total >0 {
		check = true;
	}

	switch err {
	case nil:
		return check, nil
	default:
		return false, err
	}
}




func (r *CommentRepo) GetByCollectible(ctx context.Context, db database.QueryExecer, collectible_id int64, limit int, offset int) ([]*entities.Comment, error) {

	//var args []interface{}

			stmt :=
		`
		SELECT
			%s
		FROM
			%s
			where collectible_id = $1
			order by id desc
			LIMIT 
			$2
			OFFSET 
		    $3
		`

		comment := &entities.Comment{}
		columnNames, columnValues := comment.FieldsAndValues()
		for i, columnName := range columnNames {
			columnNames[i] = comment.TableName() + "." + columnName
		}
		columnValues = append(columnValues, &comment.Id)

		stmt = fmt.Sprintf(
			stmt,
			strings.Join(columnNames, ", "),
			comment.TableName(),
		)
		
		//fmt.Println(stmt)
	
	rows, err := db.Query(ctx, stmt, collectible_id, limit, offset)
	switch err {
		case nil:

			comments := make([]*entities.Comment, 0, limit)
		for rows.Next() {
			comment := &entities.Comment{}
			columnNames, columnValues := comment.FieldsAndValues()
			for i, columnName := range columnNames {
				columnNames[i] = comment.TableName() + "." + columnName
			}

			err := rows.Scan(columnValues...)
			if err != nil {
				return nil, err
			}
			comments = append(comments, comment)
		}

		fmt.Printf("%v", comments)

		return comments, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}
