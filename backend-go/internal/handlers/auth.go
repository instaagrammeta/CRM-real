package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/auth"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB     *gorm.DB
	JWTMgr *auth.Manager
}

type loginReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login — POST /api/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	var u models.User
	if err := h.DB.Where("login = ?", req.Login).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Неверный логин или пароль"})
		return
	}
	if !auth.CheckPassword(u.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Неверный логин или пароль"})
		return
	}
	token, exp, err := h.JWTMgr.Generate(u.ID, u.Login, u.Role, u.FullName)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "could not sign token")
		return
	}
	// также ставим httpOnly cookie для совместимости
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", token, int(h.JWTMgr.TTL().Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"role":    u.Role,
		"token":   token,
		"expires": exp,
		"user": gin.H{
			"id":        u.ID,
			"full_name": u.FullName,
			"login":     u.Login,
			"role":      u.Role,
			"category":  u.Category,
			"photo":     u.Photo,
		},
	})
}

// Logout — POST /api/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// CheckAuth — GET /api/check-auth
func (h *AuthHandler) CheckAuth(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"user_id":   uid,
		"role":      middleware.CurrentRole(c),
		"full_name": middleware.CurrentUserName(c),
	})
}

// Me — GET /api/me — корбари ҳозира.
func (h *AuthHandler) Me(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	var u models.User
	if err := h.DB.First(&u, uid).Error; err != nil {
		JSONError(c, http.StatusNotFound, "user not found")
		return
	}
	c.JSON(http.StatusOK, u)
}
