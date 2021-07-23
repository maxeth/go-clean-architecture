package graph

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/99designs/gqlgen/graphql"
	"github.com/maxeth/go-account-api/graph/generated"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func stringFromMap(key string, inMap interface{}) (string, error) {
	argsMap, ok := inMap.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("provided input is not a map")
	}

	keyVal, ok := argsMap[key].(string)
	if !ok {
		return "", fmt.Errorf("provided key is not in map")
	}
	return keyVal, nil
}

var SchemaDirectives = generated.DirectiveRoot{
	ValidateEmail: func(ctx context.Context, obj interface{}, next graphql.Resolver, allowDuplicate bool) (res interface{}, err error) {
		email, err := stringFromMap("email", obj)
		if err != nil {
			return nil, err
		}

		if _, err := mail.ParseAddress(email); err != nil {
			return nil, &gqlerror.Error{
				Path:    graphql.GetPath(ctx),
				Message: "invalid email format",
			}
		}

		return next(ctx)
	},
	Length: func(ctx context.Context, obj interface{}, next graphql.Resolver, keyName string, minLength, maxLength int) (res interface{}, err error) {
		arg, err := stringFromMap(keyName, obj)
		if err != nil {
			return nil, err
		}
		if len(arg) < minLength || len(arg) > maxLength {
			return nil, &gqlerror.Error{
				Path:    graphql.GetPath(ctx),
				Message: fmt.Sprintf("%v is too short", keyName),
			}
		}
		return next(ctx)
	},
}
