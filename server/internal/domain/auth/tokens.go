package auth

import (
	"time"

	"github.com/bug-breeder/2fair/server/internal/domain/models"
	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your_secret_key")

// GenerateAccessToken generates a JWT access token
func GenerateAccessToken(userID, email string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &models.Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken generates a JWT refresh token
func GenerateRefreshToken(userID, email string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &models.Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenStr string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}