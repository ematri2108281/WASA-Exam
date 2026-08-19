package requests

import "regexp"

// LoginRequest is used for the login/register endpoint
type LoginRequest struct {
	Name string `json:"name"`
}

// IsValid checks if the username meets format requirements
func (r *LoginRequest) IsValid() bool {
	if len(r.Name) < 3 || len(r.Name) > 20 {
		return false
	}
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_ ]+$`, r.Name)
	return match
}
