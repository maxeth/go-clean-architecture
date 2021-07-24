package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxeth/go-account-api/model"
)

func (h *Handler) Me(c *gin.Context) {

	user, ok := c.Get("user")
	if !ok {
		// no user in context
		errM := model.NewAuthorization("not signed in")
		errorResponse(c, *errM)
		return
	}

	uid := user.(*model.User).UID

	ctx := c.Request.Context()
	user, err := h.UserService.Get(ctx, uid)
	if err != nil {
		errM := model.NewNotFound("user", uid.String())
		errorResponse(c, *errM)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

type signupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"` // 6 <= password <= 30
}

func (h *Handler) Signup(c *gin.Context) {
	//b, _ := ioutil.ReadAll(c.Request.Body)
	//	fmt.Println("GOT IN HANDLER BODY: ", string(b))

	var req signupReq
	if ok := bindData(c, &req); !ok {
		// request binding error, error message already sent to c
		return
	}

	ctx := c.Request.Context()
	// will create a user with email and a password. rest of the fields will remain empty strings
	user, err := h.UserService.Signup(ctx, req.Email, req.Password)
	if err != nil {
		basicErrorResponse(c, model.Status(err), err)
		return
	}

	tokenPair, err := h.TokenService.NewPairFromUser(ctx, user, "") // dont pass any prevTokenID because this is a (first) signup and not signin
	if err != nil {
		basicErrorResponse(c, model.Status(err), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokenPair,
	})
}

type signinReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"` // 6 <= password <= 30
}

func (h *Handler) Signin(c *gin.Context) {
	var req signinReq
	if ok := bindData(c, &req); !ok {
		return
	}

	// extract the "actual" http context as this is the context we want to pass down the callchain in every handler, not the gin context
	ctx := c.Request.Context()

	user, err := h.UserService.Signin(ctx, req.Email, req.Password)
	if err != nil {
		errM := model.NewAuthorization("Invalid password or email.")
		errorResponse(c, *errM)
	}

	tokens, err := h.TokenService.NewPairFromUser(ctx, user, "")
	if err != nil {
		log.Printf("Failed to create tokens when signing in user: %v\n", err.Error())
		errM := model.NewInternal()
		errorResponse(c, *errM)
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}

// Signout handler
func (h *Handler) Signout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's signout",
	})
}

// Tokens handler
func (h *Handler) Tokens(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's tokens",
	})
}

// Image handler
func (h *Handler) Image(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's image",
	})
}

// DeleteImage handler
func (h *Handler) DeleteImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's deleteImage",
	})
}

// Details handler
func (h *Handler) Details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's details",
	})
}
