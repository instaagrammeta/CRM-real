package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/config"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"github.com/instaagrammeta/crm-real/backend-go/internal/services"
	"github.com/instaagrammeta/crm-real/backend-go/internal/uploads"
	"github.com/instaagrammeta/crm-real/backend-go/internal/ws"
	"gorm.io/gorm"
)

type ChatHandler struct {
	DB       *gorm.DB
	Cfg      *config.Config
	Notifier *services.Notifier
}

// List — GET /api/messages
func (h *ChatHandler) List(c *gin.Context) {
	limit := QueryInt(c, "limit", 200)
	if limit > 1000 {
		limit = 1000
	}
	var rows []models.Message
	if err := h.DB.Order("id DESC").Limit(limit).Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	// reverse to chronological
	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}
	c.JSON(http.StatusOK, rows)
}

// Send — POST /api/messages (json or multipart with optional file).
func (h *ChatHandler) Send(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	name := middleware.CurrentUserName(c)

	msg := models.Message{
		UserID:    uid,
		UserName:  name,
		CreatedAt: time.Now(),
	}

	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		msg.Message = c.PostForm("message")
		if file, err := c.FormFile("file"); err == nil && file != nil {
			path, err := uploads.SaveUpload(h.Cfg.UploadDir, "chat", file)
			if err == nil {
				msg.FilePath = path
				msg.FileName = file.Filename
				msg.FileType = file.Header.Get("Content-Type")
			}
		}
	} else {
		var body struct {
			Message  string `json:"message"`
			FilePath string `json:"file_path"`
			FileName string `json:"file_name"`
			FileType string `json:"file_type"`
		}
		_ = c.ShouldBindJSON(&body)
		msg.Message = body.Message
		msg.FilePath = body.FilePath
		msg.FileName = body.FileName
		msg.FileType = body.FileType
	}

	if msg.Message == "" && msg.FilePath == "" {
		JSONError(c, http.StatusBadRequest, "empty message")
		return
	}
	if err := h.DB.Create(&msg).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	if h.Notifier != nil {
		h.Notifier.Broadcast(ws.EventChatMessage, msg)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": msg})
}

// Update — PUT /api/messages/:id
func (h *ChatHandler) Update(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	uid := middleware.CurrentUserID(c)
	role := middleware.CurrentRole(c)

	var m models.Message
	if err := h.DB.First(&m, id).Error; err != nil {
		JSONError(c, http.StatusNotFound, "message not found")
		return
	}
	if role != "admin" && m.UserID != uid {
		JSONError(c, http.StatusForbidden, "forbidden")
		return
	}
	var body struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&m).Update("message", body.Message).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if h.Notifier != nil {
		h.Notifier.Broadcast(ws.EventChatUpdate, gin.H{"id": id, "message": body.Message})
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Delete — DELETE /api/messages/:id
func (h *ChatHandler) Delete(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	uid := middleware.CurrentUserID(c)
	role := middleware.CurrentRole(c)
	var m models.Message
	if err := h.DB.First(&m, id).Error; err != nil {
		JSONError(c, http.StatusNotFound, "message not found")
		return
	}
	if role != "admin" && m.UserID != uid {
		JSONError(c, http.StatusForbidden, "forbidden")
		return
	}
	if err := h.DB.Delete(&m).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if h.Notifier != nil {
		h.Notifier.Broadcast(ws.EventChatDelete, gin.H{"id": id})
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
