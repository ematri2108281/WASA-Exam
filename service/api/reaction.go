package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ematri2108281/wasatext/service/api/reqcontext"
	"github.com/ematri2108281/wasatext/service/components/requests"
)

// commentMessage adds a reaction to a message
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	var req requests.AddReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode reaction request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, "Invalid reaction type", http.StatusBadRequest)
		return
	}

	if err := rt.db.AddReactionToMessage(messageID, userID, req.Type); err != nil {
		ctx.Logger.WithError(err).Error("Failed to add reaction")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// uncommentMessage removes a reaction from a message
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	reactionType := r.URL.Query().Get("type")
	if reactionType == "" {
		http.Error(w, "Missing reaction type", http.StatusBadRequest)
		return
	}

	validTypes := map[string]bool{
		"like": true, "heart": true, "laugh": true,
		"sad_face": true, "angry_face": true,
	}
	if !validTypes[reactionType] {
		http.Error(w, "Invalid reaction type", http.StatusBadRequest)
		return
	}

	if err := rt.db.DeleteReactionFromMessage(messageID, userID, reactionType); err != nil {
		ctx.Logger.WithError(err).Error("Failed to remove reaction")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
