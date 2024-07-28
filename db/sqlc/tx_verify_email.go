package db

import (
	"context"
	"database/sql"
)

// VerifyEmailTxParams contains the input parameters of verify email transaction
type VerifyEmailTxParams struct {
	EmailID    int64
	SecretCode string
}

// VerifyEmailTxResult is the result of verify email transaction.
type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

// VerifyEmailTx updates secret code and ID after email verification and also updates user data in database.
func (s *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailID,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			Username: result.VerifyEmail.Username,
		})
		return err
	})
	return result, err
}
