package models

import "time"

// Request — заявкаи мизоҷ (legacy table; используется в /api/requests).
type Request struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PropertyType string    `gorm:"size:64" json:"property_type"`
	Address      string    `gorm:"size:255" json:"address"`
	Area         float64   `json:"area"`
	Rooms        int       `json:"rooms"`
	Windows      int       `json:"windows"`
	Floor        int       `json:"floor"`
	TotalFloors  int       `json:"total_floors"`
	Documents    string    `gorm:"size:128" json:"documents"`
	TotalPrice   float64   `json:"total_price"`
	PricePerM2   float64   `json:"price_per_m2"`
	Phone        string    `gorm:"size:64" json:"phone"`
	ClientName   string    `gorm:"size:255" json:"client_name"`
	Manager      string    `gorm:"size:255" json:"manager"`
	SMM          string    `gorm:"size:255" json:"smm"`
	Comment      string    `gorm:"type:text" json:"comment"`
	AuthorID     uint      `gorm:"index" json:"author_id"`
	ExecutorID   uint      `gorm:"index" json:"executor_id"`
	Files        JSONB     `gorm:"type:jsonb" json:"files"`
	Status       string    `gorm:"size:32;default:new" json:"status"`
	CreatedAt    time.Time `json:"created_date"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Request) TableName() string { return "requests" }

// RequestsBoard — доска для канбан заявок.
type RequestsBoard struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255" json:"title"`
	Color     string    `gorm:"size:32;default:#0f172a" json:"color"`
	IsPublic  bool      `gorm:"default:true" json:"is_public"`
	AuthorID  uint      `gorm:"index" json:"author_id"`
	CreatedAt time.Time `json:"created_date"`
	UpdatedAt time.Time `json:"updated_date"`
}

func (RequestsBoard) TableName() string { return "requests_boards" }

// RequestsColumn — колонка канбан заявок.
type RequestsColumn struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	BoardID    uint      `gorm:"index" json:"board_id"`
	Title      string    `gorm:"size:255" json:"title"`
	Color      string    `gorm:"size:32;default:#3b82f6" json:"color"`
	OrderIndex int       `gorm:"default:0" json:"order_index"`
	CreatedAt  time.Time `json:"created_date"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (RequestsColumn) TableName() string { return "requests_columns" }

// RequestsItem — карточка заявки в колонке.
type RequestsItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ColumnID     uint      `gorm:"index" json:"column_id"`
	BoardID      uint      `gorm:"index" json:"board_id"`
	PropertyType string    `gorm:"size:64" json:"property_type"`
	Address      string    `gorm:"size:255" json:"address"`
	Area         float64   `json:"area"`
	Rooms        int       `json:"rooms"`
	Windows      int       `json:"windows"`
	Floor        int       `json:"floor"`
	TotalFloors  int       `json:"total_floors"`
	TotalPrice   float64   `json:"total_price"`
	PricePerM2   float64   `json:"price_per_m2"`
	Phone        string    `gorm:"size:64" json:"phone"`
	ClientName   string    `gorm:"size:255" json:"client_name"`
	Comment      string    `gorm:"type:text" json:"comment"`
	AuthorID     uint      `gorm:"index" json:"author_id"`
	AuthorName   string    `gorm:"size:255" json:"author_name"`
	ExecutorID   uint      `gorm:"index" json:"executor_id"`
	ExecutorName string    `gorm:"size:255" json:"executor_name"`
	Files        JSONB     `gorm:"type:jsonb" json:"files"`
	OrderIndex   int       `gorm:"default:0" json:"order_index"`
	CreatedAt    time.Time `json:"created_date"`
	UpdatedAt    time.Time `json:"updated_date"`
}

func (RequestsItem) TableName() string { return "requests_items" }
