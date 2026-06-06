package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/auth"
	"github.com/instaagrammeta/crm-real/backend-go/internal/config"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"github.com/instaagrammeta/crm-real/backend-go/internal/services/telegram"
	"github.com/instaagrammeta/crm-real/backend-go/internal/uploads"
	"gorm.io/gorm"
)

type UsersHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
	TG  *telegram.Bot
}

// List — GET /api/users
func (h *UsersHandler) List(c *gin.Context) {
	role := middleware.CurrentRole(c)
	uid := middleware.CurrentUserID(c)
	var users []models.User
	q := h.DB.Order("id ASC")
	if role != "admin" {
		q = q.Where("id = ?", uid)
	}
	if err := q.Find(&users).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

// Listing — GET /api/users/list (тонкий список для селекторов).
func (h *UsersHandler) Listing(c *gin.Context) {
	type row struct {
		ID       uint   `json:"id"`
		FullName string `json:"full_name"`
		Category string `json:"category"`
		Photo    string `json:"photo"`
		Role     string `json:"role"`
	}
	var rows []row
	if err := h.DB.Model(&models.User{}).
		Select("id, full_name, category, photo, role").
		Order("full_name ASC").
		Scan(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

type userInput struct {
	FullName       string   `json:"full_name"`
	Age            int      `json:"age"`
	PersonalPhones []string `json:"personal_phones"`
	WorkPhones     []string `json:"work_phones"`
	Login          string   `json:"login"`
	Password       string   `json:"password"`
	Photo          string   `json:"photo"`
	Category       string   `json:"category"`
	Role           string   `json:"role"`
}

// Create — POST /api/users (admin).
func (h *UsersHandler) Create(c *gin.Context) {
	var data userInput

	// поддерживаем и multipart, и json
	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		data.FullName = c.PostForm("full_name")
		data.Login = c.PostForm("login")
		data.Password = c.PostForm("password")
		data.Category = c.PostForm("category")
		data.Role = c.DefaultPostForm("role", "employee")
		if v := c.PostForm("age"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				data.Age = n
			}
		}
		data.PersonalPhones = splitCSV(c.PostForm("personal_phones"))
		data.WorkPhones = splitCSV(c.PostForm("work_phones"))
		// upload фото
		if file, err := c.FormFile("photo_file"); err == nil && file != nil {
			path, err := uploads.SaveUpload(h.Cfg.UploadDir, "photos", file)
			if err == nil {
				data.Photo = path
			}
		}
	} else {
		if err := c.ShouldBindJSON(&data); err != nil {
			JSONError(c, http.StatusBadRequest, "invalid body")
			return
		}
	}

	if data.FullName == "" || data.Login == "" || data.Password == "" {
		JSONError(c, http.StatusBadRequest, "full_name, login, password required")
		return
	}

	hash, err := auth.HashPassword(data.Password)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	role := data.Role
	if role == "" {
		role = "employee"
	}
	u := models.User{
		FullName:       data.FullName,
		Age:            data.Age,
		PersonalPhones: data.PersonalPhones,
		WorkPhones:     data.WorkPhones,
		Login:          data.Login,
		Password:       hash,
		Photo:          data.Photo,
		Category:       data.Category,
		Role:           role,
		CreatedAt:      time.Now(),
	}
	if err := h.DB.Create(&u).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": u.ID, "message": "Сотрудник добавлен"})
}

// Update — PUT /api/users/:id (admin).
func (h *UsersHandler) Update(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var u models.User
	if err := h.DB.First(&u, id).Error; err != nil {
		JSONError(c, http.StatusNotFound, "user not found")
		return
	}
	var data userInput
	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		data.FullName = c.PostForm("full_name")
		data.Login = c.PostForm("login")
		data.Password = c.PostForm("password")
		data.Category = c.PostForm("category")
		data.Role = c.PostForm("role")
		if v := c.PostForm("age"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				data.Age = n
			}
		}
		data.PersonalPhones = splitCSV(c.PostForm("personal_phones"))
		data.WorkPhones = splitCSV(c.PostForm("work_phones"))
		if file, err := c.FormFile("photo_file"); err == nil && file != nil {
			path, err := uploads.SaveUpload(h.Cfg.UploadDir, "photos", file)
			if err == nil {
				data.Photo = path
			}
		}
	} else {
		_ = c.ShouldBindJSON(&data)
	}

	patch := map[string]any{}
	if data.FullName != "" {
		patch["full_name"] = data.FullName
	}
	if data.Login != "" {
		patch["login"] = data.Login
	}
	if data.Category != "" {
		patch["category"] = data.Category
	}
	if data.Role != "" {
		patch["role"] = data.Role
	}
	if data.Photo != "" {
		patch["photo"] = data.Photo
	}
	if data.Age != 0 {
		patch["age"] = data.Age
	}
	if data.PersonalPhones != nil {
		patch["personal_phones"] = models.StringSlice(data.PersonalPhones)
	}
	if data.WorkPhones != nil {
		patch["work_phones"] = models.StringSlice(data.WorkPhones)
	}
	if data.Password != "" {
		hash, err := auth.HashPassword(data.Password)
		if err != nil {
			JSONError(c, http.StatusInternalServerError, err.Error())
			return
		}
		patch["password"] = hash
	}
	if len(patch) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	if err := h.DB.Model(&u).Updates(patch).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Delete — DELETE /api/users/:id (admin).
func (h *UsersHandler) Delete(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.User{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =========================================================
// TELEGRAM LINKING
// =========================================================

// GenerateTelegramLink — POST /api/me/telegram/link  -> ngz returns URL.
func (h *UsersHandler) GenerateTelegramLink(c *gin.Context) {
	if h.TG == nil || !h.TG.Enabled() {
		JSONError(c, http.StatusServiceUnavailable, "telegram bot is not configured")
		return
	}
	uid := middleware.CurrentUserID(c)
	token, err := h.TG.GenerateLinkToken(h.DB, uid)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"token":    token,
		"link_url": h.TG.LinkURL(token),
		"command":  "/start " + token,
	})
}

// ListTelegramSubscribers — GET /api/me/telegram.
func (h *UsersHandler) ListTelegramSubscribers(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	var subs []models.TelegramSubscriber
	if err := h.DB.Where("user_id = ?", uid).Find(&subs).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, subs)
}

// UnlinkTelegram — DELETE /api/me/telegram/:id.
func (h *UsersHandler) UnlinkTelegram(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	uid := middleware.CurrentUserID(c)
	if err := h.DB.Where("id = ? AND user_id = ?", id, uid).
		Delete(&models.TelegramSubscriber{}).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =========================================================
// helpers
// =========================================================

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
