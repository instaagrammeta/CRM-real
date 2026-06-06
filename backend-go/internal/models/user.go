package models

import "time"

// User — корманди система (admin / employee).
type User struct {
	ID             uint        `gorm:"primaryKey" json:"id"`
	FullName       string      `gorm:"size:255" json:"full_name"`
	Age            int         `json:"age"`
	PersonalPhones StringSlice `gorm:"type:jsonb" json:"personal_phones"`
	WorkPhones     StringSlice `gorm:"type:jsonb" json:"work_phones"`
	Login          string      `gorm:"size:128;uniqueIndex" json:"login"`
	Password       string      `gorm:"size:255" json:"-"`
	Photo          string      `gorm:"size:512" json:"photo"`
	Category       string      `gorm:"size:128" json:"category"`
	Role           string      `gorm:"size:32;default:employee" json:"role"`
	TelegramChatID int64       `gorm:"index" json:"telegram_chat_id"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

func (User) TableName() string { return "users" }
