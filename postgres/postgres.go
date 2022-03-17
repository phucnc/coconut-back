package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Config struct {
	Host    string
	Port    string
	Usr     string
	Pwd     string
	Db      string
	SSLMode string
}

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, connURI string) (*Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(connURI)
	if err != nil {
		return nil, errors.Wrap(err, "ParseConfig")
	}
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &shopspring.Numeric{},
			Name:  "numeric",
			OID:   pgtype.NumericOID,
		})
		return nil
	}

	pool, err := pgxpool.Connect(ctx, poolConfig.ConnString())
	if err != nil {
		return nil, err
	}

	postgres := &Postgres{
		Pool: pool,
	}

	return postgres, nil
}

func (c *Config) ToConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", c.Usr, c.Pwd, c.Host, c.Port, c.Db, c.SSLMode)
}

func (pg *Postgres) Shutdown() {
	pg.Pool.Close()
}
