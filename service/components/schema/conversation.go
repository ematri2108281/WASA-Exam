package schema

import "time"

type Conversation struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GroupPhoto []byte `json:"groupPhoto,omitempty"`
	IsGroup    bool   `json:"isGroup"`

	Participants []User    `json:"participants"`
	CreatedAt    time.Time `json:"createdAt"`
	LastMessage  *Message  `json:"lastMessage,omitempty"`

	Messages []Message `json:"messages,omitempty"`
}
