package repositories

import (
	"context"
	"github.com/pkg/errors"
	"nft-backend/internal/database"
	"nft-backend/internal/entities"
	"time"
)

type CollectibleCategoryRepo struct{}

//
/*func (r *CollectibleCategoryRepo) Get(ctx context.Context, db database.QueryExecer, categoryNames []string) ([]*entities.CollectibleCategory, error) {
	stmt :=
		`
		SELECT 
			collectible_id, category_id
		FROM 
			collectible_category
		WHERE 
			name IN ($1)
			delete_at IS NULL
		`

	rows, err := db.Query(ctx, stmt, categoryNames)
	if err != nil {
		return nil, err
	}

	categories := make([]*entities.CollectibleCategory, 0, len(categoryNames))
	for rows.Next() {
		category := &entities.CollectibleCategory{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}*/

func (r *CollectibleCategoryRepo) Insert(ctx context.Context, db database.QueryExecer, collectible *entities.Collectible, category *entities.Category) (*entities.CollectibleCategory, error) {
	t := time.Now().Truncate(time.Microsecond)
	collectibleCategory := &entities.CollectibleCategory{
		CollectibleID: collectible.Id,
		CategoryId:    category.ID,
		CreatedAt:     t,
		UpdatedAt:     t,
		DeletedAt:     nil,
	}

	//fieldNames, fieldValues := database.GetFieldNameAndValues(collectibleCategory)
	fieldNames, fieldValues := collectibleCategory.FieldsAndValues()
	cmd, err := database.Insert(ctx, db, collectibleCategory, fieldNames, fieldValues)
	if err != nil {
		return nil, err
	}
	if cmd.RowsAffected() < 1 {
		return nil, errors.New("no rows affected")
	}
	return collectibleCategory, nil
}
