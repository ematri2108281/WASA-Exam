package requests

// MessageContent represents the content of a message
type MessageContent struct {
	Type  string `json:"type"`  // "text" or "photo"
	Value string `json:"value"` // text content or base64 encoded photo
}

// IsValid checks if the message content is valid
func (c *MessageContent) IsValid() bool {
	if c.Type != "text" && c.Type != "photo" && c.Type != "image" {
		return false
	}
	if c.Type == "text" {
		return len(c.Value) >= 1 && len(c.Value) <= 1000
	}
	return len(c.Value) >= 1 && len(c.Value) <= (10*1024*1024)
}

// SendMessageRequest is used when sending a message to a conversation
type SendMessageRequest struct {
	Content MessageContent `json:"content"`
}

// IsValid checks if the message request is valid
func (r *SendMessageRequest) IsValid() bool {
	return r.Content.IsValid()
}

// ForwardMessageRequest is used when forwarding a message to another conversation
type ForwardMessageRequest struct {
	TargetConversationID string `json:"targetConversationId"`
}

// IsValid checks if the forward request has a valid target conversation ID
func (r *ForwardMessageRequest) IsValid() bool {
	return len(r.TargetConversationID) > 0 && len(r.TargetConversationID) <= 36
}
