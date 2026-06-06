package models

import "time"

// ObjektProject — лоиҳаи бинокор (шахматка).
type ObjektProject struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Address     string    `gorm:"size:255" json:"address"`
	Developer   string    `gorm:"size:255" json:"developer"`
	AuthorID    uint      `gorm:"index" json:"author_id"`
	CreatedAt   time.Time `json:"created_date"`
	UpdatedAt   time.Time `json:"updated_date"`
}

func (ObjektProject) TableName() string { return "objekt_projects" }

// ObjektBlock — блок дар шахматка.
type ObjektBlock struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ProjectID         uint      `gorm:"index" json:"project_id"`
	Name              string    `gorm:"size:128" json:"name"`
	FloorFrom         int       `gorm:"default:1" json:"floor_from"`
	FloorTo           int       `gorm:"default:1" json:"floor_to"`
	OrderIndex        int       `gorm:"default:0" json:"order_index"`
	DefaultArea       float64   `gorm:"default:0" json:"default_area"`
	DefaultRooms      int       `gorm:"default:0" json:"default_rooms"`
	DefaultWindows    int       `gorm:"default:0" json:"default_windows"`
	DefaultPricePerM2 float64   `gorm:"default:0" json:"default_price_per_m2"`
	DefaultBalcony    string    `gorm:"size:32" json:"default_balcony"`
	DefaultBathroom   string    `gorm:"size:32" json:"default_bathroom"`
	DefaultPlanImage  string    `gorm:"size:1024" json:"default_plan_image"`
	CreatedAt         time.Time `json:"created_date"`
	UpdatedAt         time.Time `json:"updated_date"`
}

func (ObjektBlock) TableName() string { return "objekt_blocks" }

// ObjektApartment — хонаи алоҳида.
type ObjektApartment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	BlockID    uint      `gorm:"index" json:"block_id"`
	Floor      int       `json:"floor"`
	Area       float64   `gorm:"default:0" json:"area"`
	Rooms      int       `gorm:"default:0" json:"rooms"`
	Windows    int       `gorm:"default:0" json:"windows"`
	PricePerM2 float64   `gorm:"default:0" json:"price_per_m2"`
	TotalPrice float64   `gorm:"default:0" json:"total_price"`
	Status     string    `gorm:"size:32;default:free" json:"status"`
	Balcony    string    `gorm:"size:32" json:"balcony"`
	Bathroom   string    `gorm:"size:32" json:"bathroom"`
	PlanImage  string    `gorm:"size:1024" json:"plan_image"`
	CreatedAt  time.Time `json:"created_date"`
	UpdatedAt  time.Time `json:"updated_date"`
}

func (ObjektApartment) TableName() string { return "objekt_apartments" }
