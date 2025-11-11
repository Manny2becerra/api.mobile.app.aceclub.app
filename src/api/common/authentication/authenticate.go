package auth

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token format")
	ErrExpiredToken = errors.New("token has expired")
)

// AccessClaims represents the JWT claims expected by the API
// EstablishmentID is optional and present for establishment users.
type AccessClaims struct {
	jwt.RegisteredClaims
	Role            string `json:"role"`
	EstablishmentID string `json:"establishmentId,omitempty"`
	UserID          string `json:"userId,omitempty"`
}

// Authenticate validates a JWT access token and returns the parsed claims
func Authenticate(accessToken string, jwtSecret string) (*AccessClaims, error) {
	if !strings.HasPrefix(accessToken, "Bearer ") {
		return nil, ErrInvalidToken
	}

	// Remove "Bearer " prefix
	tokenString := strings.TrimPrefix(accessToken, "Bearer ")

	var claims AccessClaims
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	// Extract claims
	if c, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		if c.Subject == "" || c.Role == "" {
			return nil, ErrInvalidToken
		}
		return c, nil
	}

	return nil, ErrInvalidToken
}
