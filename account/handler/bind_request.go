package handler

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/maxeth/go-account-api/model"
)

// used to help extract validation errors of http request body
type InvalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// type of the http reponse when sending an invalid request
type InvalidRequestResponse struct {
	Error       model.Error       `json:"error"`
	InvalidArgs []InvalidArgument `json:"invalidArgs"`
}

// bindData tries to bind the expected request struct to the actual http request body inside gin.Context
// if binding is successful return true. Else, send an http error and return false
func bindData(c *gin.Context, req interface{}) bool {
	//b, _ := ioutil.ReadAll(c.Request.Body)
	//fmt.Println("inside binder: ", string(b))

	fmt.Println("request interface: ", req)
	// return error of request header is of any other type than json
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("Content-Type for %s must be application/json", c.FullPath())
		err := model.NewUnsupportedMediaType(msg)

		errorResponse(c, *err)
		return false
	}

	// Bind incoming json to struct and check for validation errors
	err := c.ShouldBind(req)
	if err == nil {
		return true
	}

	log.Printf("Error binding data: %+v\n", err)

	if errs, ok := err.(validator.ValidationErrors); ok {
		// validation error from the validation library that gin uses
		var invalidArgs []InvalidArgument

		for _, err := range errs {
			invalidArgs = append(invalidArgs, InvalidArgument{
				err.Field(),
				err.Value().(string),
				err.Tag(),
				err.Param(),
			})
		}

		err := model.NewBadRequest("Invalid request parameters. See invalidArgs")
		c.JSON(err.Status(), gin.H{
			"error":       err,
			"invalidArgs": invalidArgs,
		})

		return false
	}

	// later we'll add code for validating max body size here!

	// if we aren't able to properly extract validation errors,
	// we'll fallback and return an internal server error
	fallBack := model.NewInternal()
	c.JSON(fallBack.Status(), gin.H{"error": fallBack})

	return false
}
