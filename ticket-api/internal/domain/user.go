package domain

import (
	"encoding/json"

	"ticket-api/internal/types"
	"ticket-api/pkg/gcp"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StripeCustomer struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Description string `json:"description"`
}

type User struct {
	Id             types.BinaryUUID `json:"id"`
	FirebaseId     *string          `json:"firebase_id"`
	Email          *string          `json:"email" form:"email"`
	NickName       *string          `json:"nick_name" form:"nick_name"`
	Username       *string          `json:"username" form:"username"`
	Bio            *string          `json:"bio" form:"bio"`
	Genre          json.RawMessage  `json:"genre" form:"genre"`
	Gender         *string          `json:"gender" form:"gender"`
	Age            *string          `json:"age" form:"age"`
	StripeID       *string          `json:"stripe_id" form:"stripe_id"`
	ProfilePicture *gcp.SignedURL   `json:"profile_picture"`
	Location       *string          `json:"location"`
	TempEmail      *bool            `json:"temp_email"`
	BaseDate
}

type ShortUserDetail struct {
	LessUserDetail
	FirebaseId string `json:"firebase_id"`
}

type LessUserDetail struct {
	Id             types.BinaryUUID `json:"id"`
	ProfilePicture *gcp.SignedURL   `json:"profile_picture"`
	NickName       string           `json:"nick_name"`
	Username       string           `json:"username"`
}

type UserService interface {
	CreateUser(*User) error

	GetAllUsers() (map[string]interface{}, error)

	GetUserByFirebaseId(string) (User, error)

	GetUserByStripeId(string) (User, error)

	SavePicture(User, types.BinaryUUID) error

	UpdateUser(userID types.BinaryUUID, user *map[string]interface{}) error

	GetUserById(userID types.BinaryUUID) (User, error)

	GetUserWithLessDetail(userID []string) ([]*ShortUserDetail, error)

	DeleteUser(uuid types.BinaryUUID) error

	SearchUsers(query string) (map[string]interface{}, error)

	SetPaginationScope(scope func(*gorm.DB) *gorm.DB) UserService
}

// BeforeCreate run this before creating
func (t *User) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	t.Id = types.BinaryUUID(id)
	return err
}
