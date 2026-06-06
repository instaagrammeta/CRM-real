package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/auth"
	"github.com/rs/zerolog/log"
)

const (
	// CtxUserID — ключ ID-и корбар дар Gin context.
	CtxUserID    = "uid"
	CtxUserRole  = "role"
	CtxUserName  = "name"
	CtxUserLogin = "login"
)

// CORS configures permissive CORS for the API.
func CORS(origins []string) gin.HandlerFunc {
	cfg := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	if len(origins) == 1 && origins[0] == "*" {
		cfg.AllowAllOrigins = true
		cfg.AllowOrigins = nil
		cfg.AllowCredentials = false
	}
	return cors.New(cfg)
}

// Logger — простой реквест-логгер на zerolog.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		entry := log.Info()
		if c.Writer.Status() >= 500 {
			entry = log.Error()
		} else if c.Writer.Status() >= 400 {
			entry = log.Warn()
		}
		entry.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Msg("http")
	}
}

// Recovery converts panics into a 500 JSON response.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Interface("panic", r).Str("path", c.Request.URL.Path).Msg("panic")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}

// Auth парсит JWT из заголовка / cookie / query и кладёт claims в context.
func Auth(jwtMgr *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		claims, err := jwtMgr.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUserRole, claims.Role)
		c.Set(CtxUserName, claims.FullName)
		c.Set(CtxUserLogin, claims.Login)
		c.Next()
	}
}

// AdminOnly требует роль admin (использовать после Auth).
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(CtxUserRole)
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	// 1. Authorization: Bearer <token>
	if h := c.GetHeader("Authorization"); h != "" {
		if strings.HasPrefix(h, "Bearer ") {
			return strings.TrimPrefix(h, "Bearer "), nil
		}
	}
	// 2. cookie
	if v, err := c.Cookie("auth_token"); err == nil && v != "" {
		return v, nil
	}
	// 3. query (для WebSocket)
	if v := c.Query("token"); v != "" {
		return v, nil
	}
	return "", errors.New("no token")
}

// CurrentUserID — helper-и кӯтоҳ.
func CurrentUserID(c *gin.Context) uint {
	v, _ := c.Get(CtxUserID)
	if id, ok := v.(uint); ok {
		return id
	}
	return 0
}

// CurrentRole — helper-и кӯтоҳ.
func CurrentRole(c *gin.Context) string {
	v, _ := c.Get(CtxUserRole)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// CurrentUserName — helper-и кӯтоҳ.
func CurrentUserName(c *gin.Context) string {
	v, _ := c.Get(CtxUserName)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
