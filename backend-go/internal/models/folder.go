package models

import "time"

// Folder — папкаи "база".
type Folder struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:255" json:"name"`
	ParentID   uint      `gorm:"default:0;index" json:"parent_id"`
	AuthorID   uint      `gorm:"index" json:"author_id"`
	AuthorName string    `gorm:"size:255" json:"author_name"`
	CreatedAt  time.Time `json:"created_date"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Folder) TableName() string { return "folders" }

// FolderFile — файл дар папка.
type FolderFile struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FolderID     uint      `gorm:"index" json:"folder_id"`
	Filename     string    `gorm:"size:512" json:"filename"`
	OriginalName string    `gorm:"size:512" json:"original_name"`
	Filepath     string    `gorm:"size:1024" json:"filepath"`
	Filetype     string    `gorm:"size:64" json:"filetype"`
	Filesize     int64     `json:"filesize"`
	AuthorID     uint      `gorm:"index" json:"author_id"`
	AuthorName   string    `gorm:"size:255" json:"author_name"`
	CreatedAt    time.Time `json:"created_date"`
}

func (FolderFile) TableName() string { return "folder_files" }
