package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/maxeth/go-account-api/library"
	"github.com/maxeth/go-account-api/model"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func TestMain(m *testing.M) {

	dbConn, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres dbname=accounts_db password=password sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic("cannot connect to db in test mode")
	}
	db = dbConn

	os.Exit(m.Run())
}

func randomUser() (user *model.User) {
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
