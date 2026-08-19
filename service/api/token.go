package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// Secret key for signing JWT tokens
var jwtKey = []byte("your-very-secret-key-change-this-in-production")

// tokenClaims represents the JWT payload
type tokenClaims struct {
	Sub string `json:"sub"` // Subject (user ID)
	Exp int64  `json:"exp"` // Expiration time
}

// sign creates an HMAC signature for the given data
func sign(data string) []byte {
	h := hmac.New(sha256.New, jwtKey)
	h.Write([]byte(data))
	return h.Sum(nil)
}

// createToken generates a new JWT token for the given user ID
func createToken(userID string) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	claims := tokenClaims{
		Sub: userID,
		Exp: time.Now().Add(24 * time.Hour).Unix(),
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	sig := sign(header + "." + payload)
	signature := base64.RawURLEncoding.EncodeToString(sig)

	return header + "." + payload + "." + signature, nil
}

// ParseToken validates and parses a JWT token, returning the user ID
func ParseToken(tokenString string) (string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", ErrUnauthorized
	}

	header, payload, sigB64 := parts[0], parts[1], parts[2]

	sig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return "", ErrUnauthorized
	}

	expected := sign(header + "." + payload)
	if !hmac.Equal(sig, expected) {
		return "", ErrUnauthorized
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", ErrUnauthorized
	}

	var claims tokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return "", ErrUnauthorized
	}

	if claims.Exp < time.Now().Unix() {
		return "", ErrUnauthorized
	}

	if claims.Sub == "" {
		return "", ErrUnauthorized
	}

	return claims.Sub, nil
}
