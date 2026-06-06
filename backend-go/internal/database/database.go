package database

import (
	"fmt"
	"time"

	"github.com/instaagrammeta/crm-real/backend-go/internal/config"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// New opens a GORM connection to PostgreSQL and runs AutoMigrate for all models.
func New(cfg *config.Config) (*gorm.DB, error) {
	gormLogLevel := gormlogger.Warn
	if !cfg.IsProduction() {
		gormLogLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{
		Logger:                 gormlogger.Default.LogMode(gormLogLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	log.Info().Msg("database connected and migrated")

	return db, nil
}

// AutoMigrate creates / updates tables for all models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Request{},
		&models.Post{},
		&models.Message{},

		&models.RequestsBoard{},
		&models.RequestsColumn{},
		&models.RequestsItem{},

		&models.CompanyPhone{},
		&models.SimCard{},
		&models.SimTariff{},
		&models.TariffPayment{},

		&models.House{},

		&models.RealtyObject{},
		&models.RealtyBlock{},
		&models.RealtyPricing{},
		&models.RealtyLayout{},

		&models.ObjektProject{},
		&models.ObjektBlock{},
		&models.ObjektApartment{},

		&models.Lid{},

		&models.Folder{},
		&models.FolderFile{},

		&models.KanbanBoard{},
		&models.KanbanBoardMember{},
		&models.KanbanColumn{},
		&models.KanbanLead{},
		&models.KanbanLeadInteraction{},

		&models.Bank{},
		&models.MortgageCondition{},
		&models.InstallmentObject{},
		&models.InstallmentCondition{},

		&models.Notification{},
		&models.TelegramSubscriber{},
	)
}
