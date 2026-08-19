package schema

// User represents an application user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Photo    []byte `json:"photo,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	User  User
	Token string `json:"token"`
}

type UsernameUpdateRequest struct {
	Username string `json:"username"`
}

type UsernameUpdateResponse = User

type ProfilePhotoUpdateRequest struct {
	Photo []byte `json:"photo"`
}

type ProfilePhotoUpdateResponse struct {
	Photo []byte `json:"photo"`
}
