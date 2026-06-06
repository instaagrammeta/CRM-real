package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"gorm.io/gorm"
)

type RasrochkaHandler struct {
	DB *gorm.DB
}

// ============== Objects ==============

func (h *RasrochkaHandler) ListObjects(c *gin.Context) {
	q := h.DB.Order("order_index ASC, id ASC")
	if c.Query("active") == "true" {
		q = q.Where("is_active = ?", true)
	}
	var rows []models.InstallmentObject
	if err := q.Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *RasrochkaHandler) GetObject(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var o models.InstallmentObject
	if err := h.DB.First(&o, id).Error; err != nil {
		JSONError(c, http.StatusNotFound, "object not found")
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *RasrochkaHandler) GetObjectBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var o models.InstallmentObject
	if err := h.DB.Where("slug = ?", slug).First(&o).Error; err != nil {
		JSONError(c, http.StatusNotFound, "object not found")
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *RasrochkaHandler) CreateObject(c *gin.Context) {
	var o models.InstallmentObject
	if err := c.ShouldBindJSON(&o); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if o.Slug == "" {
		o.Slug = makeSlug(o.Name)
	}
	o.Slug = ensureUniqueSlug(h.DB, "installment_objects", o.Slug, 0)
	o.CreatedAt = time.Now()
	if err := h.DB.Create(&o).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": o.ID, "object": o})
}

func (h *RasrochkaHandler) UpdateObject(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var o models.InstallmentObject
	if err := c.ShouldBindJSON(&o); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if o.Slug == "" {
		o.Slug = makeSlug(o.Name)
	}
	o.Slug = ensureUniqueSlug(h.DB, "installment_objects", o.Slug, id)
	if err := h.DB.Model(&models.InstallmentObject{}).Where("id = ?", id).Updates(&o).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *RasrochkaHandler) DeleteObject(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	tx.Where("object_id = ?", id).Delete(&models.InstallmentCondition{})
	if err := tx.Delete(&models.InstallmentObject{}, id).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ============== Conditions ==============

func (h *RasrochkaHandler) ListConditions(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var rows []models.InstallmentCondition
	if err := h.DB.Where("object_id = ?", id).Order("order_index, id").
		Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *RasrochkaHandler) CreateCondition(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var cond models.InstallmentCondition
	if err := c.ShouldBindJSON(&cond); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	cond.ObjectID = id
	cond.CreatedAt = time.Now()
	if err := h.DB.Create(&cond).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": cond.ID, "condition": cond})
}

func (h *RasrochkaHandler) UpdateCondition(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var cond models.InstallmentCondition
	if err := c.ShouldBindJSON(&cond); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.InstallmentCondition{}).Where("id = ?", id).
		Updates(&cond).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *RasrochkaHandler) DeleteCondition(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.InstallmentCondition{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
