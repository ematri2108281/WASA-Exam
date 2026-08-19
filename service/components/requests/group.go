package requests

// AddToGroupRequest is used when adding a user to a group
type AddToGroupRequest struct {
	ConversationID string `json:"conversationId"`
	UserID         string `json:"userId"`
}

// IsValid checks if the request has valid IDs
func (r *AddToGroupRequest) IsValid() bool {
	return len(r.ConversationID) > 0 && len(r.ConversationID) <= 36 &&
		len(r.UserID) > 0 && len(r.UserID) <= 36
}

// SetGroupNameRequest is used when changing a group name
type SetGroupNameRequest struct {
	Name string `json:"name"`
}

// IsValid checks if the new group name is valid
func (r *SetGroupNameRequest) IsValid() bool {
	return len(r.Name) >= 3 && len(r.Name) <= 50
}

// SetGroupPhotoRequest is used when updating group photo
type SetGroupPhotoRequest struct {
	Photo []byte `json:"photo"`
}

// IsValid checks if the photo data is present and within size limits
func (r *SetGroupPhotoRequest) IsValid() bool {
	return len(r.Photo) > 0 && len(r.Photo) <= (10*1024*1024)
}
