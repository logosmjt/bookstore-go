package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/logosmjt/bookstore-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Name:           util.RandomUserName(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		Role:           "seller",
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Role, user.Role)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestUpdateUserName(t *testing.T) {
	user := createRandomUser(t)
	newName := util.RandomUserName()

	arg := UpdateUserParams{
		ID: user.ID,
		Name: pgtype.Text{
			String: newName,
			Valid:  true,
		},
	}
	updatedUser, err := testStore.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEqual(t, user.Name, updatedUser.Name)
	require.Equal(t, newName, updatedUser.Name)
	require.Equal(t, user.Email, user.Email)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserEmail(t *testing.T) {
	user := createRandomUser(t)
	newEmail := util.RandomEmail()

	arg := UpdateUserParams{
		ID: user.ID,
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
	}
	updatedUser, err := testStore.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEqual(t, user.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, user.Name, user.Name)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)

	newEmail := util.RandomEmail()
	newName := util.RandomUserName()
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := UpdateUserParams{
		ID: user.ID,
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
		Name: pgtype.Text{
			String: newName,
			Valid:  true,
		},
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
	}
	updatedUser, err := testStore.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEqual(t, user.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, user.Name, updatedUser.Name)
	require.Equal(t, newName, updatedUser.Name)
	require.NotEqual(t, user.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
}
