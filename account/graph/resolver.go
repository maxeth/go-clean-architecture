package graph

import "github.com/maxeth/go-account-api/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService  model.UserService
	TokenService model.TokenService
}
