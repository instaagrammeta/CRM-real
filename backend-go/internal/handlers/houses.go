package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"gorm.io/gorm"
)

type HousesHandler struct {
	DB *gorm.DB
}

// List — GET /api/houses
func (h *HousesHandler) List(c *gin.Context) {
	var rows []models.House
	if err := h.DB.Order("id DESC").Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

// Create — POST /api/houses
func (h *HousesHandler) Create(c *gin.Context) {
	var x models.House
	if err := c.ShouldBindJSON(&x); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	x.AuthorID = middleware.CurrentUserID(c)
	x.CreatedAt = time.Now()
	if err := h.DB.Create(&x).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": x.ID, "house": x})
}

// Update — PUT /api/houses/:id
func (h *HousesHandler) Update(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var x models.House
	if err := c.ShouldBindJSON(&x); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.House{}).Where("id = ?", id).Updates(&x).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Delete — DELETE /api/houses/:id
func (h *HousesHandler) Delete(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.House{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
