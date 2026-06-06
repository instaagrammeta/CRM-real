package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"gorm.io/gorm"
)

type RealtyHandler struct {
	DB *gorm.DB
}

// =================== Objects ===================

func (h *RealtyHandler) ListObjects(c *gin.Context) {
	var rows []models.RealtyObject
	if err := h.DB.Order("id DESC").Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *RealtyHandler) CreateObject(c *gin.Context) {
	var o models.RealtyObject
	if err := c.ShouldBindJSON(&o); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	o.AuthorID = middleware.CurrentUserID(c)
	o.CreatedAt = time.Now()
	if err := h.DB.Create(&o).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": o.ID, "object": o})
}

func (h *RealtyHandler) UpdateObject(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var o models.RealtyObject
	if err := c.ShouldBindJSON(&o); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.RealtyObject{}).Where("id = ?", id).Updates(&o).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *RealtyHandler) DeleteObject(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	tx.Where("object_id = ?", id).Delete(&models.RealtyLayout{})
	tx.Where("object_id = ?", id).Delete(&models.RealtyPricing{})
	tx.Where("object_id = ?", id).Delete(&models.RealtyBlock{})
	if err := tx.Delete(&models.RealtyObject{}, id).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =================== Blocks ===================

func (h *RealtyHandler) ListBlocks(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var rows []models.RealtyBlock
	if err := h.DB.Where("object_id = ?", id).Order("order_index, id").Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *RealtyHandler) CreateBlock(c *gin.Context) {
	var b models.RealtyBlock
	if err := c.ShouldBindJSON(&b); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	b.CreatedAt = time.Now()
	if err := h.DB.Create(&b).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": b.ID, "block": b})
}

func (h *RealtyHandler) UpdateBlock(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var b models.RealtyBlock
	if err := c.ShouldBindJSON(&b); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.RealtyBlock{}).Where("id = ?", id).Updates(&b).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *RealtyHandler) DeleteBlock(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.RealtyBlock{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =================== Pricing ===================

func (h *RealtyHandler) ListPricing(c *gin.Context) {
	q := h.DB.Order("id DESC")
	if v := c.Query("object_id"); v != "" {
		q = q.Where("object_id = ?", v)
	}
	if v := c.Query("block_id"); v != "" {
		q = q.Where("block_id = ?", v)
	}
	var rows []models.RealtyPricing
	if err := q.Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *RealtyHandler) CreatePricing(c *gin.Context) {
	var p models.RealtyPricing
	if err := c.ShouldBindJSON(&p); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	p.CreatedAt = time.Now()
	if err := h.DB.Create(&p).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": p.ID, "pricing": p})
}

func (h *RealtyHandler) UpdatePricing(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var p models.RealtyPricing
	if err := c.ShouldBindJSON(&p); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.RealtyPricing{}).Where("id = ?", id).Updates(&p).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *RealtyHandler) DeletePricing(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.RealtyPricing{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =================== Layouts ===================

func (h *RealtyHandler) ListLayouts(c *gin.Context) {
	q := h.DB.Order("id DESC")
	if v := c.Query("object_id"); v != "" {
		q = q.Where("object_id = ?", v)
	}
	var rows []models.RealtyLayout
	if err := q.Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *RealtyHandler) CreateLayout(c *gin.Context) {
	var l models.RealtyLayout
	if err := c.ShouldBindJSON(&l); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	l.CreatedAt = time.Now()
	if err := h.DB.Create(&l).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": l.ID, "layout": l})
}

func (h *RealtyHandler) UpdateLayout(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var l models.RealtyLayout
	if err := c.ShouldBindJSON(&l); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.RealtyLayout{}).Where("id = ?", id).Updates(&l).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *RealtyHandler) DeleteLayout(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.RealtyLayout{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
