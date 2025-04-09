package models

import (
	"time"
)

type Role struct {
	ID          int      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string   `json:"name" gorm:"size:50;unique;not null"`
	Permissions []string `json:"permissions" gorm:"type:text[];default:'{}'"`
}

type User struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Email        string    `json:"email" gorm:"size:255;unique;not null"`
	PasswordHash string    `json:"-" gorm:"type:text;not null"` // Пароль исключаем из JSON
	Username     string    `json:"username" gorm:"type:text;unique;not null"`
	AvatarURL    string    `json:"avatar_url,omitempty" gorm:"type:text"`
	RoleID       int       `json:"role_id" gorm:"not null"`
	Role         Role      `json:"role" gorm:"foreignKey:RoleID;constraint:OnDelete:RESTRICT"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// UserResponse - структура для безопасного возврата данных пользователя
type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request body"`
	Details string `json:"details,omitempty" example:"invalid character '}' looking for beginning of value"`
}

type ValidationErrorResponse struct {
	Error  string            `json:"error" example:"Validation failed"`
	Errors map[string]string `json:"errors" example:"email:Invalid email format,password:Password must be 8-100 characters"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Username string `json:"username" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// ToResponse преобразует User в UserResponse (без пароля)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
