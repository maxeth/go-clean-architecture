package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxeth/go-account-api/model"
)

type Handler struct {
	UserService     model.UserService
	TokenService    model.TokenService
	TimeOutDuration time.Duration
}

type Config struct {
	R               *gin.Engine // type of the gin/http router
	UserService     model.UserService
	TokenService    model.TokenService
	TimeOutDuration time.Duration
}

func NewHandler(c *Config) {
	h := &Handler{
		TokenService:    c.TokenService,
		UserService:     c.UserService,
		TimeOutDuration: c.TimeOutDuration,
	}

	g := c.R.Group("/account")

	//if gin.Mode() != gin.TestMode {
	// g.Use(middleware.Timeout(c.TimeOutDuration, model.NewServiceUnavailable()))
	//}
	// sleeping endpoint for testing the timeout middleware
	if gin.Mode() == gin.TestMode {
		g.GET("/sleep", func(c *gin.Context) {
			time.Sleep(100000)
		})
	}

	g.GET("/me", h.Me)
	g.POST("/signup", h.Signup)
	g.POST("/signin", h.Signin)
	g.POST("/signout", h.Signout)
	g.POST("/tokens", h.Tokens)
	g.POST("/image", h.Image)
	g.DELETE("/image", h.DeleteImage)
	g.PUT("/details", h.Details)
}
