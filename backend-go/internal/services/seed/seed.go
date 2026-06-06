package seed

import (
	"errors"

	"github.com/instaagrammeta/crm-real/backend-go/internal/auth"
	"github.com/instaagrammeta/crm-real/backend-go/internal/config"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// EnsureAdmin creates the default admin user if not present.
func EnsureAdmin(db *gorm.DB, cfg *config.Config) error {
	var existing models.User
	err := db.Where("login = ?", cfg.AdminLogin).First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hash, err := auth.HashPassword(cfg.AdminPassword)
	if err != nil {
		return err
	}
	u := models.User{
		FullName:       cfg.AdminFullName,
		Login:          cfg.AdminLogin,
		Password:       hash,
		Role:           "admin",
		Category:       "Руководство",
		Age:            30,
		PersonalPhones: models.StringSlice{},
		WorkPhones:     models.StringSlice{},
	}
	if err := db.Create(&u).Error; err != nil {
		return err
	}
	log.Info().Str("login", cfg.AdminLogin).Msg("admin user seeded")
	return nil
}
