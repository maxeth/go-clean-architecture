package handler

import (
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/maxeth/go-account-api/graph"
	"github.com/maxeth/go-account-api/graph/generated"
	"github.com/maxeth/go-account-api/handler/middleware"
	"github.com/maxeth/go-account-api/model"
)

type Handler struct {
	UserService     model.UserService
	TokenService    model.TokenService
	TimeOutDuration time.Duration
	OAuthService    model.OAuthService
}

type Config struct {
	R               *gin.Engine // type of the gin/http router
	UserService     model.UserService
	TokenService    model.TokenService
	TimeOutDuration time.Duration
	OAuthService    model.OAuthService
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// returns the handler function for the graphql-playground end point
func graphqlHandler(c *Config) gin.HandlerFunc {
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

func applyMiddleware(c *Config) {
	c.R.Use(middleware.Timeout(c.TimeOutDuration, model.NewInternal()))
	c.R.Use(middleware.Cors("*"))
}

func newGraphqlHandler(c *Config) {
	fmt.Println("gql handler being init'ed")

	c.R.POST("/graphql", graphqlHandler(c))
	c.R.GET("/playground", playgroundHandler())
}

func NewHandler(c *Config) {
	h := &Handler{
		TokenService:    c.TokenService,
		UserService:     c.UserService,
		TimeOutDuration: c.TimeOutDuration,
		OAuthService:    c.OAuthService,
	}

	noMd := c.R.Group("/")
	noMd.GET("/auth/test", func(c *gin.Context) {
		c.Redirect(301, "http://www.google.com/test")
		c.Abort()
	})
	noMd.GET("/auth", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "http://www.google.com/test"})
	})
	noMd.GET("/auth/twitch/callback", h.SigninTwitch)
	noMd.GET("/auth/twitch", h.RedirectTwitch)

	g := c.R.Group("/")

	if gin.Mode() == gin.TestMode {
		g.GET("/sleep", func(c *gin.Context) {
			time.Sleep(100000) // in order to e2e test the timeout middleware
		})
	}

	g.Use(middleware.Timeout(c.TimeOutDuration, model.NewInternal()))
	g.Use(middleware.Cors("*"))

	g.GET("/me", h.Me)
	g.POST("/signup", h.Signup)
	g.POST("/signin", h.Signin)
	g.POST("/signout", h.Signout)
	g.POST("/tokens", h.Tokens)
	g.POST("/image", h.Image)
	g.DELETE("/image", h.DeleteImage)
	g.PUT("/details", h.Details)

	gql := c.R.Group("/")

	gql.Use(middleware.Timeout(c.TimeOutDuration, model.NewInternal()))
	gql.Use(middleware.Cors("*"))

	gql.POST("/graphql", graphqlHandler(c))
	gql.GET("/playground", playgroundHandler())
}
