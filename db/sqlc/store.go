package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store 组合 *Queries 和 *sql.DB，便于执行事务
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParam) error
}

// SQLStore 组合 *Queries 和 *sql.DB，便于执行事务
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore 创建新的 Store 对象用于执行事务
func NewStore (db *sql.DB) Store {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

// execTx 执行 fn 定义的命令集合
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err %v, rb err %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
