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

type KanbanHandler struct {
	DB       *gorm.DB
	Notifier *services.Notifier
}

// =================== Boards ===================

// ListBoards — GET /api/kanban/boards
func (h *KanbanHandler) ListBoards(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	var boards []models.KanbanBoard
	if err := h.DB.
		Where("is_archived = ? AND (is_public = ? OR author_id = ?)", false, true, uid).
		Order("id DESC").Find(&boards).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, boards)
}

// CreateBoard — POST /api/kanban/boards
func (h *KanbanHandler) CreateBoard(c *gin.Context) {
	var b models.KanbanBoard
	if err := c.ShouldBindJSON(&b); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	b.AuthorID = middleware.CurrentUserID(c)
	b.CreatedAt = time.Now()
	if err := h.DB.Create(&b).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": b.ID, "board": b})
}

// UpdateBoard — PUT /api/kanban/boards/:id
func (h *KanbanHandler) UpdateBoard(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var b models.KanbanBoard
	if err := c.ShouldBindJSON(&b); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.KanbanBoard{}).Where("id = ?", id).Updates(map[string]any{
		"title":     b.Title,
		"color":     b.Color,
		"is_public": b.IsPublic,
	}).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteBoard — DELETE /api/kanban/boards/:id
func (h *KanbanHandler) DeleteBoard(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	tx.Where("board_id = ?", id).Delete(&models.KanbanLead{})
	tx.Where("board_id = ?", id).Delete(&models.KanbanColumn{})
	tx.Where("board_id = ?", id).Delete(&models.KanbanBoardMember{})
	if err := tx.Delete(&models.KanbanBoard{}, id).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ArchiveBoard — PUT /api/kanban/boards/:id/archive
func (h *KanbanHandler) ArchiveBoard(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Model(&models.KanbanBoard{}).Where("id = ?", id).
		Update("is_archived", true).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =================== Columns ===================

// ListColumns — GET /api/kanban/boards/:id/columns
func (h *KanbanHandler) ListColumns(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var cols []models.KanbanColumn
	if err := h.DB.Where("board_id = ? AND is_archived = ?", id, false).
		Order("order_index ASC, id ASC").Find(&cols).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, cols)
}

// CreateColumn — POST /api/kanban/columns
func (h *KanbanHandler) CreateColumn(c *gin.Context) {
	var col models.KanbanColumn
	if err := c.ShouldBindJSON(&col); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	col.CreatedAt = time.Now()
	if err := h.DB.Create(&col).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": col.ID, "column": col})
}

// UpdateColumn — PUT /api/kanban/columns/:id
func (h *KanbanHandler) UpdateColumn(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var col models.KanbanColumn
	if err := c.ShouldBindJSON(&col); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	patch := map[string]any{}
	if col.Title != "" {
		patch["title"] = col.Title
	}
	if col.Color != "" {
		patch["color"] = col.Color
	}
	patch["order_index"] = col.OrderIndex
	if err := h.DB.Model(&models.KanbanColumn{}).Where("id = ?", id).Updates(patch).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteColumn — DELETE /api/kanban/columns/:id
func (h *KanbanHandler) DeleteColumn(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	tx.Where("column_id = ?", id).Delete(&models.KanbanLead{})
	if err := tx.Delete(&models.KanbanColumn{}, id).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdateColumnOrder — PUT /api/kanban/columns/:id/order { order_index }
func (h *KanbanHandler) UpdateColumnOrder(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var body struct {
		OrderIndex int `json:"order_index"`
	}
	_ = c.ShouldBindJSON(&body)
	if err := h.DB.Model(&models.KanbanColumn{}).Where("id = ?", id).
		Update("order_index", body.OrderIndex).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =================== Leads ===================

// ListLeads — GET /api/kanban/leads?board_id=&column_id=
func (h *KanbanHandler) ListLeads(c *gin.Context) {
	q := h.DB.Order("order_index ASC, id DESC")
	if v := c.Query("board_id"); v != "" {
		q = q.Where("board_id = ?", v)
	}
	if v := c.Query("column_id"); v != "" {
		q = q.Where("column_id = ?", v)
	}
	var leads []models.KanbanLead
	if err := q.Find(&leads).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, leads)
}

// CreateLead — POST /api/kanban/leads
func (h *KanbanHandler) CreateLead(c *gin.Context) {
	var l models.KanbanLead
	if err := c.ShouldBindJSON(&l); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	l.AuthorID = middleware.CurrentUserID(c)
	l.AuthorName = middleware.CurrentUserName(c)
	l.CreatedAt = time.Now()
	if err := h.DB.Create(&l).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if h.Notifier != nil {
		h.Notifier.Broadcast(ws.EventLeadCreated, l)
		h.Notifier.NotifyAdmins(services.Payload{
			Title:      "Лиди нав дар канбан",
			Body:       l.ClientName + " — " + l.Phone,
			Type:       "lead",
			EntityType: "kanban_lead",
			EntityID:   l.ID,
		})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": l.ID, "lead": l})
}

// UpdateLead — PUT /api/kanban/leads/:id
func (h *KanbanHandler) UpdateLead(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var l models.KanbanLead
	if err := c.ShouldBindJSON(&l); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	patch := map[string]any{
		"client_name": l.ClientName,
		"phone":       l.Phone,
		"topic":       l.Topic,
		"comment":     l.Comment,
		"source":      l.Source,
		"mortgage":    l.Mortgage,
		"box":         l.Box,
	}
	if err := h.DB.Model(&models.KanbanLead{}).Where("id = ?", id).Updates(patch).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if h.Notifier != nil {
		h.Notifier.Broadcast(ws.EventLeadUpdated, gin.H{"id": id})
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteLead — DELETE /api/kanban/leads/:id
func (h *KanbanHandler) DeleteLead(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.KanbanLead{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// MoveLead — POST /api/kanban/leads/move { lead_id, column_id, order_index }
func (h *KanbanHandler) MoveLead(c *gin.Context) {
	var body struct {
		LeadID     uint `json:"lead_id"`
		ColumnID   uint `json:"column_id"`
		OrderIndex int  `json:"order_index"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.KanbanLead{}).Where("id = ?", body.LeadID).Updates(map[string]any{
		"column_id":   body.ColumnID,
		"order_index": body.OrderIndex,
		"updated_at":  time.Now(),
	}).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if h.Notifier != nil {
		h.Notifier.Broadcast(ws.EventLeadMoved, body)
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdateLeadOrder — PUT /api/kanban/leads/:id/order
func (h *KanbanHandler) UpdateLeadOrder(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var body struct {
		OrderIndex int `json:"order_index"`
	}
	_ = c.ShouldBindJSON(&body)
	if err := h.DB.Model(&models.KanbanLead{}).Where("id = ?", id).
		Update("order_index", body.OrderIndex).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
