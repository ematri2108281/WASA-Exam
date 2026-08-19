package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ematri2108281/wasatext/service/api/reqcontext"
	"github.com/ematri2108281/wasatext/service/components/requests"
	"github.com/ematri2108281/wasatext/service/components/schema"
)

// searchBy handles searching for users and conversations by query parameter
func (rt *_router) searchBy(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userQuery := r.URL.Query().Get("user")
	conversationQuery := r.URL.Query().Get("conversation")

	if userQuery == "" && conversationQuery == "" {
		http.Error(w, "Missing query parameters", http.StatusBadRequest)
		return
	}

	var users []schema.User
	var conversations []*schema.Conversation
	var err error

	if userQuery != "" {
		users, err = rt.db.SearchUserByUsername(userQuery)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to search users by username")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if conversationQuery != "" {
		conversations, err = rt.db.SearchConversationByName(conversationQuery)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to search conversations by name")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := struct {
		Users         []schema.User          `json:"users,omitempty"`
		Conversations []*schema.Conversation `json:"conversations,omitempty"`
	}{
		Users:         users,
		Conversations: conversations,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode response")
	}
}

// getAllUsers returns all registered users
func (rt *_router) getAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	users, err := rt.db.GetAllUsers()
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to get users")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"users": users}); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode response")
	}
}

// setMyUserName updates the authenticated user's username
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.UsernameUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, "Username must be 3-20 characters", http.StatusBadRequest)
		return
	}

	if err := rt.db.UpdateUsername(userID, req.Name); err != nil {
		if errors.Is(err, schema.ErrUserDoesNotExist) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, schema.ErrUsernameTaken) {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		ctx.Logger.WithError(err).Error("Failed to update username")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user, err := rt.db.GetUserById(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get updated user")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

// setMyPhoto updates the authenticated user's profile photo
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.ProfilePhotoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, "Invalid photo data", http.StatusBadRequest)
		return
	}

	if err := rt.db.UpdateUserPhoto(userID, req.Photo); err != nil {
		if errors.Is(err, schema.ErrUserDoesNotExist) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("Failed to update user photo")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
