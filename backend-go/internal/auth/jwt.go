package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims — JWT-и стандартии корбар.
type Claims struct {
	UserID   uint   `json:"uid"`
	Login    string `json:"login"`
	Role     string `json:"role"`
	FullName string `json:"name"`
	jwt.RegisteredClaims
}

// Manager — менеджери JWT.
type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, hours int) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    time.Duration(hours) * time.Hour,
	}
}

func (m *Manager) Generate(userID uint, login, role, name string) (string, time.Time, error) {
	exp := time.Now().Add(m.ttl)
	claims := &Claims{
		UserID:   userID,
		Login:    login,
		Role:     role,
		FullName: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   login,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	return signed, exp, err
}

func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
