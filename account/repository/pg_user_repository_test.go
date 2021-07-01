package repository

import (
	"context"
	"testing"

	"github.com/maxeth/go-account-api/library"
	"github.com/maxeth/go-account-api/model"
	"github.com/stretchr/testify/require"
)

func randomCreateUser() *model.User {
	u := &model.User{
		Email:    library.RandomString(8),
		Password: library.RandomString(10),
	}

	return u
}

func TestCreateUser(t *testing.T) {
	okUser := randomCreateUser()
	duplicateUser := &model.User{
		Email: okUser.Email,
	}

	repo := NewUserRepository(db)
	gotUser, err := repo.Create(context.Background(), okUser)
	require.NoError(t, err)
	require.Equal(t, okUser.Email, gotUser.Email)
	require.Equal(t, okUser.Name, gotUser.Name)
	require.NotEmpty(t, gotUser.UID)
	require.Empty(t, gotUser.Website, gotUser.ImageURL)

	gotUser2, err := repo.Create(context.Background(), duplicateUser)

	errM, ok := err.(*model.Error)

	require.True(t, ok)
	require.Equal(t, 409, errM.Status())
	require.Empty(t, gotUser2)
}
