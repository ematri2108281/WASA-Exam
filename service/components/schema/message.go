package schema

import "time"

type Message struct {
	ID             string         `json:"id"`
	ConversationID string         `json:"conversationId"`
	SenderID       string         `json:"senderId"`
	Sender         User           `json:"sender"`
	Content        MessageContent `json:"content"`
	Timestamp      time.Time      `json:"timestamp"`
	MessageStatus  string         `json:"message_status"`
	Reactions      []Reaction     `json:"reactions,omitempty"`
	ForwardedFrom  string         `json:"forwarded_from,omitempty"`
}

type MessageContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Reaction struct {
	Type     string `json:"type"`
	AuthorID string `json:"authorId"`
}
