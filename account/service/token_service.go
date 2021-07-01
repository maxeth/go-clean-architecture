package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/maxeth/go-account-api/model"
)

// TokenService used for injecting an implementation of TokenRepository
// for use in service methods along with keys and secrets for
// signing JWTs
type tokenService struct {
	TokenRepository     model.TokenRepository
	PrivKey             *rsa.PrivateKey
	PubKey              *rsa.PublicKey
	RefreshSecret       string
	AccessTokenExpSecs  int64
	RefreshTokenExpSecs int64
}

// TSConfig will hold repositories that will eventually be injected into this
// this service layer
type TokenServiceConfig struct {
	TokenRepository     model.TokenRepository
	PrivKey             *rsa.PrivateKey
	PubKey              *rsa.PublicKey
	RefreshSecret       string
	AccessTokenExpSecs  int64
	RefreshTokenExpSecs int64
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(c *TokenServiceConfig) model.TokenService {
	return &tokenService{
		TokenRepository:     c.TokenRepository,
		PrivKey:             c.PrivKey,
		PubKey:              c.PubKey,
		RefreshSecret:       c.RefreshSecret,
		AccessTokenExpSecs:  c.AccessTokenExpSecs,
		RefreshTokenExpSecs: c.RefreshTokenExpSecs,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
func (s *tokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	// No need to use a repository for idToken as it is unrelated to any data source
	accessToken, err := generateAccessToken(u, s.PrivKey, s.AccessTokenExpSecs)
	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, model.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret, s.RefreshTokenExpSecs)
	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, model.NewInternal()
	}

	// save the refresh token associated to this user id in redis.
	if err := s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("error saving refresh token in redis: %v\n", err.Error())
		return nil, model.NewInternal()
	}

	// delete the users previous refresh token from redis if an prevTokenID was provided
	if len(prevTokenID) > 0 {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			log.Printf("error deleting user's previous refresh token in redis: %v\n", err.Error())
			return nil, model.NewInternal()
		}
	}

	// TODO: store refresh tokens by calling TokenRepository methods

	tp := &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.SignedRefreshToken,
	}
	return tp, nil
}
