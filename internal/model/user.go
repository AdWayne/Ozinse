package model

import "time"

type User struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	UserAvatarURL *string   `json:"user_avatar_url"`
	FullName      string    `json:"full_name"`
	Phone         *string   `json:"phone"`
	BirthDate     *string   `json:"birth_date"`
	RoleID        *int      `json:"role_id"`
	RoleName      string    `json:"role_name"`
	CreatedAt     time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	RepeatPassword string `json:"repeat_password" binding:"required"`
	FullName       string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateProfileRequest struct {
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	BirthDate string `json:"birth_date"`
}

type ChangePasswordRequest struct {
	OldPassword       string `json:"old_password" binding:"required"`
	NewPassword       string `json:"new_password" binding:"required,min=6"`
}

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Permissions string `json:"permissions"`
}