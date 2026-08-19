package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ematri2108281/wasatext/service/api/reqcontext"
	"github.com/ematri2108281/wasatext/service/components/requests"
	"github.com/ematri2108281/wasatext/service/components/schema"
)

// sendMessage sends a new message to a conversation
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	var req requests.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode message")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, "Invalid message content", http.StatusBadRequest)
		return
	}

	message := &schema.Message{
		Sender: schema.User{ID: userID},
		Content: schema.MessageContent{
			Type:  req.Content.Type,
			Value: req.Content.Value,
		},
		Timestamp:     time.Now().UTC(),
		MessageStatus: "sent",
	}

	if err := rt.db.SendMessage(message, conversationID); err != nil {
		ctx.Logger.WithError(err).Error("Failed to send message")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	stored, err := rt.db.GetMessageByID(message.ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get stored message")
		w.WriteHeader(http.StatusCreated)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(stored)
}

// forwardMessage forwards an existing message to another conversation
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "Missing message ID", http.StatusBadRequest)
		return
	}

	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.ForwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode forward message request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, "Invalid target conversation ID", http.StatusBadRequest)
		return
	}

	originalMessage, err := rt.db.GetMessageByID(messageID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to fetch original message")
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	forwardedMessage := &schema.Message{
		Sender: schema.User{ID: userID},
		Content: schema.MessageContent{
			Type:  originalMessage.Content.Type,
			Value: originalMessage.Content.Value,
		},
		Timestamp:      time.Now().UTC(),
		MessageStatus:  "sent",
		ForwardedFrom:  originalMessage.ID,
		ConversationID: req.TargetConversationID,
	}

	if err := rt.db.SendMessage(forwardedMessage, req.TargetConversationID); err != nil {
		ctx.Logger.WithError(err).Error("Failed to send forwarded message")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	stored, err := rt.db.GetMessageByID(forwardedMessage.ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get stored forwarded message")
		w.WriteHeader(http.StatusCreated)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(stored)
}

// deleteMessage deletes a message from a conversation
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationID := ps.ByName("conversationId")
	messageID := ps.ByName("messageId")

	if conversationID == "" || messageID == "" {
		http.Error(w, "Missing conversation or message ID", http.StatusBadRequest)
		return
	}

	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := rt.db.DeleteMessage(conversationID, messageID, userID); err != nil {
		ctx.Logger.WithError(err).
			WithField("conversation_id", conversationID).
			WithField("message_id", messageID).
			WithField("user_id", userID).
			Error("Failed to delete message")
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// setMessageStatus updates the delivery/read status of a message
func (rt *_router) setMessageStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationID := ps.ByName("conversationId")
	messageID := ps.ByName("messageId")

	if conversationID == "" || messageID == "" {
		http.Error(w, "Missing conversation or message ID", http.StatusBadRequest)
		return
	}

	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode message status request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status != "delivered" && req.Status != "read" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	if err := rt.db.MarkMessageStatus(messageID, userID, req.Status); err != nil {
		ctx.Logger.WithError(err).Error("Failed to update message status")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
