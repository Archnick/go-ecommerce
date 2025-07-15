package models

// UserPayload is the expected JSON body for registration.
type UserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
