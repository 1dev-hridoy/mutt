package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null;type:varchar(70)"`
	Email    string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,email,max=255"`
	Password string `json:"-" gorm:"not null" validate:"required"`
	Phone    string `json:"phone" gorm:"type:varchar(20);uniqueIndex;not null" validate:"required,max=20"`
}

type SignupRequest struct {
	Username string `json:"username" validate:"required,min=3,max=70"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Phone    string `json:"phone" validate:"required,max=20"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
	}
}
