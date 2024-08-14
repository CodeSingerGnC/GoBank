package db

import (
	"context"
	"testing"
	"time"

	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User{
	userAccount := util.RandomUser()
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	
	arg := CreateUserParams{
		UserAccount: userAccount,
		HashPassword: hashedPassword,
		Username: util.RandomString(10),
		Email: util.RandomEmail(),
	}

	result, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	rowsAffected, _ := result.RowsAffected()
	require.True(t, rowsAffected == 1)

	var user User
	query := `SELECT     
				user_account, 
				hash_password, 
				username, 
				email, 
				password_chaged_at, 
				created_at 
			FROM users 
			WHERE user_account = ?`

	err = testdb.QueryRow(query, userAccount).
		Scan(&user.UserAccount, 
			&user.HashPassword,
			&user.Username,
			&user.Email,
			&user.PasswordChagedAt,
			&user.CreatedAt)
			
	require.NoError(t, err)
	require.Equal(t, userAccount, user.UserAccount)
	require.Equal(t, arg.HashPassword, user.HashPassword)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChagedAt.Equal(time.Unix(1,0)))
	require.NotZero(t, user.CreatedAt)

	return user
}


func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.UserAccount)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserAccount, user2.UserAccount)
	require.Equal(t, user1.HashPassword, user2.HashPassword)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChagedAt, user2.PasswordChagedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}