package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	audience = "be-template-backend"
	issuer   = "be-template-sa"
	subject  = "api-auth"
)

var (
	ErrInvalidTokenClaims           = errors.New("invalid token claims")
	ErrUnexpectedTokenSigningMethod = errors.New("unexpected token signing method")
)

// TokenManager is a JSON web token manager
type tokenManager struct {
	accessTokenSecretKey  string
	refreshTokenSecretKey string
	accessTokenDuration   time.Duration
	refreshTokenDuration  time.Duration
}

type TokenManager interface {
	GenerateAccessToken(ctx context.Context, user *UserForToken) (string, error)
	GenerateRefreshToken(ctx context.Context, user *UserForToken) (string, error)
	Verify(accessToken string) (*CustomClaims, error)
	VerifyRefresh(refreshToken string) (*CustomClaims, error)
}

// UserForToken is bridge for bussinessUser to authUser
type UserForToken struct {
	ID uint64 `json:"id"`
}

// CustomClaims is a custom JWT claims that contains some user's information
type CustomClaims struct {
	jwt.StandardClaims
	UserForToken
}

// NewJWTManager returns a new JWT manager
func NewJWTManager(accessTokenSecretKey, refreshTokenSecretKey string, accessTokenDuration,
	refreshTokenDuration time.Duration) TokenManager {
	return &tokenManager{accessTokenSecretKey, refreshTokenSecretKey,
		accessTokenDuration, refreshTokenDuration}
}

// GenerateAccessToken generates and signs a new token for a user
func (manager *tokenManager) GenerateAccessToken(ctx context.Context, user *UserForToken) (string, error) {

	t := time.Now().UTC()
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			Issuer:    issuer,
			Audience:  audience,
			IssuedAt:  t.Unix(),
			ExpiresAt: t.Add(manager.accessTokenDuration).Unix(),
		},
		UserForToken: *user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.accessTokenSecretKey))
}

// GenerateRefreshToken generates and signs a new token for a user
func (manager *tokenManager) GenerateRefreshToken(ctx context.Context, user *UserForToken) (string, error) {
	t := time.Now().UTC()
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			ExpiresAt: t.Add(manager.refreshTokenDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.refreshTokenSecretKey))
}

// Verify verifies the access token string and return a user claim if the token is valid
func (manager *tokenManager) Verify(accessToken string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrUnexpectedTokenSigningMethod
			}

			return []byte(manager.accessTokenSecretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}

// VerifyRefresh verifies the refresh token string and return new access token if the token is valid
func (manager *tokenManager) VerifyRefresh(refreshToken string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		refreshToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrUnexpectedTokenSigningMethod
			}

			return []byte(manager.refreshTokenSecretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}
