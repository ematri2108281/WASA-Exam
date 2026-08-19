/*
Package database is the middleware between the app database and the code. All data (de)serialization
(save/load) from a persistent database are handled here. Database specific logic should never escape
this package.

To use this package you need to connect to it (using the database data source name from config), and
then initialize an instance of AppDatabase from the DB connection.

Example:

	db, err := sql.Open("sqlite3", cfg.DB.Filename)
	if err != nil { ... }
	defer db.Close()
	appdb, err := database.New(db)
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ematri2108281/wasatext/service/components/schema"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error

	// user related
	GetAllUsers() ([]schema.User, error)
	SearchUserByUsername(username string) ([]schema.User, error)
	GetUserByName(username string) (*schema.User, error)
	GetUserById(id string) (*schema.User, error)
	CreateUser(user *schema.User) (string, error)
	UpdateUsername(userID, newUsername string) error
	UpdateUserPhoto(userID string, photo []byte) error

	// conversation related
	GetMyConversations(userID string) ([]schema.Conversation, error)
	GetConversationByID(userID, conversationID string) (*schema.Conversation, error)
	SearchConversationByName(name string) ([]*schema.Conversation, error)
	CreateConversation(conversation *schema.Conversation) error
	GetLastMessageByConversationID(conversationID string) (*schema.Message, error)
	EnsureDirectConversation(userID, peerUserID string) (*schema.Conversation, error)
	GetConversationMembers(conversationID string) ([]schema.User, error)
	UpdateGroupName(conversationID, newName string) error
	UpdateGroupPhoto(conversationID string, photo []byte) error
	AddUserToGroup(conversationID, userID string) error
	LeaveGroup(conversationID, userID string) error

	// message related
	SendMessage(message *schema.Message, conversationID string) error
	GetMessagesByConversationID(conversationID string, currentUserID string) ([]*schema.Message, error)
	GetMessageByID(messageID string) (*schema.Message, error)
	ForwardMessage(message *schema.Message, userID string) error
	DeleteMessage(conversationID, messageID, userID string) error
	MarkMessageStatus(messageID, userID, status string) error

	// reaction related
	AddReactionToMessage(messageID, userID, reactionType string) error
	DeleteReactionFromMessage(messageID, userID, reactionType string) error
	GetReactionsForMessage(messageID string) ([]schema.Reaction, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist yet
	var tableName string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		creationQueries := []string{
			`CREATE TABLE users (
				id TEXT NOT NULL PRIMARY KEY,
				username TEXT NOT NULL UNIQUE,
				photo BLOB
			);`,
			`CREATE TABLE conversations (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				is_group BOOLEAN NOT NULL CHECK (is_group IN (0,1)),
				group_photo BLOB,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);`,
			`CREATE TABLE conversation_members (
				conversationId TEXT NOT NULL,
				userId TEXT NOT NULL,
				PRIMARY KEY (conversationId, userId),
				FOREIGN KEY (conversationId) REFERENCES conversations(id) ON DELETE CASCADE,
				FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
			);`,
			`CREATE TABLE messages (
				id TEXT NOT NULL PRIMARY KEY,
				conversationId TEXT NOT NULL,
				senderId TEXT NOT NULL,
				content_type TEXT NOT NULL,
				content_value TEXT NOT NULL,
				timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
				status TEXT NOT NULL,
				forwarded_from TEXT,
				FOREIGN KEY (conversationId) REFERENCES conversations(id) ON DELETE CASCADE,
				FOREIGN KEY (senderId) REFERENCES users(id) ON DELETE CASCADE
			);`,
			`CREATE TABLE reactions (
				messageId TEXT NOT NULL,
				userId TEXT NOT NULL,
				type TEXT NOT NULL,
				PRIMARY KEY (messageId, userId, type),
				FOREIGN KEY (messageId) REFERENCES messages(id) ON DELETE CASCADE,
				FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
			);`,
			`CREATE TABLE message_status (
				messageId TEXT NOT NULL,
				userId TEXT NOT NULL,
				deliveredAt DATETIME,
				readAt DATETIME,
				PRIMARY KEY (messageId, userId),
				FOREIGN KEY (messageId) REFERENCES messages(id) ON DELETE CASCADE,
				FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
			);`,
		}
		for _, stmt := range creationQueries {
			if _, err = db.Exec(stmt); err != nil {
				return nil, fmt.Errorf("error creating table: %w", err)
			}
		}
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
