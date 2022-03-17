package repositories

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"nft-backend/internal/entities"
	"strings"
	"testing"
)

func TestCollectibleCategoryRepo_Insert(t *testing.T) {
	test := require.New(t)

	ctx := context.Background()
	pg, err := initPostgres()
	test.NoError(err)

	collectible := &entities.Collectible{
		Title:               "title-TestCollectibleCategoryRepo_Insert",
		Description:         "description-TestCollectibleCategoryRepo_Insert",
		UploadFile:          "upload-file-TestCollectibleCategoryRepo_Insert",
		Royalties:           decimal.NewFromFloat(0.3),
		UnlockOncePurchased: true,
		Properties:          nil,
		QuoteToken: &entities.Token{
			ID: 1,
		},
	}
	err = (&CollectibleRepository{}).Insert(ctx, pg.Pool, collectible)
	test.NoError(err)
	//fmt.Println(collectible)

	category := &entities.Category{ID: 1}

	collectibleCategory, err := (&CollectibleCategoryRepo{}).Insert(ctx, pg.Pool, collectible, category)
	test.NoError(err)
	//fmt.Println(collectibleCategory)

	// Raw query to check
	queriedCollectibleCategory := &entities.CollectibleCategory{}
	stmt :=
		`
		SELECT
			%s
		FROM
			%s
		WHERE
			collectible_id = $1
			AND category_id = $2
		`
	names, values := queriedCollectibleCategory.FieldsAndValues()
	stmt = fmt.Sprintf(
		stmt,
		strings.Join(names, ", "),
		queriedCollectibleCategory.TableName(),
	)
	row := pg.Pool.QueryRow(ctx, stmt, collectibleCategory.CollectibleID, collectibleCategory.CategoryId)
	err = row.Scan(values...)
	test.NoError(err)
	//fmt.Println(queriedCollectibleCategory)

	test.Equal(queriedCollectibleCategory.CollectibleID, collectibleCategory.CollectibleID)
	test.Equal(queriedCollectibleCategory.CategoryId, collectibleCategory.CategoryId)
	test.Equal(queriedCollectibleCategory.CreatedAt, collectibleCategory.CreatedAt)
	test.Equal(queriedCollectibleCategory.UpdatedAt, collectibleCategory.UpdatedAt)
	test.Equal(queriedCollectibleCategory.DeletedAt, collectibleCategory.DeletedAt)
}
