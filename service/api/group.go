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

// createGroup creates a new group conversation
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		GroupName  string   `json:"groupName"`
		Members    []string `json:"members"`
		GroupPhoto []byte   `json:"groupPhoto,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if body.GroupName == "" || len(body.Members) == 0 {
		http.Error(w, "Group name and members are required", http.StatusBadRequest)
		return
	}

	// Ensure creator is included in members
	foundCreator := false
	for _, m := range body.Members {
		if m == userID {
			foundCreator = true
			break
		}
	}
	if !foundCreator {
		body.Members = append(body.Members, userID)
	}

	// Convert member IDs to User objects
	var participants []schema.User
	for _, memberID := range body.Members {
		user, err := rt.db.GetUserById(memberID)
		if err != nil {
			ctx.Logger.WithError(err).WithField("user_id", memberID).Error("Failed to get user")
			http.Error(w, "Invalid member ID: "+memberID, http.StatusBadRequest)
			return
		}
		participants = append(participants, *user)
	}

	conversation := &schema.Conversation{
		Name:         body.GroupName,
		IsGroup:      true,
		GroupPhoto:   body.GroupPhoto,
		Participants: participants,
	}

	if err := rt.db.CreateConversation(conversation); err != nil {
		ctx.Logger.WithError(err).Error("Failed to create group in database")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	conv, err := rt.db.GetConversationByID(userID, conversation.ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Group created but failed to load conversation")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(conv); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode group response")
	}
}

// addToGroup adds a user to a group
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupID := ps.ByName("groupId")
	if groupID == "" {
		http.Error(w, "Missing group ID", http.StatusBadRequest)
		return
	}

	if _, err := rt.getAuthenticatedUserID(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.AddToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if err := rt.db.AddUserToGroup(groupID, req.UserID); err != nil {
		if errors.Is(err, schema.ErrUserAlreadyInGroup) {
			http.Error(w, "User already in group", http.StatusBadRequest)
			return
		}
		ctx.Logger.WithError(err).Error("Failed to add user to group")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// leaveGroup removes a user from a group
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupID := ps.ByName("groupId")
	if groupID == "" {
		http.Error(w, "Missing group ID", http.StatusBadRequest)
		return
	}

	if _, err := rt.getAuthenticatedUserID(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetUserID := r.URL.Query().Get("userId")
	if targetUserID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	if err := rt.db.LeaveGroup(groupID, targetUserID); err != nil {
		if errors.Is(err, schema.ErrUserNotInGroup) {
			http.Error(w, "User not in group", http.StatusBadRequest)
			return
		}
		ctx.Logger.WithError(err).Error("Failed to remove user from group")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// setGroupName updates a group's name
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupID := ps.ByName("groupId")
	if groupID == "" {
		http.Error(w, "Missing group ID", http.StatusBadRequest)
		return
	}

	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.SetGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, "Invalid group name", http.StatusBadRequest)
		return
	}

	if err := rt.db.UpdateGroupName(groupID, req.Name); err != nil {
		if errors.Is(err, schema.ErrGroupNotFound) {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("Failed to update group name")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	conv, err := rt.db.GetConversationByID(userID, groupID)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(conv)
}

// setGroupPhoto updates a group's photo
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupID := ps.ByName("groupId")
	if groupID == "" {
		http.Error(w, "Missing group ID", http.StatusBadRequest)
		return
	}

	if _, err := rt.getAuthenticatedUserID(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.SetGroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, "Invalid photo data", http.StatusBadRequest)
		return
	}

	fileType := http.DetectContentType(req.Photo)
	if fileType != "image/jpeg" && fileType != "image/png" {
		http.Error(w, "Invalid file type. Only JPEG and PNG are supported.", http.StatusUnsupportedMediaType)
		return
	}

	if err := rt.db.UpdateGroupPhoto(groupID, req.Photo); err != nil {
		ctx.Logger.WithError(err).Error("Failed to update group photo")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Photo updated successfully",
		"groupPhoto": req.Photo,
	}); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode photo update response")
	}
}
