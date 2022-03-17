package repositories

import (
	"context"
	"nft-backend/internal/database"
	"nft-backend/postgres"
)

func initPostgres() (*postgres.Postgres, error) {
	initCtx := context.Background()
	connURI := database.NewPostgreConfigFromEnv().ToConnStr()

	return postgres.NewPostgres(initCtx, connURI)
}