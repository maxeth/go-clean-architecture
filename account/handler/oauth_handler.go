package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxeth/go-account-api/model"
)

func (h *Handler) RedirectTwitch(c *gin.Context) {
	url := h.OAuthService.GetTwitchRedirectURL()
	c.Redirect(http.StatusMovedPermanently, url)
	c.Abort()
}

func (h *Handler) SigninTwitch(c *gin.Context) {
	code := c.Query("code")
	if len(code) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": model.NewBadRequest("Expected an oauth code as 'code' query parameter."),
		})
		return
	}

	twitchOICD, err := h.OAuthService.GetTwitchCredentials(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": twitchOICD})

}
