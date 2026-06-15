// Package auth provides functions for password hashing, JWT token generation and validation, and extracting Bearer tokens from HTTP headers.
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
)

// Create hashed password from plain text password
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	// WARNING: argon2id.DefaultParams are a good starting point,
	// but you may want to adjust them based on your security requirements and performance needs
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return hash, nil
}

// Check if the provided password matches the hashed password
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("failed to compare password and hash: %w", err)
	}
	return match, nil
}

// Generate a JWT token for the given user ID
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Issuer:    "crawltrip-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error signing JWT: %w", err)
	}
	return ss, nil
}

// Generate a random refresh token
func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}

// Validate the provided JWT token and return the user ID if valid
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing JWT: %w", err)
	} else if claim, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		if claim.ExpiresAt.Time.Before(time.Now().UTC()) {
			return uuid.Nil, fmt.Errorf("JWT token has expired")
		}
		userID, err := uuid.Parse(claim.Subject)
		if err != nil {
			return uuid.Nil, fmt.Errorf("error parsing user ID from JWT claims: %w", err)
		}
		return userID, nil
	} else {
		return uuid.Nil, fmt.Errorf("error getting subject from JWT claims: %w", err)
	}
}

// Extract the Bearer token from the Authorization header
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) < 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("malformed authorization header")
	}
	return parts[1], nil
}
