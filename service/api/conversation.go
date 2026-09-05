package api

import (
	"encoding/json"
	"net/http"

	"github.com/ematri2108281/wasatext/service/api/reqcontext"
	"github.com/ematri2108281/wasatext/service/components/schema"
	"github.com/julienschmidt/httprouter"
)

// getMyConversations returns all conversations for the authenticated user
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversations, err := rt.db.GetMyConversations(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get conversations")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode conversations")
	}
}

// getConversation returns a single conversation with all messages
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Missing conversation ID", http.StatusBadRequest)
		return
	}

	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversation, err := rt.db.GetConversationByID(userID, conversationID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get conversation")
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	messages, err := rt.db.GetMessagesByConversationID(conversationID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get messages")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Mark messages as delivered for this user
	for _, msg := range messages {
		if msg.Sender.ID != userID {
			if err := rt.db.MarkMessageStatus(msg.ID, userID, "delivered"); err != nil {
				ctx.Logger.WithError(err).WithField("message_id", msg.ID).Error("Failed to mark message as delivered")
			}
		}
	}

	conversation.Messages = make([]schema.Message, len(messages))
	for i, msg := range messages {
		conversation.Messages[i] = *msg
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode conversation")
	}
}

// getConversationMembers returns all members of a conversation
func (rt *_router) getConversationMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Missing conversation ID", http.StatusBadRequest)
		return
	}

	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify user has access to this conversation
	if _, err := rt.db.GetConversationByID(userID, conversationID); err != nil {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	members, err := rt.db.GetConversationMembers(conversationID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get conversation members")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(members)
}

// createDirectConversation creates or retrieves a direct conversation between two users
func (rt *_router) createDirectConversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		PeerUserID string `json:"peerUserId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.PeerUserID == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	conv, err := rt.db.EnsureDirectConversation(userID, body.PeerUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to ensure direct conversation")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(conv)
}

// deleteConversation removes the authenticated user from a conversation
func (rt *_router) deleteConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Missing conversation ID", http.StatusBadRequest)
		return
	}

	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := rt.db.LeaveGroup(conversationID, userID); err != nil {
		ctx.Logger.WithError(err).Error("Failed to delete/leave conversation")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
