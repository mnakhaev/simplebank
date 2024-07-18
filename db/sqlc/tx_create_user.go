package db

import (
	"context"
)

// CreateUserTxParams contains the input parameters of create user transaction
type CreateUserTxParams struct {
	CreateUserParams
	// AfterCreate will be executed after the user is created, in same transaction
	AfterCreate func(user User) error
}

// CreateUserTxResult is the result of create user transaction.
type CreateUserTxResult struct {
	User User
}

func (s *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		// execute callback function
		return arg.AfterCreate(result.User)
	})
	return result, err
}
