package domain

import (
	"time"
)

type User struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"unique;not null"`
	ExternalID string    `json:"external_id" gorm:"column:external_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
