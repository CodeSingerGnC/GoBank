package db

import "context"

// CreateUserTxParam 定义转账操作需要的参数
type CreateUserTxParam struct {
	CreateUserParams
	AfterCreate func() error 
}

// CreateUserTx 用于完成转账操作，将产生一条转账记录，两条账户条目，以及更新两个账户的余额。
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParam) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		err = arg.AfterCreate()
		return err
	})
	if err != nil {
		return err
	}

	return nil
}