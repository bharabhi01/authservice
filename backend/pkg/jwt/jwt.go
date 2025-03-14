package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Global configuration for JWT
var (
	jwtSecret      string
	expirationHours int
)

// Init initializes the JWT package with configuration
// This should be called during application startup
func Init(secret string, expHours int) {
	jwtSecret = secret
	expirationHours = expHours
}

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
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	// Set expiration time to 24 hours
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// Create the JWT claims
	claims := &Claims {
		UserId: userId,
		Username: username,
		Role: role,
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: expirationTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "authservice",
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

