package models

import "time"

// Message — паёми чат.
type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	UserName  string    `gorm:"size:255" json:"user_name"`
	Message   string    `gorm:"type:text" json:"message"`
	FilePath  string    `gorm:"size:1024" json:"file_path"`
	FileName  string    `gorm:"size:512" json:"file_name"`
	FileType  string    `gorm:"size:64" json:"file_type"`
	CreatedAt time.Time `json:"created_date"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Message) TableName() string { return "messages" }
