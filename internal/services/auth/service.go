package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ash543210/go-chat/mongo"
	"github.com/ash543210/go-chat/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("super-secret-key") // store this securely (env var)

// CustomClaims defines your JWT claims structure
type CustomClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

type AuthService struct {
	Logger *logger.Logger
}

// GenerateAuthTokens creates an access and refresh token
func (a *AuthService) GenerateAuthTokens(userID string, ctx *context.Context) (string, string, error) {
	// Access Token: 15 minutes expiry
	accessTokenClaims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := generateJWT(accessTokenClaims)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken := uuid.NewString()

	return accessToken, refreshToken, nil
}

// Helper to sign JWT
func generateJWT(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		// Check for validation errors
		switch err {
		case jwt.ErrTokenExpired:
			return nil, jwt.ErrTokenExpired
		default:
			return nil, jwt.ErrTokenUnverifiable
		}
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenUnverifiable
	}

	return claims, nil
}

func (s *AuthService) InsertRefreshToken(userID string, refreshToken string, ctx context.Context) error {
	collection, err := mongo.GetCollection("refresh_tokens")
	if err != nil {
		s.Logger.Error(ctx, "Error while getting colection: %s", err)
		return nil
	}
	_, err = collection.InsertOne(ctx, RefreshToken{
		UserID:       userID,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now().Unix(),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	if err != nil {
		s.Logger.Error(ctx, "Error saving refresh token: %s", err)
		return nil
	}
	s.Logger.Info(ctx, "Refresh token saved for user: %s", userID)
	return nil
}
