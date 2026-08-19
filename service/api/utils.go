package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ematri2108281/wasatext/service/components/schema"
)

var ErrUnauthorized = errors.New("unauthorized")

// getAuthenticatedUserID extracts and validates the user ID from the Authorization header
func (rt *_router) getAuthenticatedUserID(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrUnauthorized
	}

	tokenString := authHeader[7:]

	// Parse and validate token
	userID, err := ParseToken(tokenString)
	if err != nil {
		return "", ErrUnauthorized
	}

	// Verify user exists in database
	if _, err := rt.db.GetUserById(userID); err != nil {
		if errors.Is(err, schema.ErrUserDoesNotExist) {
			return "", ErrUnauthorized
		}
		return "", err
	}

	return userID, nil
}
