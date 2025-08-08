package users

import (
	"gorm.io/gorm"
	"popcorntime-project/internal/models"
	"popcorntime-project/pkg/utils"
)

const RoleDefault = 1
const DefaultAvatar = 1

type Service interface {
	GetUser(id uint) (*models.User, error)
	Register(email, password, username string) (*models.User, error)
	Authenticate(email, password string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUsers() ([]models.User, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (u *service) GetUser(id uint) (*models.User, error) {
	var user models.User
	result := u.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *service) GetUsers() ([]models.User, error) {
	var users []models.User

	result := u.db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (u *service) Register(email, password, username string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Username:     username,
		RoleID:       RoleDefault,
		AvatarID:     DefaultAvatar,
	}

	if err := u.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *service) Authenticate(email, password string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}

func (u *service) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
