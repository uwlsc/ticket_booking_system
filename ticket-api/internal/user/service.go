package user

import (
	"ticket-api/internal/domain"
	"ticket-api/internal/types"
	"ticket-api/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	logger     pkg.Logger
	repository Repository
}

func NewService(logger pkg.Logger, repository Repository) domain.UserService {
	return &Service{logger: logger, repository: repository}
}

// SetPaginationScope .
func (s Service) SetPaginationScope(scope func(*gorm.DB) *gorm.DB) domain.UserService {
	s.repository.DB = s.repository.Scopes(scope)
	return s
}

func (s Service) CreateUser(user *domain.User) error {
	return s.repository.Create(&user).Error
}

func (s Service) GetUserByFirebaseId(firebaseId string) (user domain.User, err error) {
	return user, s.repository.Where("firebase_id=? AND deleted_at IS NULL", firebaseId).First(&user).Error
}

func (s Service) GetUserByStripeId(StripeId string) (user domain.User, err error) {
	return user, s.repository.Where("stripe_id=?", StripeId).First(&user).Error
}

func (s Service) SavePicture(user domain.User, userID types.BinaryUUID) error {
	return s.repository.Model(&domain.User{}).Where("id=?", userID).Updates(user).Error
}

// UpdateUser updates the user
func (s Service) UpdateUser(userID types.BinaryUUID, user *map[string]interface{}) error {
	return s.repository.Model(&domain.User{}).Where("id=?", userID).Updates(user).Error
}

// GetUserById Get user by ID
func (s Service) GetUserById(userID types.BinaryUUID) (user domain.User, err error) {
	return user, s.repository.Where("id=? AND deleted_at IS NULL", userID).First(&user).Error
}

// GetUserWithLessDetail Get user by ID
func (s Service) GetUserWithLessDetail(userIDs []string) (users []*domain.ShortUserDetail, err error) {
	err = s.repository.Table("users").Where("firebase_id IN ?", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetAllUsers get all users
func (s Service) GetAllUsers() (map[string]interface{}, error) {
	var users []domain.User
	var count int64
	err := s.repository.
		Where("deleted_at IS NULL AND temp_email = false").
		Order("created_at desc").
		Find(&users).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error
	if err != nil {
		return nil, err
	}
	return gin.H{"data": users, "count": count}, nil
}

// DeleteUser deletes the user
func (s Service) DeleteUser(uuid types.BinaryUUID) error {
	return s.repository.Where("id = ?", uuid).Delete(&domain.User{}).Error
}

// SearchUsers Search Users based on nickname, username, email
func (s Service) SearchUsers(query string) (map[string]interface{}, error) {
	var users []domain.User
	var count int64
	err := s.repository.
		Where(
			"lower(nick_name) LIKE ? OR lower(username) LIKE ? OR email LIKE ?", "%"+query+"%",
			"%"+query+"%", "%"+query+"%").
		Where("deleted_at IS NULL AND temp_email = false").
		Find(&users).Order("created_at desc").Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return nil, err
	}

	return gin.H{"data": users, "count": count}, nil
}
