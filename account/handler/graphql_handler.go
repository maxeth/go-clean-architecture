package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxeth/go-account-api/graph"
	"github.com/maxeth/go-account-api/graph/generated"
	"github.com/maxeth/go-account-api/handler/middleware"
	"github.com/maxeth/go-account-api/model"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type GraphQLHandler struct {
	UserService     model.UserService
	TokenService    model.TokenService
	TimeOutDuration time.Duration
}

type GraphQLConfig struct {
	R               *gin.Engine // type of the gin/http router
	UserService     model.UserService
	TokenService    model.TokenService
	TimeOutDuration time.Duration
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// returns the handler function for the graphql-playground end point
func graphqlHandler(c *GraphQLConfig) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	generatedConfig := generated.Config{
		Resolvers: &graph.Resolver{
			UserService:  c.UserService,
			TokenService: c.TokenService,
		},
		Directives: graph.SchemaDirectives,
	}

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generatedConfig))
	if h == nil {
		panic("GraphQL handlerfunction is nil.")
	}

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func NewGraphQLHandler(c *GraphQLConfig) {
	fmt.Println("gql handler being init'ed")

	c.R.Use(middleware.Timeout(c.TimeOutDuration, model.NewInternal()))

	c.R.POST("/graphql", graphqlHandler(c))
	c.R.GET("/playground", playgroundHandler())
}
