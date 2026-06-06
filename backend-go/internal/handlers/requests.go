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

type RequestsHandler struct {
	DB       *gorm.DB
	Notifier *services.Notifier
}

// ====================== legacy /api/requests ======================

func (h *RequestsHandler) List(c *gin.Context) {
	var rows []models.Request
	if err := h.DB.Order("id DESC").Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *RequestsHandler) Create(c *gin.Context) {
	var r models.Request
	if err := c.ShouldBindJSON(&r); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	r.AuthorID = middleware.CurrentUserID(c)
	r.CreatedAt = time.Now()
	if err := h.DB.Create(&r).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if h.Notifier != nil {
		h.Notifier.Broadcast(ws.EventRequestCreated, r)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": r.ID, "request": r})
}

func (h *RequestsHandler) Update(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var r models.Request
	if err := c.ShouldBindJSON(&r); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.Request{}).Where("id = ?", id).Updates(&r).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *RequestsHandler) Delete(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.Request{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ====================== new kanban-style requests ======================

// ListBoards — GET /api/requests-boards
func (h *RequestsHandler) ListBoards(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	var boards []models.RequestsBoard
	if err := h.DB.Where("is_public = ? OR author_id = ?", true, uid).
		Order("id DESC").Find(&boards).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, boards)
}

// CreateBoard — POST /api/requests-boards
func (h *RequestsHandler) CreateBoard(c *gin.Context) {
	var b models.RequestsBoard
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

// UpdateBoard — PUT /api/requests-boards/:id
func (h *RequestsHandler) UpdateBoard(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var b models.RequestsBoard
	if err := c.ShouldBindJSON(&b); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.RequestsBoard{}).Where("id = ?", id).Updates(map[string]any{
		"title":     b.Title,
		"color":     b.Color,
		"is_public": b.IsPublic,
	}).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteBoard — DELETE /api/requests-boards/:id
func (h *RequestsHandler) DeleteBoard(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	tx.Where("board_id = ?", id).Delete(&models.RequestsItem{})
	tx.Where("board_id = ?", id).Delete(&models.RequestsColumn{})
	if err := tx.Delete(&models.RequestsBoard{}, id).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ListColumnsForBoard — GET /api/requests-boards/:id/columns
func (h *RequestsHandler) ListColumnsForBoard(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var cols []models.RequestsColumn
	if err := h.DB.Where("board_id = ?", id).Order("order_index ASC, id ASC").
		Find(&cols).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, cols)
}

// CreateColumn — POST /api/requests-columns
func (h *RequestsHandler) CreateColumn(c *gin.Context) {
	var col models.RequestsColumn
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

// UpdateColumn — PUT /api/requests-columns/:id
func (h *RequestsHandler) UpdateColumn(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var col models.RequestsColumn
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
	if err := h.DB.Model(&models.RequestsColumn{}).Where("id = ?", id).Updates(patch).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteColumn — DELETE /api/requests-columns/:id
func (h *RequestsHandler) DeleteColumn(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	tx.Where("column_id = ?", id).Delete(&models.RequestsItem{})
	if err := tx.Delete(&models.RequestsColumn{}, id).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ListItemsForBoard — GET /api/requests-board/:id/requests
func (h *RequestsHandler) ListItemsForBoard(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var items []models.RequestsItem
	if err := h.DB.Where("board_id = ?", id).
		Order("column_id ASC, order_index ASC, id DESC").
		Find(&items).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateItem — POST /api/requests-new
func (h *RequestsHandler) CreateItem(c *gin.Context) {
	var r models.RequestsItem
	if err := c.ShouldBindJSON(&r); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	r.AuthorID = middleware.CurrentUserID(c)
	r.AuthorName = middleware.CurrentUserName(c)
	r.CreatedAt = time.Now()
	if err := h.DB.Create(&r).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if h.Notifier != nil {
		h.Notifier.Broadcast(ws.EventRequestCreated, r)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": r.ID, "request": r})
}

// UpdateItem — PUT /api/requests-update/:id
func (h *RequestsHandler) UpdateItem(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var r models.RequestsItem
	if err := c.ShouldBindJSON(&r); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.RequestsItem{}).Where("id = ?", id).Updates(&r).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteItem — DELETE /api/requests-delete/:id
func (h *RequestsHandler) DeleteItem(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.RequestsItem{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// MoveItem — POST /api/requests-move
func (h *RequestsHandler) MoveItem(c *gin.Context) {
	var body struct {
		RequestID  uint `json:"request_id"`
		ColumnID   uint `json:"column_id"`
		OrderIndex int  `json:"order_index"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.RequestsItem{}).Where("id = ?", body.RequestID).Updates(map[string]any{
		"column_id":   body.ColumnID,
		"order_index": body.OrderIndex,
		"updated_at":  time.Now(),
	}).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if h.Notifier != nil {
		h.Notifier.Broadcast(ws.EventRequestMoved, body)
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
