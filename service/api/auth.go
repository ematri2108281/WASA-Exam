package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ematri2108281/wasatext/service/api/reqcontext"
	"github.com/ematri2108281/wasatext/service/components/requests"
	"github.com/ematri2108281/wasatext/service/components/schema"
	"github.com/julienschmidt/httprouter"
)

// doLogin handles user login or registration
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	var req requests.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode login request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, "Username must be 3-20 characters", http.StatusBadRequest)
		return
	}

	user, err := rt.db.GetUserByName(req.Name)
	if errors.Is(err, schema.ErrUserDoesNotExist) {
		// Create new user
		newUser := &schema.User{Username: req.Name}

		userID, err := rt.db.CreateUser(newUser)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to create user")
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}
		newUser.ID = userID

		tokenString, err := createToken(userID)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to create token")
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"identifier": tokenString,
			"id":         newUser.ID,
			"username":   newUser.Username,
			"photo":      newUser.Photo,
		})
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get user by name")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tokenString, err := createToken(user.ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create token")
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"identifier": tokenString,
		"id":         user.ID,
		"username":   user.Username,
		"photo":      user.Photo,
	})
}
