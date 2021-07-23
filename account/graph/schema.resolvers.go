package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/maxeth/go-account-api/graph/generated"
	gql_model "github.com/maxeth/go-account-api/graph/model"
	"github.com/maxeth/go-account-api/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) SignUp(ctx context.Context, input gql_model.SignUpDto) (*gql_model.SignUpResponse, error) {
	user, err := r.UserService.Signup(ctx, input.Email, input.Password)
	if err != nil {
		e := &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: err.Error(),
		}
		appError, ok := err.(*model.Error)
		if ok && len(appError.Field) > 0 {
			e.Extensions = map[string]interface{}{
				"field": appError.Field,
				"type":  appError.Type,
			}
		}
		return nil, e
	}

	tokenPair, err := r.TokenService.NewPairFromUser(ctx, user, "")
	if err != nil {
		return nil, &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: "internal error",
		}
	}

	return &gql_model.SignUpResponse{
		Errors:    nil,
		TokenPair: (*gql_model.TokenPair)(tokenPair),
	}, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, input gql_model.SignUpDto) (*gql_model.SignUpResponse, error) {
	user, err := r.UserService.Signin(ctx, input.Email, input.Password)
	if err != nil {
		e := &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: err.Error(),
		}
		return nil, e
	}

	tokenPair, err := r.TokenService.NewPairFromUser(ctx, user, "")
	if err != nil {
		return nil, &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: "internal error",
		}
	}
	return &gql_model.SignUpResponse{
		Errors:    nil,
		TokenPair: (*gql_model.TokenPair)(tokenPair),
	}, nil
}

func (r *queryResolver) Me(ctx context.Context) (*gql_model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id int) (*gql_model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
