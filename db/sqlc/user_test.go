package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nhan-ngo-usf/NBank/util"
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

func TestUpdateFullName(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newFullName := util.RandomUserName()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		FullName: sql.NullString{
			String: newFullName, 
			Valid: true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
}

func TestUpdateEmail(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newEmail := util.RandomEmail()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Email: sql.NullString{
			String: newEmail, 
			Valid: true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
}

func TestUpdateUserPassword(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword,err := util.HashPassword(newPassword)
	require.NoError(t, err)

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword, 
			Valid: true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.Email, newUser.Email)
}

func TestUpdateAll(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newPassword := util.RandomString(6)
	newFullName := util.RandomUserName()
	newEmail := util.RandomEmail()
	
	newHashedPassword,err := util.HashPassword(newPassword)
	require.NoError(t, err)

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword, 
			Valid: true,
		},
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid: true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid: true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, oldUser.Username, newUser.Username)
	require.NotEqual(t, oldUser.FullName, newUser.FullName)
	require.NotEqual(t, oldUser.Email, newUser.Email)
}