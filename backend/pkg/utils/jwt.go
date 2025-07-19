package utils

import (
	"fmt"
	"time"

	"document-reader-chatbot/configs"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenClaims represents the JWT claims structure
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(cfg configs.JWTConfig, userID, email, role string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    cfg.Issuer,
			Subject:   userID,
			Audience:  []string{"document-reader-chatbot"},
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.ExpirationTime)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(cfg configs.JWTConfig, tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// RefreshToken generates a new token with updated expiration time
func RefreshToken(cfg configs.JWTConfig, oldToken string) (string, error) {
	claims, err := ValidateToken(cfg, oldToken)
	if err != nil {
		return "", fmt.Errorf("failed to validate old token: %w", err)
	}

	// Generate new token with same claims but updated timestamps
	return GenerateToken(cfg, claims.UserID, claims.Email, claims.Role)
}
