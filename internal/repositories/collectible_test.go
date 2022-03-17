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

func TestCollectibleRepository_Insert(t *testing.T) {
	test := require.New(t)

	ctx := context.Background()
	pg, err := initPostgres()
	test.NoError(err)

	// Try inserting with repo
	collectible := &entities.Collectible{
		Title:               "test-name",
		Description:         "test-description",
		UploadFile:          "test-upload-file",
		Royalties:           decimal.NewFromFloat(0.3),
		UnlockOncePurchased: true,
		InstantSalePrice:    decimal.NewFromInt(5),
		QuoteToken: &entities.Token{
			ID: 1,
		},
	}
	err = (&CollectibleRepository{}).Insert(ctx, pg.Pool, collectible)
	test.NoError(err)

	// Raw query to check
	queriedCollectible := &entities.Collectible{}
	stmt :=
		`
		SELECT
			%s
		FROM
			%s
		WHERE
			id = $1
		`
	names, values := queriedCollectible.FieldsAndValues()
	stmt = fmt.Sprintf(
		stmt,
		strings.Join(names, ", "),
		queriedCollectible.TableName(),
	)
	row := pg.Pool.QueryRow(ctx, stmt, collectible.Id)
	err = row.Scan(values...)
	test.NoError(err)

	test.Equal(queriedCollectible.Id, collectible.Id)
	test.Equal(queriedCollectible.Title, collectible.Title)
	test.Equal(queriedCollectible.Description, collectible.Description)
	test.Equal(queriedCollectible.UploadFile, collectible.UploadFile)
	test.Equal(queriedCollectible.Royalties.String(), collectible.Royalties.String())
	test.Equal(queriedCollectible.UnlockOncePurchased, collectible.UnlockOncePurchased)
	test.Equal(queriedCollectible.InstantSalePrice.String(), collectible.InstantSalePrice.String())
	test.Equal(queriedCollectible.CreatedAt, collectible.CreatedAt)
	test.Equal(queriedCollectible.UpdatedAt, collectible.UpdatedAt)
	test.Equal(queriedCollectible.DeletedAt, collectible.DeletedAt)
}

func TestCollectibleRepository_Get(t *testing.T) {
	test := require.New(t)

	ctx := context.Background()
	pg, err := initPostgres()
	test.NoError(err)

	collectible := &entities.Collectible{
		Title:               "test-name",
		Description:         "test-description",
		UploadFile:          "test-upload-file",
		Royalties:           decimal.NewFromFloat(0.3),
		UnlockOncePurchased: true,
		InstantSalePrice: decimal.NewFromInt(7),
		QuoteToken: &entities.Token{
			ID: 1,
		},
	}
	err = (&CollectibleRepository{}).Insert(ctx, pg.Pool, collectible)
	test.NoError(err)

	categories := map[int64]*entities.Category{
		1: {ID: 1},
	}

	for _, category := range categories {
		_, err = (&CollectibleCategoryRepo{}).Insert(ctx, pg.Pool, collectible, category)
		test.NoError(err)
	}

	queriedCollectible, err := (&CollectibleRepository{}).Get(ctx, pg.Pool, collectible.GUID)
	test.NoError(err)
	//fmt.Println(queriedCollectible)

	test.Equal(collectible.Id, queriedCollectible.Id)
	test.Equal(collectible.Title, queriedCollectible.Title)
	test.Equal(collectible.Description, queriedCollectible.Description)
	test.Equal(collectible.Royalties.String(), queriedCollectible.Royalties.String())
	test.Equal(collectible.UnlockOncePurchased, queriedCollectible.UnlockOncePurchased)
	test.Equal(collectible.InstantSalePrice.String(), queriedCollectible.InstantSalePrice.String())
	test.Equal(collectible.UploadFile, queriedCollectible.UploadFile)
	test.Equal(collectible.CreatedAt, queriedCollectible.CreatedAt)
	test.Equal(collectible.UpdatedAt, queriedCollectible.UpdatedAt)
	test.Equal(collectible.DeletedAt, queriedCollectible.DeletedAt)

	test.Len(queriedCollectible.Categories, len(categories))

	for _, category := range queriedCollectible.Categories {
		_, ok := categories[category.ID]
		test.True(ok)
	}
}
