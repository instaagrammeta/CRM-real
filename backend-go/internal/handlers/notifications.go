package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"gorm.io/gorm"
)

type NotificationsHandler struct {
	DB *gorm.DB
}

// List — GET /api/notifications?unread=1
func (h *NotificationsHandler) List(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	q := h.DB.Where("user_id = ?", uid).Order("id DESC")
	if c.Query("unread") == "1" {
		q = q.Where("is_read = ?", false)
	}
	limit := QueryInt(c, "limit", 100)
	var rows []models.Notification
	if err := q.Limit(limit).Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

// MarkRead — POST /api/notifications/:id/read
func (h *NotificationsHandler) MarkRead(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	uid := middleware.CurrentUserID(c)
	if err := h.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, uid).
		Update("is_read", true).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// MarkAllRead — POST /api/notifications/read-all
func (h *NotificationsHandler) MarkAllRead(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	if err := h.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", uid, false).
		Update("is_read", true).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UnreadCount — GET /api/notifications/unread-count
func (h *NotificationsHandler) UnreadCount(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	var n int64
	h.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", uid, false).Count(&n)
	c.JSON(http.StatusOK, gin.H{"unread": n})
}

// Delete — DELETE /api/notifications/:id
func (h *NotificationsHandler) Delete(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	uid := middleware.CurrentUserID(c)
	if err := h.DB.Where("id = ? AND user_id = ?", id, uid).
		Delete(&models.Notification{}).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
