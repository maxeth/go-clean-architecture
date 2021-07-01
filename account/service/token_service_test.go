package service

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPairFromUser(t *testing.T) {
	privFile, err := ioutil.ReadFile("../rsa_private_test.pem")
	require.NoError(t, err)
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privFile)
	require.NoError(t, err)

	pubFile, err := ioutil.ReadFile("../rsa_public_test.pem")
	require.NoError(t, err)
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubFile)
	require.NoError(t, err)

	secret := "secret1sdsadasdasdasdasda23"

	issuedAt := time.Now()
	atExpiry := issuedAt.Add(15 * time.Minute)    // expected access token expiry
	rtExpiry := issuedAt.Add(30 * time.Hour * 24) // expected refresh token expiry

	tsc := &TokenServiceConfig{
		PrivKey:             privKey,
		PubKey:              pubKey,
		RefreshSecret:       secret,
		AccessTokenExpSecs:  60 * 15,           // 15 min
		RefreshTokenExpSecs: 60 * 60 * 24 * 30, // 30 days
	}
	tokenService := NewTokenService(tsc)

	user := randomUser(t)
	// generate a token pair for the user
	tokenPair, err := tokenService.NewPairFromUser(context.Background(), user, "")
	require.NoError(t, err)

	// assert tokens are a string
	var s string
	assert.IsType(t, s, tokenPair.AccessToken, tokenPair.RefreshToken)

	var accessTokenClaims AccessTokenClaims

	at, err := jwt.ParseWithClaims(tokenPair.AccessToken, &accessTokenClaims, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	require.NoError(t, err)

	// try casting received AT claims struct to the AccesTokenClaims struct
	atClaims, ok := at.Claims.(*AccessTokenClaims)
	require.True(t, ok)
	require.NotEmpty(t, atClaims.ExpiresAt, atClaims.Id, atClaims.IssuedAt, atClaims.User.Email, atClaims.User.UID)
	require.Empty(t, atClaims.User.Password)

	// check whether RT and AT expiry is as expetcted (15min, 30days)
	require.WithinDuration(t, time.Unix(atClaims.ExpiresAt, 0), atExpiry, time.Second)
	require.WithinDuration(t, time.Unix(atClaims.IssuedAt, 0), issuedAt, time.Second)

	// check rt token
	var refreshTokenClaims RefreshTokenClaims
	rt, err := jwt.ParseWithClaims(tokenPair.RefreshToken, &refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	// try casting received RT claims struct to the RefreshTokenClaims struct
	rtClaims, ok := rt.Claims.(*RefreshTokenClaims)
	require.True(t, ok)
	require.NoError(t, err)

	// ensure the UID of the RT is of type UUID
	var uuid uuid.UUID
	assert.IsType(t, rtClaims.UID, uuid)

	require.WithinDuration(t, time.Unix(rtClaims.ExpiresAt, 0), rtExpiry, time.Second)
	require.WithinDuration(t, time.Unix(rtClaims.IssuedAt, 0), issuedAt, time.Second)
}
