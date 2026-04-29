package domain

import (
	"time"
)

type User struct {
	ID         string    `json:"id" gorm:"primaryKey;column:id"`
	Email      string    `json:"email" gorm:"unique;not null;column:email"`
	Name       string    `json:"name" gorm:"column:name"`
	Username   string    `json:"username" gorm:"unique;column:username"`
	AvatarURL  string    `json:"avatar_url" gorm:"column:avatar_url"`
	Bio        string    `json:"bio" gorm:"column:bio"`
	ExternalID string    `json:"external_id" gorm:"column:external_id"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
