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

	// g := c.R.Group("/")

	//if gin.Mode() != gin.TestMode {
	// g.Use(middleware.Timeout(c.TimeOutDuration, model.NewServiceUnavailable()))
	//}
	// sleeping endpoint for testing the timeout middleware
	if gin.Mode() == gin.TestMode {
		c.R.GET("/sleep", func(c *gin.Context) {
			time.Sleep(100000)
		})
	}

	c.R.GET("/me", h.Me)
	c.R.POST("/signup", h.Signup)
	c.R.POST("/signin", h.Signin)
	c.R.POST("/signout", h.Signout)
	c.R.POST("/tokens", h.Tokens)
	c.R.POST("/image", h.Image)
	c.R.DELETE("/image", h.DeleteImage)
	c.R.PUT("/details", h.Details)
}
