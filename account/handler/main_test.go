package handler

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxeth/go-account-api/library"
	"github.com/maxeth/go-account-api/model"
)

func TestMain(m *testing.M) {
	// set gin  to test mode so it doesnt run in debug mode during test
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func randomUser(t *testing.T) (user model.User) {
	user = model.User{
		UID:      uuid.New(),
		Email:    library.RandomString(15),
		Password: library.RandomString(10),
		Name:     library.RandomString(20),
		ImageURL: library.RandomString(25),
		Website:  library.RandomString(30),
	}
	return
}
