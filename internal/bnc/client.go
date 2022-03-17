package bnc

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	NetworkEndpoint string

	NFTContractAddress      string
	ExchangeContractAddress string

	NFTCreateEventSignature    []byte
	ExchangeSellEventSignature []byte
	ExchangeBuyEventSignature  []byte
}

type Client struct {
	logs chan types.Log
	conn chan error

	config    *Config
	pgPool    *pgxpool.Pool
	zapLogger *zap.Logger

	stop chan struct{}
}

func NewClient(pgPool *pgxpool.Pool, zapLogger *zap.Logger, networkEndpoint string) (*Client, error) {
	//mnemonic := "rare fox animal view solid elder eye cushion tissue execute remember canoe"

	c := &Client{
		logs: make(chan types.Log),
		conn: make(chan error),
		config: &Config{
			NetworkEndpoint:           networkEndpoint,
			NFTContractAddress:        os.Getenv("NFT_CONTRACT_ADDRESS"),
			NFTCreateEventSignature:   []byte(`Transfer(address,address,uint256)`),
			ExchangeContractAddress:   os.Getenv("NFT_EXCHANGE_CONTRACT_ADDRESS"),
			ExchangeBuyEventSignature: []byte(`BuyToken(uint256,(uint256,enum)`),
		},
		pgPool:    pgPool,
		zapLogger: zapLogger,
		stop:      make(chan struct{}),
	}

	go c.run()
	go c.connect()

	return c, nil
}
