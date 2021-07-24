package model

type TwitchOIDCResponse struct {
	AccessToken  string   `json:"access_token"`  // a non JWT access token for the signed in user
	RefreshToken string   `json:"refresh_token"` // a refresh token for the AccessToken
	IdToken      string   `json:"id_token"`      // JWT that includes default claims: iss, sub, aud, exp, iat, nonce + all requested claims // https://dev.twitch.tv/docs/authentication/getting-tokens-oidc , IDToken cannot be refreshed. if expired, re-request data from users api
	ExpiresIn    int64    `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}
