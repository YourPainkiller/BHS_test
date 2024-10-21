package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txManagerKey struct{}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (m *TxManager) RunSerialazible(ctx context.Context, fn func(ctxTx context.Context) error) error {
	opt := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	return m.beginFunc(ctx, opt, fn)

}

func (m *TxManager) RunReadUncommited(ctx context.Context, fn func(ctxTx context.Context) error) error {
	opt := pgx.TxOptions{
		IsoLevel:   pgx.ReadUncommitted,
		AccessMode: pgx.ReadOnly,
	}
	return m.beginFunc(ctx, opt, fn)
}

func (m *TxManager) RunReadWriteCommited(ctx context.Context, fn func(ctxTx context.Context) error) error {
	opt := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}
	return m.beginFunc(ctx, opt, fn)
}

func (m *TxManager) beginFunc(ctx context.Context, opts pgx.TxOptions, fn func(ctxTx context.Context) error) error {
	tx, err := m.pool.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ctx = context.WithValue(ctx, txManagerKey{}, tx)
	if err := fn(ctx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (m *TxManager) GetQueryEngine(ctx context.Context) QueryEngine {
	v, ok := ctx.Value(txManagerKey{}).(QueryEngine)
	if ok && v != nil {
		return v
	} else {
		return m.pool
	}
}
