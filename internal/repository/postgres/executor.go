package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type poolExecutor struct{ pool *pgxpool.Pool }

type ctxTxKey struct{}

type DBExecutor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func withTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxTxKey{}, tx)
}

func getExecutor(ctx context.Context, pool *pgxpool.Pool) DBExecutor {
	if tx, ok := ctx.Value(ctxTxKey{}).(pgx.Tx); ok && tx != nil {
		return txExecutor{tx: tx}
	}
	return poolExecutor{pool: pool}
}

func (p poolExecutor) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	tag, err := p.pool.Exec(ctx, sql, args...)
	return tag, err
}
func (p poolExecutor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}
func (p poolExecutor) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

type txExecutor struct{ tx pgx.Tx }

func (t txExecutor) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	tag, err := t.tx.Exec(ctx, sql, args...)
	return tag, err
}
func (t txExecutor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return t.tx.Query(ctx, sql, args...)
}
func (t txExecutor) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return t.tx.QueryRow(ctx, sql, args...)
}
