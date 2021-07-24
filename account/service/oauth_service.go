package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/maxeth/go-account-api/model"
)

type oAuthService struct {
	TwitchAPIUrl string
	Secret       string
	ClientID     string
	Callback_URI string
}

type OAuthServiceConfig struct {
	TwitchAPIUrl string
	Secret       string
	ClientID     string
	Callback_URI string
}

func NewOAuthService(c *OAuthServiceConfig) model.OAuthService {
	return &oAuthService{
		TwitchAPIUrl: c.TwitchAPIUrl,
		Secret:       c.Secret,
		ClientID:     c.ClientID,
		Callback_URI: c.Callback_URI,
	}
}

func (s *oAuthService) GetTwitchRedirectURL() string {
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=%v&redirect_uri=%v&scope=user:read:email+openid&state=c3ab8aa609ea11e793ae92361f002671&claims={\"id_token\":{\"email\":null}}", s.ClientID, s.Callback_URI)

	return url
}

func (s *oAuthService) GetTwitchCredentials(code string) (model.TwitchOIDCResponse, error) {
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%v&client_secret=%v&grant_type=authorization_code&redirect_uri=%v&code=%v", s.ClientID, s.Secret, s.Callback_URI, code)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return model.TwitchOIDCResponse{}, model.NewInternal()

	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return model.TwitchOIDCResponse{}, model.NewInternal()
	}
	defer resp.Body.Close()

	var twitchOIDC model.TwitchOIDCResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.TwitchOIDCResponse{}, model.NewInternal()
	}

	err = json.Unmarshal(b, &twitchOIDC)
	if err != nil {
		return model.TwitchOIDCResponse{}, model.NewInternal()

	}

	return twitchOIDC, nil
}
