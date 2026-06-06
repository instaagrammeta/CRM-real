package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"github.com/instaagrammeta/crm-real/backend-go/internal/services"
	"github.com/instaagrammeta/crm-real/backend-go/internal/ws"
	"gorm.io/gorm"
)

type LidsHandler struct {
	DB       *gorm.DB
	Notifier *services.Notifier
}

// List — GET /api/lids
func (h *LidsHandler) List(c *gin.Context) {
	var lids []models.Lid
	if err := h.DB.Order("id DESC").Find(&lids).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, lids)
}

// Create — POST /api/lids
func (h *LidsHandler) Create(c *gin.Context) {
	var l models.Lid
	if err := c.ShouldBindJSON(&l); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	l.AuthorID = middleware.CurrentUserID(c)
	l.CreatedAt = time.Now()
	if err := h.DB.Create(&l).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if h.Notifier != nil {
		h.Notifier.Broadcast(ws.EventLeadCreated, l)
		h.Notifier.NotifyAdmins(services.Payload{
			Title:      "Лиди нав",
			Body:       l.ClientName + " — " + l.Phone,
			Type:       "lead",
			EntityType: "lid",
			EntityID:   l.ID,
		})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": l.ID, "lid": l})
}

// Update — PUT /api/lids/:id
func (h *LidsHandler) Update(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var data models.Lid
	if err := c.ShouldBindJSON(&data); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.Lid{}).Where("id = ?", id).Updates(map[string]any{
		"client_name": data.ClientName,
		"phone":       data.Phone,
		"topic":       data.Topic,
		"comment":     data.Comment,
		"source":      data.Source,
		"mortgage":    data.Mortgage,
		"box":         data.Box,
	}).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Delete — DELETE /api/lids/:id
func (h *LidsHandler) Delete(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.Lid{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
