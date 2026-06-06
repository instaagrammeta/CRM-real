package models

import "time"

// Post — пост маркетинги (SMM).
type Post struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	UserName    string    `gorm:"size:255" json:"user_name"`
	Title       string    `gorm:"size:512" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"size:128" json:"category"`
	ContentType string    `gorm:"size:64" json:"content_type"`
	Project     string    `gorm:"size:128" json:"project"`
	MediaPath   string    `gorm:"size:1024" json:"media_path"`
	MediaType   string    `gorm:"size:64" json:"media_type"`
	Link        string    `gorm:"size:1024" json:"link"`
	PostDate    string    `gorm:"size:64" json:"post_date"`
	CreatedAt   time.Time `json:"created_date"`
	Likes       int       `gorm:"default:0" json:"likes"`
	Comments    int       `gorm:"default:0" json:"comments"`
	Shares      int       `gorm:"default:0" json:"shares"`
	Views       int       `gorm:"default:0" json:"views"`
	Reach       int       `gorm:"default:0" json:"reach"`
	IsPublished bool      `gorm:"default:true" json:"is_published"`
	PublishedAt string    `gorm:"size:64" json:"published_at"`
	UpdatedAt   time.Time `json:"updated_date"`
}

func (Post) TableName() string { return "posts" }
