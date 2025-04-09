package users

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"popcorntime-project/internal/models"
	"regexp"
)

type CreateUserRequest struct {
	models.CreateUserRequest
}

func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email,
			validation.Required.Error("Email is required"),
			is.Email.Error("Invalid email format"),
			validation.Length(5, 255).Error("Email must be between 5-255 characters"),
		),
		validation.Field(&r.Password,
			validation.Required.Error("Password is required"),
			validation.Length(8, 100).Error("Password must be 8-100 characters"),
		),
		validation.Field(&r.Username,
			validation.Required.Error("Username is required"),
			validation.Length(3, 50).Error("Username must be 3-50 characters"),
			validation.Match(
				regexp.MustCompile(`^[a-zA-Z0-9_]+$`),
			).Error("Username can only contain letters, numbers and underscores"),
		),
	)
}
