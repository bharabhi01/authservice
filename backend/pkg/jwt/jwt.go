package jwt

import (
	"errors"
	"time"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// Claims represents the JWT claims structure
// It extends the standard JWT claims with our custom fields
type Claims struct {
	UserId string `json:"user_id"`
	Username string `json:"username"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userId, username, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	// Set expiration time to 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims
	claims := &Claims {
		UserID: userId,
		Username: username,
		Role: role,
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: expirationTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "authservice"
		},
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	// Parse the token with the secret key
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

