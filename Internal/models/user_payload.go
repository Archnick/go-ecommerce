package models

// UserPayload is the expected JSON body for registration.
type UserPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserPayload struct {
	FirstName       string `json:"first_name" binding:"omitempty,min=2"`
	LastName        string `json:"last_name" binding:"omitempty,min=2"`
	ProfileImageURL string `json:"profile_img_url" binding:"omitempty,url"`
	Role            Role   `json:"role" binding:"omitempty,oneof=admin shop customer"`
}
