package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	oauth "github.com/maxeth/go-account-api/model/oauth"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement
type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	Signup(ctx context.Context, email, password string) (*User, error)
	Signin(ctx context.Context, email, password string) (*User, error)
}

type TokenService interface {
	NewPairFromUser(ctx context.Context, u *User, prevTokenID string) (*TokenPair, error)
}

type OAuthService interface {
	SignupTwitch(ctx context.Context, tokenRes oauth.TwitchOIDCResponse) error
}

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) (*User, error)
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error
}
