package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/maxeth/go-account-api/library"
	"github.com/maxeth/go-account-api/model"
	"github.com/maxeth/go-account-api/model/mocks"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user *model.User) {
	user = &model.User{
		UID:      uuid.New(),
		Email:    library.RandomString(15),
		Password: library.RandomString(10),
		Name:     library.RandomString(20),
		ImageURL: library.RandomString(25),
		Website:  library.RandomString(30),
	}
	return
}

func TestGetUser(t *testing.T) {
	id := uuid.New()
	user := model.User{
		UID:      id,
		Email:    "email",
		Password: "pw",
		Name:     "name",
		ImageURL: "url",
		Website:  "website",
	}
	testCases := []struct {
		name          string
		buildStubs    func(repo *mocks.MockUserRepository)
		checkResponse func(t *testing.T, gotUser *model.User, gotError error)
	}{
		{
			name: "OK",
			buildStubs: func(repo *mocks.MockUserRepository) {
				repo.
					EXPECT().
					FindByID(gomock.Any(), gomock.Eq(id)).
					Return(&user, nil)

			},
			checkResponse: func(t *testing.T, gotUser *model.User, gotError error) {
				require.Equal(t, gotUser.UID, id)
				require.NoError(t, gotError)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockUserRepository(ctrl)
			c := UserServiceConfig{
				UserRepository: repo,
			}
			service := NewUserService(&c)
			tc.buildStubs(repo)

			// calling the Get method
			u, err := service.Get(context.Background(), id)
			tc.checkResponse(t, u, err)

		})

	}
}
