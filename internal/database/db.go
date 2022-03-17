package database

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type queryer interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

// execer is an interface for Exec
type execer interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type QueryExecer interface {
	queryer
	execer
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

// TxStarter is an interface to deal with transaction
type TxStarter interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// TxController is an interface to deal with transaction
type TxController interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Ext is a union interface which can bind, query, and exec
type Ext interface {
	QueryExecer
	TxStarter
}

type TxHandler = func(ctx context.Context, tx pgx.Tx) error

func ExecInTx(ctx context.Context, db Ext, txHandler TxHandler) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "db.Begin")
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		err = errors.Wrap(tx.Commit(ctx), "tx.Commit")
	}()
	err = txHandler(ctx, tx)
	return err
}

