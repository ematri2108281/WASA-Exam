package requests

// AddReactionRequest is used when adding a reaction to a message
type AddReactionRequest struct {
	Type string `json:"type"` // like, heart, laugh, sad_face, angry_face
}

// IsValid checks if the reaction type is one of the allowed values
func (r *AddReactionRequest) IsValid() bool {
	validTypes := map[string]bool{
		"like":       true,
		"heart":      true,
		"laugh":      true,
		"sad_face":   true,
		"angry_face": true,
	}
	return validTypes[r.Type]
}
