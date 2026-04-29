package domain

import (
	"time"
)

type Event struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:id"`
	CreatorID    string    `json:"creator_id" gorm:"not null;index;column:creator_id"`
	Title        string    `json:"title" gorm:"not null;column:title"`
	Description  string    `json:"description" gorm:"column:description"`
	StartTime    time.Time `json:"start_time" gorm:"not null;column:start_time"`
	EndTime      time.Time `json:"end_time" gorm:"column:end_time"`
	Location     string    `json:"location" gorm:"column:location"`
	DisplayID    string    `json:"display_id" gorm:"unique;not null;index;column:display_id"`
	Status       string    `json:"status" gorm:"not null;default:draft;column:status"` // "draft", "published"
	PopupTiming  string    `json:"popup_timing" gorm:"column:popup_timing"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}
