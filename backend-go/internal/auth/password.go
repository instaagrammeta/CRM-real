package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword - bcrypt cost 12.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword проверяет пароль. Поддерживает старый sha256-формат
// для совместимости с первоначальной Flask-базой.
func CheckPassword(stored, plain string) bool {
	if stored == "" {
		return false
	}
	if isBcrypt(stored) {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(plain)) == nil
	}
	// legacy sha256 (Flask-версия)
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:]) == stored
}

func isBcrypt(s string) bool {
	return len(s) >= 4 && (s[:4] == "$2a$" || s[:4] == "$2b$" || s[:4] == "$2y$")
}
