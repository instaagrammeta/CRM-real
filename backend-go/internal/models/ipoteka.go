package models

import "time"

// Bank — бонк барои ипотека.
type Bank struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255" json:"name"`
	Slug        string    `gorm:"size:255;uniqueIndex" json:"slug"`
	Logo        string    `gorm:"size:1024" json:"logo"`
	Description string    `gorm:"type:text" json:"description"`
	Phone       string    `gorm:"size:64" json:"phone"`
	Website     string    `gorm:"size:255" json:"website"`
	Address     string    `gorm:"size:255" json:"address"`
	OrderIndex  int       `gorm:"default:0" json:"order_index"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_date"`
	UpdatedAt   time.Time `json:"updated_date"`
}

func (Bank) TableName() string { return "banks" }

// MortgageCondition — шарти ипотека.
type MortgageCondition struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	BankID         uint      `gorm:"index" json:"bank_id"`
	Title          string    `gorm:"size:255" json:"title"`
	InterestRate   float64   `json:"interest_rate"`
	MinDownPayment float64   `json:"min_down_payment"`
	MaxAmount      float64   `json:"max_amount"`
	MaxTermYears   int       `json:"max_term_years"`
	Currency       string    `gorm:"size:8;default:TJS" json:"currency"`
	Requirements   string    `gorm:"type:text" json:"requirements"`
	Documents      string    `gorm:"type:text" json:"documents"`
	AdditionalInfo string    `gorm:"type:text" json:"additional_info"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	OrderIndex     int       `gorm:"default:0" json:"order_index"`
	CreatedAt      time.Time `json:"created_date"`
	UpdatedAt      time.Time `json:"updated_date"`
}

func (MortgageCondition) TableName() string { return "mortgage_conditions" }

// InstallmentObject — объекти рассрочка.
type InstallmentObject struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255" json:"name"`
	Slug        string    `gorm:"size:255;uniqueIndex" json:"slug"`
	Logo        string    `gorm:"size:1024" json:"logo"`
	Description string    `gorm:"type:text" json:"description"`
	Phone       string    `gorm:"size:64" json:"phone"`
	Website     string    `gorm:"size:255" json:"website"`
	Address     string    `gorm:"size:255" json:"address"`
	Developer   string    `gorm:"size:255" json:"developer"`
	OrderIndex  int       `gorm:"default:0" json:"order_index"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_date"`
	UpdatedAt   time.Time `json:"updated_date"`
}

func (InstallmentObject) TableName() string { return "installment_objects" }

// InstallmentCondition — шарти рассрочка.
type InstallmentCondition struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ObjectID       uint      `gorm:"index" json:"object_id"`
	Title          string    `gorm:"size:255" json:"title"`
	MinDownPayment float64   `json:"min_down_payment"`
	MaxTermMonths  int       `json:"max_term_months"`
	InterestRate   float64   `json:"interest_rate"`
	Currency       string    `gorm:"size:8;default:TJS" json:"currency"`
	Requirements   string    `gorm:"type:text" json:"requirements"`
	Documents      string    `gorm:"type:text" json:"documents"`
	AdditionalInfo string    `gorm:"type:text" json:"additional_info"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	OrderIndex     int       `gorm:"default:0" json:"order_index"`
	CreatedAt      time.Time `json:"created_date"`
	UpdatedAt      time.Time `json:"updated_date"`
}

func (InstallmentCondition) TableName() string { return "installment_conditions" }
