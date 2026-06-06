package models

import "time"

// Lid — лиди мизоҷ (шакли соддаи lid).
type Lid struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ClientName string    `gorm:"size:255" json:"client_name"`
	Phone      string    `gorm:"size:64" json:"phone"`
	Topic      string    `gorm:"size:255" json:"topic"`
	Comment    string    `gorm:"type:text" json:"comment"`
	Source     string    `gorm:"size:128" json:"source"`
	Mortgage   bool      `gorm:"default:false" json:"mortgage"`
	Box        bool      `gorm:"default:false" json:"box"`
	AuthorID   uint      `gorm:"index" json:"author_id"`
	CreatedAt  time.Time `json:"created_date"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Lid) TableName() string { return "lids" }

// KanbanBoard — доскаи канбан барои лидҳо.
type KanbanBoard struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"size:255" json:"title"`
	Color      string    `gorm:"size:32;default:#0079bf" json:"color"`
	IsPublic   bool      `gorm:"default:true" json:"is_public"`
	AuthorID   uint      `gorm:"index" json:"author_id"`
	IsArchived bool      `gorm:"default:false" json:"is_archived"`
	CreatedAt  time.Time `json:"created_date"`
	UpdatedAt  time.Time `json:"updated_date"`
}

func (KanbanBoard) TableName() string { return "kanban_boards" }

// KanbanBoardMember — иштирокчиёни доскаи приватӣ.
type KanbanBoardMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BoardID   uint      `gorm:"index" json:"board_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	CreatedAt time.Time `json:"added_date"`
}

func (KanbanBoardMember) TableName() string { return "kanban_board_members" }

// KanbanColumn — колонкаи канбан.
type KanbanColumn struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	BoardID    uint      `gorm:"index" json:"board_id"`
	Title      string    `gorm:"size:255" json:"title"`
	Color      string    `gorm:"size:32;default:#0079bf" json:"color"`
	OrderIndex int       `gorm:"default:0" json:"order_index"`
	IsArchived bool      `gorm:"default:false" json:"is_archived"`
	CreatedAt  time.Time `json:"created_date"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (KanbanColumn) TableName() string { return "kanban_columns" }

// KanbanLead — карточкаи лид дар колонка.
type KanbanLead struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ColumnID   uint      `gorm:"index" json:"column_id"`
	BoardID    uint      `gorm:"index" json:"board_id"`
	ClientName string    `gorm:"size:255" json:"client_name"`
	Phone      string    `gorm:"size:64" json:"phone"`
	Topic      string    `gorm:"size:255" json:"topic"`
	Comment    string    `gorm:"type:text" json:"comment"`
	Source     string    `gorm:"size:128" json:"source"`
	Mortgage   bool      `gorm:"default:false" json:"mortgage"`
	Box        bool      `gorm:"default:false" json:"box"`
	AuthorID   uint      `gorm:"index" json:"author_id"`
	AuthorName string    `gorm:"size:255" json:"author_name"`
	OrderIndex int       `gorm:"default:0" json:"order_index"`
	CreatedAt  time.Time `json:"created_date"`
	UpdatedAt  time.Time `json:"updated_date"`
}

func (KanbanLead) TableName() string { return "kanban_leads" }

// KanbanLeadInteraction — таърихи робита бо лид.
type KanbanLeadInteraction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LeadID      uint      `gorm:"index" json:"lead_id"`
	Phone       string    `gorm:"size:64" json:"phone"`
	Topic       string    `gorm:"size:255" json:"topic"`
	Source      string    `gorm:"size:128" json:"source"`
	ContactType string    `gorm:"size:64" json:"contact_type"`
	Comment     string    `gorm:"type:text" json:"comment"`
	Mortgage    bool      `gorm:"default:false" json:"mortgage"`
	Box         bool      `gorm:"default:false" json:"box"`
	CreatedAt   time.Time `json:"created_date"`
}

func (KanbanLeadInteraction) TableName() string { return "kanban_lead_interactions" }
