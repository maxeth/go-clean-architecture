package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/maxeth/go-account-api/handler"
	"github.com/maxeth/go-account-api/repository"
	"github.com/maxeth/go-account-api/service"
)

// will initialize a handler starting from data sources
// which inject into repository layer
// which inject into service layer
// which inject into handler layer
func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	userRepository := repository.NewUserRepository(d.DB)
	userService := service.NewUserService(&service.UserServiceConfig{
		UserRepository: userRepository,
	})

	// load rsa keys and config vars for tokenrepository
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := ioutil.ReadFile(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// load refresh token secret and expiry seconds from env variables
	refreshSecret := os.Getenv("REFRESH_SECRET")
	accessTokenExp := os.Getenv("ACCESS_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	accessTokenExpSecs, err := strconv.ParseInt(accessTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could parse access token exp: %w", err)
	}
	refreshtokenExpSecs, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could parse refresh token exp: %w", err)

	}

	tokenRepository := repository.NewTokenRepository(d.RedisClient)
	tokenService := service.NewTokenService(&service.TokenServiceConfig{
		TokenRepository:     tokenRepository,
		PrivKey:             privKey,
		PubKey:              pubKey,
		RefreshSecret:       refreshSecret,
		RefreshTokenExpSecs: refreshtokenExpSecs,
		AccessTokenExpSecs:  accessTokenExpSecs,
	})

	// initialize gin.Engine
	router := gin.Default()

	tc := &service.OAuthServiceConfig{
		Secret:       os.Getenv("TWITCH_SECRET"),
		ClientID:     os.Getenv("TWITCH_CLIENT"),
		Callback_URI: os.Getenv("TWITCH_CALLBACK"),
	}
	oAuthService := service.NewOAuthService(tc)

	c := &handler.Config{
		R:               router,
		UserService:     userService,
		OAuthService:    oAuthService,
		TokenService:    tokenService,
		TimeOutDuration: time.Duration(7 * time.Second),
	}
	handler.NewHandler(c)
	//handler.NewGraphQLHandler(c)

	return router, nil
}
