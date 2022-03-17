package main

import (
	"context"
	"log"
	"nft-backend/internal/bnc"
	"nft-backend/internal/database"
	"os"
	"os/signal"
	"time"

	"nft-backend/internal/services"
	"nft-backend/internal/storage"
	"nft-backend/postgres"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func main() {
	// Interrupt listener
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger: %v \n", err)
		return
	}
	zapLogger.Sugar().Info("initialized logger successfully")

	ctx := context.Background()

	initCtx := context.Background()
	pgConfig := database.NewPostgreConfigFromEnv()

	pg, err := postgres.NewPostgres(initCtx, pgConfig.ToConnStr())
	if err != nil {
		zapLogger.Sugar().Error(err)
		return
	}

	m, err := migrate.New(
		"file://migrations",
		pgConfig.ToConnStr(),
	)
	if err != nil {
		zapLogger.Sugar().Errorf("failed to init magration: %v \n", err)
	}

	switch err := m.Up(); err {
	case nil:
		zapLogger.Sugar().Info("migrate schema successfully")
	case migrate.ErrNoChange:
		zapLogger.Info("schema migration is up-to-date")
	default:
		zapLogger.Sugar().Errorf("failed to init migrate up: %v \n", err)
	}

	uploader, err := storage.NewUploader(ctx, storage.LoadConfigFromEnv())
	if err != nil {
		zapLogger.Sugar().Errorf("failed to init uploader up: %v \n", err)
		return
	}
	zapLogger.Sugar().Info("initialized uploader")

	s := services.NewServer(ctx, os.Getenv("BACKEND_PORT"), zapLogger, pg, uploader)
	zapLogger.Sugar().Info("http server running on ", os.Getenv("BACKEND_PORT"))

	//wss://data-seed-prebsc-1-s1.binance.org:8545
	//bsc
	bncClient, err := bnc.NewClient(pg.Pool, zapLogger, os.Getenv("WSS_BSC"))
	if err != nil {
		zapLogger.Sugar().Error(err)
		return
	}

	defer func() {
		// Release resources
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			zapLogger.Sugar().Error(err)
		}

		bncClient.Shutdown()

		zapLogger.Sugar().Info("graceful shutdown, shut down")
	}()

	// Graceful shutdown
	select {
	case <-interrupt:
		zapLogger.Sugar().Info("graceful shutdown, shutting down")
		return
	}
}
