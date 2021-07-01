package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/maxeth/go-account-api/model"
)

func basicErrorResponse(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{
		"error": err,
	})
}

func errorResponse(c *gin.Context, customError model.Error) {
	c.JSON(customError.Status(), gin.H{
		"error": customError,
	})
}
