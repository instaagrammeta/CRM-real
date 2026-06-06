package models

import "time"

// Task — задача корманд.
type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	AuthorID    uint      `gorm:"index" json:"author_id"`
	ExecutorID  uint      `gorm:"index" json:"executor_id"`
	Photo       string    `gorm:"size:512" json:"photo"`
	Status      string    `gorm:"size:32;default:new" json:"status"`
	CreatedAt   time.Time `json:"created_date"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Task) TableName() string { return "tasks" }
