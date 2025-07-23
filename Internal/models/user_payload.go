package models

// UserPayload is the expected JSON body for registration.
type UserPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
