package repositories

import (
	"context"
	"github.com/pkg/errors"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"time"
	"fmt"
	//"strconv"
)

type CollectibleLikeRepo struct{}


func (r *CollectibleLikeRepo) Insert(ctx context.Context, db database.QueryExecer, collectible_id int64, account_id int64) ( error) {
	t := time.Now().Truncate(time.Microsecond)
	collectibleLike := &entities.CollectibleLike{
		CollectibleID: collectible_id,
		AccountID:    account_id,
		CreatedAt:     t,
		UpdatedAt:     t,
	}

	fieldNames, fieldValues := collectibleLike.FieldsAndValues()
	cmd, err := database.Insert(ctx, db, collectibleLike, fieldNames, fieldValues)
	if err != nil {
		return  err
	}
	if cmd.RowsAffected() < 1 {
		return  errors.New("no rows affected")
	}
	return  nil
}

func (r *CollectibleLikeRepo) Delete(ctx context.Context, db database.QueryExecer, collectible_id int64, account_id int64) ( error) {
    
    stmt := `
	DELETE FROM %s
	WHERE collectible_id = $1 AND account_id = $2`

	stmt = fmt.Sprintf(
		stmt,
		"collectible_like",
	)
	_, err := db.Exec(ctx, stmt, collectible_id, account_id)
	if err != nil {
  	panic(err)
  	return err
	}

	return  nil
}

func (r *CollectibleLikeRepo) Check(ctx context.Context, db database.QueryExecer, collectible_id int64, account_id int64) (bool, error) {
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
			collectible_id = $1
			AND account_id = $2
		`
    
    stmt = fmt.Sprintf(
		stmt,
		"collectible_like",
	)

	row := db.QueryRow(ctx, stmt, collectible_id, account_id)
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


func (r *CollectibleLikeRepo) Count(ctx context.Context, db database.QueryExecer, collectible_id int64) (int) {
	var total int


	stmt :=
		`
		SELECT
			count(*) 
		FROM
			%s
		WHERE
			collectible_id = $1
		`
    
    stmt = fmt.Sprintf(
		stmt,
		"collectible_like",
	)

	row := db.QueryRow(ctx, stmt, collectible_id)
	err := row.Scan(&total)


	switch err {
	case nil:
		return total
	default:
		return 0
	}
	}




/*func (r *CollectibleLikeRepo) Delete(ctx context.Context, db database.QueryExecer, collectible *entities.Collectible, account *entities.Account) (*entities.CollectibleLike, error) {
	t := time.Now().Truncate(time.Microsecond)
	collectibleLike := &entities.CollectibleLike{
		CollectibleID: collectible.Id,
		AccountID:    account.ID,
		CreatedAt:     t,
		UpdatedAt:     t,
	}

	fieldNames, fieldValues := collectibleLike.FieldsAndValues()
	cmd, err := database.Insert(ctx, db, collectibleLike, fieldNames, fieldValues)
	if err != nil {
		return nil, err
	}
	if cmd.RowsAffected() < 1 {
		return nil, errors.New("no rows affected")
	}
	return collectibleLike, nil
}*/

