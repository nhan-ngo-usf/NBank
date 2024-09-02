package db

import (
	"context"
	"testing"
	"time"

	"github.com/nhan-ngo-usf/NBank/db/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T)(User) {
	HashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams {
		Username: util.RandomUserName(),
		HashedPassword: HashedPassword,
		FullName: util.RandomUserName(),
		Email: util.RandomEmail(),
	}
	
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user) 

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())
	return user
}
func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	randomUser := CreateRandomUser(t)
	user, err := testQueries.GetUser(context.Background(), randomUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, randomUser.Username)
	require.Equal(t, user.HashedPassword, randomUser.HashedPassword)
	require.Equal(t, user.FullName, randomUser.FullName)
	require.Equal(t, user.Email, randomUser.Email)
	
	require.WithinDuration(t, user.CreatedAt, randomUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, randomUser.PasswordChangedAt, time.Second)
}