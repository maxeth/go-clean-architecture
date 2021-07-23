package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/maxeth/go-account-api/model"
)

type userService struct {
	UserRepository model.UserRepository
}

type UserServiceConfig struct {
	UserRepository model.UserRepository
}

func NewUserService(c *UserServiceConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

func (us *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := us.UserRepository.FindByID(ctx, uid)
	return u, err
}

func (us *userService) Signup(ctx context.Context, email, password string) (*model.User, error) {
	empty := &model.User{}

	hashedPw, err := HashPassword(password)
	if err != nil {
		return empty, model.NewInternal()
	}

	// create a user struct with the passed email and hashed password to be stored in the db
	u := &model.User{
		Email:    email,
		Password: hashedPw,
	}
	user, err := us.UserRepository.Create(ctx, u)
	if err != nil {
		return empty, err
	}

	return user, err
}

func (us *userService) Signin(ctx context.Context, email, password string) (*model.User, error) {
	empty := &model.User{}

	user, err := us.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		return empty, model.NewInternal()
	}

	if user == nil {
		return empty, model.NewNotFound("user", email)
	}

	if err := ComparePassword(user.Password, password); err != nil {
		return empty, model.NewAuthorization("password and email do not match")
	}

	return user, nil
}
