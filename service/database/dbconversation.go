package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ematri2108281/wasatext/service/components/schema"
	"github.com/gofrs/uuid"
)

// GetMyConversations returns a summary list of conversations for a user
func (db *appdbimpl) GetMyConversations(userID string) ([]schema.Conversation, error) {
	query := `
		SELECT c.id, c.name, c.is_group, c.group_photo
		FROM conversations c
		JOIN conversation_members cm ON cm.conversationId = c.id
		WHERE cm.userId = ?`

	rows, err := db.c.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations: %w", err)
	}
	defer rows.Close()

	var conversations []schema.Conversation

	for rows.Next() {
		var conv schema.Conversation
		var groupPhoto []byte
		if err := rows.Scan(&conv.ID, &conv.Name, &conv.IsGroup, &groupPhoto); err != nil {
			return nil, err
		}

		if conv.IsGroup {
			conv.GroupPhoto = groupPhoto
		} else {
			var otherUser schema.User
			err := db.c.QueryRow(`
				SELECT u.id, u.username, u.photo
				FROM users u
				JOIN conversation_members cm ON cm.userId = u.id
				WHERE cm.conversationId = ? AND cm.userId != ?
			`, conv.ID, userID).Scan(&otherUser.ID, &otherUser.Username, &otherUser.Photo)
			if err != nil {
				return nil, fmt.Errorf("failed to get direct conversation info: %w", err)
			}
			conv.Participants = []schema.User{otherUser}
		}

		// Fetch last message preview
		var msgPreview string
		var ts time.Time
		err = db.c.QueryRow(`
			SELECT content_value, timestamp
			FROM messages
			WHERE conversationId = ?
			ORDER BY timestamp DESC LIMIT 1
		`, conv.ID).Scan(&msgPreview, &ts)
		if err == nil {
			conv.LastMessage = &schema.Message{
				Content:   schema.MessageContent{Value: msgPreview},
				Timestamp: ts,
			}
		}

		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

// GetConversationByID returns the full conversation with participants and messages
func (db *appdbimpl) GetConversationByID(userID, conversationID string) (*schema.Conversation, error) {
	var conv schema.Conversation
	var groupPhoto []byte
	err := db.c.QueryRow(`
		SELECT id, name, is_group, group_photo
		FROM conversations
		WHERE id = ?
	`, conversationID).Scan(&conv.ID, &conv.Name, &conv.IsGroup, &groupPhoto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation info: %w", err)
	}

	if conv.IsGroup {
		conv.GroupPhoto = groupPhoto
	}

	// Get participants
	memberRows, err := db.c.Query(`
		SELECT u.id, u.username, u.photo
		FROM users u
		JOIN conversation_members cm ON cm.userId = u.id
		WHERE cm.conversationId = ?
	`, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation members: %w", err)
	}
	defer memberRows.Close()

	for memberRows.Next() {
		var u schema.User
		if err := memberRows.Scan(&u.ID, &u.Username, &u.Photo); err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		conv.Participants = append(conv.Participants, u)
	}
	if err := memberRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating members: %w", err)
	}

	// Get all messages
	messageRows, err := db.c.Query(`
		SELECT m.id, m.senderId, m.content_value, m.timestamp
		FROM messages m
		WHERE m.conversationId = ?
		ORDER BY m.timestamp ASC
	`, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer messageRows.Close()

	for messageRows.Next() {
		var msg schema.Message
		var senderID, content string
		var ts time.Time
		if err := messageRows.Scan(&msg.ID, &senderID, &content, &ts); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		msg.Timestamp = ts
		msg.Content = schema.MessageContent{Value: content}

		var sender schema.User
		err := db.c.QueryRow(`SELECT id, username, photo FROM users WHERE id = ?`, senderID).Scan(&sender.ID, &sender.Username, &sender.Photo)
		if err != nil {
			return nil, fmt.Errorf("failed to get sender info: %w", err)
		}
		msg.Sender = sender

		conv.Messages = append(conv.Messages, msg)
	}
	if err := messageRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return &conv, nil
}

// CreateConversation inserts a new conversation and its members
func (db *appdbimpl) CreateConversation(conv *schema.Conversation) error {
	if conv.ID == "" {
		id, _ := uuid.NewV4()
		conv.ID = id.String()
	}

	_, err := db.c.Exec(`
		INSERT INTO conversations (id, name, is_group, group_photo, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
	`, conv.ID, conv.Name, conv.IsGroup, conv.GroupPhoto)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	for _, u := range conv.Participants {
		_, err := db.c.Exec(`INSERT INTO conversation_members (conversationId, userId) VALUES (?, ?)`, conv.ID, u.ID)
		if err != nil {
			return fmt.Errorf("failed to add member: %w", err)
		}
	}

	return nil
}

// SearchConversationByName finds conversations containing the given name
func (db *appdbimpl) SearchConversationByName(name string) ([]*schema.Conversation, error) {
	rows, err := db.c.Query(`SELECT id, name, is_group, group_photo FROM conversations WHERE name LIKE '%' || ? || '%'`, name)
	if err != nil {
		return nil, fmt.Errorf("failed to search conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*schema.Conversation
	for rows.Next() {
		var conv schema.Conversation
		if err := rows.Scan(&conv.ID, &conv.Name, &conv.IsGroup, &conv.GroupPhoto); err != nil {
			return nil, err
		}
		conversations = append(conversations, &conv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating conversations: %w", err)
	}

	return conversations, nil
}

// GetLastMessageByConversationID fetches the most recent message of a conversation
func (db *appdbimpl) GetLastMessageByConversationID(conversationID string) (*schema.Message, error) {
	var msg schema.Message
	var senderID, content string
	var ts time.Time

	err := db.c.QueryRow(`
		SELECT id, senderId, content_value, timestamp
		FROM messages
		WHERE conversationId = ?
		ORDER BY timestamp DESC LIMIT 1
	`, conversationID).Scan(&msg.ID, &senderID, &content, &ts)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // no messages yet
		}
		return nil, fmt.Errorf("failed to fetch last message: %w", err)
	}

	msg.Timestamp = ts
	msg.Content = schema.MessageContent{Value: content}

	var sender schema.User
	err = db.c.QueryRow(`SELECT id, username, photo FROM users WHERE id = ?`, senderID).Scan(&sender.ID, &sender.Username, &sender.Photo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sender info: %w", err)
	}
	msg.Sender = sender

	return &msg, nil
}

// EnsureDirectConversation finds or creates a 1-on-1 conversation between two users
func (db *appdbimpl) EnsureDirectConversation(userID, peerUserID string) (*schema.Conversation, error) {
	var conversationID string
	err := db.c.QueryRow(`
        SELECT c.id
        FROM conversations c
        JOIN conversation_members cm1 ON cm1.conversationId = c.id AND cm1.userId = ?
        JOIN conversation_members cm2 ON cm2.conversationId = c.id AND cm2.userId = ?
        WHERE c.is_group = 0
        LIMIT 1
    `, userID, peerUserID).Scan(&conversationID)
	if err == nil {
		return db.GetConversationByID(userID, conversationID)
	}

	convUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	convID := convUUID.String()

	_, err = db.c.Exec(`INSERT INTO conversations (id, name, is_group) VALUES (?, ?, 0)`, convID, "")
	if err != nil {
		return nil, err
	}

	_, err = db.c.Exec(`
        INSERT INTO conversation_members (conversationId, userId)
        VALUES (?, ?), (?, ?)
    `, convID, userID, convID, peerUserID)
	if err != nil {
		return nil, err
	}

	return db.GetConversationByID(userID, convID)
}

// GetConversationMembers returns all members of a conversation
func (db *appdbimpl) GetConversationMembers(conversationID string) ([]schema.User, error) {
	rows, err := db.c.Query(`
		SELECT u.id, u.username, u.photo
		FROM users u
		JOIN conversation_members cm ON cm.userId = u.id
		WHERE cm.conversationId = ?
		ORDER BY u.username ASC
	`, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversation members: %w", err)
	}
	defer rows.Close()

	var members []schema.User
	for rows.Next() {
		var u schema.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Photo); err != nil {
			return nil, fmt.Errorf("failed to scan member row: %w", err)
		}
		members = append(members, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating members: %w", err)
	}

	return members, nil
}

// UpdateGroupName renames a group conversation
func (db *appdbimpl) UpdateGroupName(conversationID, newName string) error {
	_, err := db.c.Exec(`UPDATE conversations SET name = ? WHERE id = ? AND is_group = 1`, newName, conversationID)
	if err != nil {
		return fmt.Errorf("failed to update group name: %w", err)
	}
	return nil
}

// UpdateGroupPhoto updates a group's profile photo
func (db *appdbimpl) UpdateGroupPhoto(conversationID string, photo []byte) error {
	_, err := db.c.Exec(`UPDATE conversations SET group_photo = ? WHERE id = ? AND is_group = 1`, photo, conversationID)
	if err != nil {
		return fmt.Errorf("failed to update group photo: %w", err)
	}
	return nil
}

// AddUserToGroup inserts a new participant into a group
func (db *appdbimpl) AddUserToGroup(conversationID, userID string) error {
	_, err := db.c.Exec(`
		INSERT INTO conversation_members (conversationId, userId)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING
	`, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}
	return nil
}

// LeaveGroup removes a user from a group
func (db *appdbimpl) LeaveGroup(conversationID, userID string) error {
	_, err := db.c.Exec(`
		DELETE FROM conversation_members
		WHERE conversationId = ? AND userId = ?
	`, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}
	return nil
}
