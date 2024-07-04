package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/mnakhaev/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user2.Username, user1.Username)
	require.Equal(t, user2.HashedPassword, user1.HashedPassword)
	require.Equal(t, user2.FullName, user1.FullName)
	require.Equal(t, user2.Email, user1.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Username, updatedUser.Username)
}

func TestUpdateUserEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Username, updatedUser.Username)
}

func TestUpdateUserPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(10)
	newHashPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Username, updatedUser.Username)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(10)
	newHashPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newHashPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)

	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashPassword, updatedUser.HashedPassword)

	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)

	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
}
