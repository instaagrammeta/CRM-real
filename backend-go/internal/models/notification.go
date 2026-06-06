package models

import "time"

// Notification — оғоҳномаи дохилӣ барои корбар.
type Notification struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Title      string    `gorm:"size:255" json:"title"`
	Body       string    `gorm:"type:text" json:"body"`
	Type       string    `gorm:"size:64;index" json:"type"`
	EntityType string    `gorm:"size:64" json:"entity_type"`
	EntityID   uint      `json:"entity_id"`
	Link       string    `gorm:"size:512" json:"link"`
	IsRead     bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Notification) TableName() string { return "notifications" }

// TelegramSubscriber — пайвандсози Telegram chat-а бо корбар.
type TelegramSubscriber struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ChatID    int64     `gorm:"uniqueIndex" json:"chat_id"`
	Username  string    `gorm:"size:128" json:"username"`
	FirstName string    `gorm:"size:128" json:"first_name"`
	LastName  string    `gorm:"size:128" json:"last_name"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	LinkToken string    `gorm:"size:64;index" json:"link_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TelegramSubscriber) TableName() string { return "telegram_subscribers" }
