package database

import (
	"fmt"

	"github.com/ematri2108281/wasatext/service/components/schema"
)

// AddReactionToMessage inserts a reaction (idempotent)
func (db *appdbimpl) AddReactionToMessage(messageID, userID, reactionType string) error {
	if messageID == "" || userID == "" || reactionType == "" {
		return fmt.Errorf("invalid reaction parameters")
	}

	_, err := db.c.Exec(`
		INSERT INTO reactions (messageId, userId, type)
		VALUES (?, ?, ?)
		ON CONFLICT(messageId, userId, type) DO NOTHING
	`, messageID, userID, reactionType)

	return err
}

// DeleteReactionFromMessage removes a specific reaction
func (db *appdbimpl) DeleteReactionFromMessage(messageID, userID, reactionType string) error {
	_, err := db.c.Exec(`
		DELETE FROM reactions
		WHERE messageId = ? AND userId = ? AND type = ?
	`, messageID, userID, reactionType)

	return err
}

// GetReactionsForMessage returns all reactions on a message
func (db *appdbimpl) GetReactionsForMessage(messageID string) ([]schema.Reaction, error) {
	rows, err := db.c.Query(`
		SELECT userId, type
		FROM reactions
		WHERE messageId = ?
	`, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []schema.Reaction
	for rows.Next() {
		var r schema.Reaction
		if err := rows.Scan(&r.AuthorID, &r.Type); err != nil {
			return nil, err
		}
		reactions = append(reactions, r)
	}

	return reactions, rows.Err()
}
