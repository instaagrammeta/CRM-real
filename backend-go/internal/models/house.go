package models

import "time"

// House — обекти амлок (oddiy roʻyxat).
type House struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Title               string    `gorm:"size:255" json:"title"`
	ConstructionType    string    `gorm:"size:64" json:"construction_type"`
	District            string    `gorm:"size:128" json:"district"`
	Address             string    `gorm:"size:255" json:"address"`
	Area                float64   `json:"area"`
	Rooms               int       `json:"rooms"`
	Windows             int       `json:"windows"`
	Floor               int       `json:"floor"`
	TotalFloors         int       `json:"total_floors"`
	PricePerM2          float64   `json:"price_per_m2"`
	TotalPrice          float64   `json:"total_price"`
	Developer           string    `gorm:"size:255" json:"developer"`
	ContactPhone        string    `gorm:"size:64" json:"contact_phone"`
	HasTechPassport     string    `gorm:"size:16" json:"has_tech_passport"`
	HasRenovationPermit string    `gorm:"size:16" json:"has_renovation_permit"`
	Files               JSONB     `gorm:"type:jsonb" json:"files"`
	AuthorID            uint      `gorm:"index" json:"author_id"`
	CreatedAt           time.Time `json:"created_date"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (House) TableName() string { return "houses" }

// RealtyObject — обекти амлок бо координат (барои харита).
type RealtyObject struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Name                string    `gorm:"size:255" json:"name"`
	Description         string    `gorm:"type:text" json:"description"`
	Address             string    `gorm:"size:255" json:"address"`
	Lat                 float64   `gorm:"default:38.5598" json:"lat"`
	Lng                 float64   `gorm:"default:68.7870" json:"lng"`
	ConstructionType    string    `gorm:"size:64;default:новостройка" json:"construction_type"`
	District            string    `gorm:"size:128;default:н.Сино" json:"district"`
	Area                float64   `gorm:"default:0" json:"area"`
	Rooms               int       `gorm:"default:0" json:"rooms"`
	Windows             int       `gorm:"default:0" json:"windows"`
	Floor               int       `gorm:"default:0" json:"floor"`
	TotalFloors         int       `gorm:"default:0" json:"total_floors"`
	PricePerM2          float64   `gorm:"default:0" json:"price_per_m2"`
	TotalPrice          float64   `gorm:"default:0" json:"total_price"`
	Developer           string    `gorm:"size:255" json:"developer"`
	ContactPhone        string    `gorm:"size:64" json:"contact_phone"`
	HasTechPassport     string    `gorm:"size:16;default:нет" json:"has_tech_passport"`
	HasRenovationPermit string    `gorm:"size:16;default:нет" json:"has_renovation_permit"`
	AuthorID            uint      `gorm:"index" json:"author_id"`
	CreatedAt           time.Time `json:"created_date"`
	UpdatedAt           time.Time `json:"updated_date"`
}

func (RealtyObject) TableName() string { return "realty_objects" }

// RealtyBlock — блоки обекти амлок.
type RealtyBlock struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ObjectID   uint      `gorm:"index" json:"object_id"`
	Name       string    `gorm:"size:128" json:"name"`
	Code       string    `gorm:"size:64" json:"code"`
	OrderIndex int       `gorm:"default:0" json:"order_index"`
	CreatedAt  time.Time `json:"created_date"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (RealtyBlock) TableName() string { return "realty_blocks" }

// RealtyPricing — нархгузорӣ.
type RealtyPricing struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ObjectID        uint      `gorm:"index" json:"object_id"`
	BlockID         *uint     `gorm:"index" json:"block_id"`
	FloorNumber     int       `json:"floor_number"`
	FloorRangeStart int       `json:"floor_range_start"`
	FloorRangeEnd   int       `json:"floor_range_end"`
	PercentValue    int       `json:"percent_value"`
	PriceUSD        float64   `gorm:"default:0" json:"price_usd"`
	PriceTJS        float64   `gorm:"default:0" json:"price_tjs"`
	Currency        string    `gorm:"size:8;default:USD" json:"currency"`
	CreatedAt       time.Time `json:"created_date"`
	UpdatedAt       time.Time `json:"updated_date"`
}

func (RealtyPricing) TableName() string { return "realty_pricing" }

// RealtyLayout — планировкаҳо.
type RealtyLayout struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ObjectID     uint      `gorm:"index" json:"object_id"`
	RoomType     string    `gorm:"size:64" json:"room_type"`
	WindowsCount int       `json:"windows_count"`
	Area         float64   `json:"area"`
	PriceUSD     float64   `gorm:"default:0" json:"price_usd"`
	PriceTJS     float64   `gorm:"default:0" json:"price_tjs"`
	Description  string    `gorm:"type:text" json:"description"`
	CreatedAt    time.Time `json:"created_date"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (RealtyLayout) TableName() string { return "realty_layouts" }
