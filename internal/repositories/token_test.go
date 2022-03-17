package repositories

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTokenRepository_GetByName(t *testing.T) {
	test := require.New(t)

	ctx := context.Background()
	pg, err := initPostgres()
	test.NoError(err)

	token, err := (&TokenRepository{}).GetByName(ctx, pg.Pool, "BMP")
	//fmt.Println(token)
	test.NoError(err)
	test.Equal(token.Name, "BMP")
}
