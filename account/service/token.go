package service

import (
	"crypto/rsa"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/maxeth/go-account-api/model"
)

type AccessTokenClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

func generateAccessToken(u *model.User, key *rsa.PrivateKey, exp int64) (string, error) {
	unixTime := time.Now().Unix()
	expTime := unixTime + exp // 15 min

	claims := &AccessTokenClaims{
		User: u,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: expTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	ss, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return ss, nil
}

// the refresh token holds the jwt signed string token
type RefreshToken struct {
	SignedRefreshToken string        // signed refresh token string that is beign  returned to the user
	ID                 string        // tokens unoque id,  used for internal utility
	ExpiresIn          time.Duration // used for internal utility
}

// claims of the refresh token whose signed string is contained in the RefreshToken struct
type RefreshTokenClaims struct {
	UID uuid.UUID `json:"uid"` // the users uuid
	jwt.StandardClaims
}

func generateRefreshToken(uid uuid.UUID, key string, exp int64) (*RefreshToken, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp * int64(time.Second))) // 30 days

	tokenID, err := uuid.NewRandom()

	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	claims := &RefreshTokenClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println("Failed to sign refresh token", err)
		return nil, err
	}

	rt := &RefreshToken{
		SignedRefreshToken: ss, // This is the only thing that is actually being returned to the user inside the service. the other properties are just for internal utility
		ID:                 tokenID.String(),
		ExpiresIn:          tokenExp.Sub(currentTime),
	}

	return rt, nil
}
