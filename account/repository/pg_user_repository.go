package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/maxeth/go-account-api/model"
)

type pgUserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &pgUserRepository{
		DB: db,
	}
}

func (r *pgUserRepository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	q := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"

	user := &model.User{}
	if err := r.DB.GetContext(ctx, user, q, u.Email, u.Password); err != nil {
		fmt.Println("got error when creating user:", err)
		// check whether its a unique constrain viloation pg error
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			errM := model.NewConflict("email", u.Email)
			return &model.User{}, errM
		}

		// something else went wrong
		return &model.User{}, model.NewInternal()
	}

	fmt.Println("repo returning user: ", user)
	return user, nil
}

func (r *pgUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	q := "SELECT * FROM users u WHERE id = $1 LIMIT 1"

	user := &model.User{}
	if err := r.DB.GetContext(ctx, user, q, uid); err != nil {
		return &model.User{}, model.NewInternal()
	}

	return user, nil
}

func (r *pgUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	q := "SELECT * FROM users u WHERE email = $1 LIMIT 1"

	user := &model.User{}
	if err := r.DB.GetContext(ctx, user, q, email); err != nil {
		return &model.User{}, model.NewInternal()
	}

	return user, nil
}
