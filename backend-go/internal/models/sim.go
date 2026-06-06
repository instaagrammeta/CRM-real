package models

import "time"

// CompanyPhone — телефони ширкат.
type CompanyPhone struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Model       string    `gorm:"size:128" json:"model"`
	PhoneID     string    `gorm:"size:128;uniqueIndex" json:"phone_id"`
	AssignedTo  *uint     `gorm:"index" json:"assigned_to"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:32;default:free" json:"status"`
	CreatedAt   time.Time `json:"created_date"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CompanyPhone) TableName() string { return "company_phones" }

// SimCard — SIM карта.
type SimCard struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PhoneNumber string    `gorm:"size:64;uniqueIndex" json:"phone_number"`
	Operator    string    `gorm:"size:64" json:"operator"`
	AssignedTo  *uint     `gorm:"index" json:"assigned_to"`
	PhoneID     *uint     `gorm:"index" json:"phone_id"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:32;default:active" json:"status"`
	CreatedAt   time.Time `json:"created_date"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SimCard) TableName() string { return "sim_cards" }

// SimTariff — тарифи SIM-карта.
type SimTariff struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SimID     uint      `gorm:"index" json:"sim_id"`
	Minutes   int       `gorm:"default:0" json:"minutes"`
	GB        float64   `gorm:"default:0" json:"gb"`
	SMS       int       `gorm:"default:0" json:"sms"`
	Cost      float64   `gorm:"default:0" json:"cost"`
	StartDate string    `gorm:"size:32" json:"start_date"`
	EndDate   string    `gorm:"size:32" json:"end_date"`
	Status    string    `gorm:"size:32;default:active" json:"status"`
	CreatedAt time.Time `json:"created_date"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SimTariff) TableName() string { return "sim_tariffs" }

// TariffPayment — пардохти тариф.
type TariffPayment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SimID       uint      `gorm:"index" json:"sim_id"`
	TariffID    *uint     `gorm:"index" json:"tariff_id"`
	Amount      float64   `gorm:"default:0" json:"amount"`
	PaymentDate string    `gorm:"size:32" json:"payment_date"`
	StartDate   string    `gorm:"size:32" json:"start_date"`
	EndDate     string    `gorm:"size:32" json:"end_date"`
	Status      string    `gorm:"size:32;default:paid" json:"status"`
	CreatedAt   time.Time `json:"created_date"`
}

func (TariffPayment) TableName() string { return "tariff_payments" }
