package user

import (
	"ticket-api/internal/domain"
)

type DetailSerializer struct {
	domain.User
}

type CreateUser struct {
	domain.User
	ProfilePictureUrl *string `json:"profile_picture_url,omitempty"`
}
