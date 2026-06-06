package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"gorm.io/gorm"
)

type ObjektHandler struct {
	DB *gorm.DB
}

// =================== Projects (шахматка) ===================

func (h *ObjektHandler) ListProjects(c *gin.Context) {
	var rows []models.ObjektProject
	if err := h.DB.Order("id DESC").Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *ObjektHandler) CreateProject(c *gin.Context) {
	var p models.ObjektProject
	if err := c.ShouldBindJSON(&p); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	p.AuthorID = middleware.CurrentUserID(c)
	p.CreatedAt = time.Now()
	if err := h.DB.Create(&p).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": p.ID, "project": p})
}

func (h *ObjektHandler) UpdateProject(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var p models.ObjektProject
	if err := c.ShouldBindJSON(&p); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.ObjektProject{}).Where("id = ?", id).Updates(&p).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ObjektHandler) DeleteProject(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	var blocks []models.ObjektBlock
	tx.Where("project_id = ?", id).Find(&blocks)
	for _, b := range blocks {
		tx.Where("block_id = ?", b.ID).Delete(&models.ObjektApartment{})
	}
	tx.Where("project_id = ?", id).Delete(&models.ObjektBlock{})
	if err := tx.Delete(&models.ObjektProject{}, id).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =================== Blocks ===================

func (h *ObjektHandler) ListBlocks(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var rows []models.ObjektBlock
	if err := h.DB.Where("project_id = ?", id).Order("order_index, id").
		Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *ObjektHandler) CreateBlock(c *gin.Context) {
	var b models.ObjektBlock
	if err := c.ShouldBindJSON(&b); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	b.CreatedAt = time.Now()
	if err := h.DB.Create(&b).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	// auto-create apartments
	h.autoCreateApartments(b)
	c.JSON(http.StatusOK, gin.H{"success": true, "id": b.ID, "block": b})
}

func (h *ObjektHandler) UpdateBlock(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var b models.ObjektBlock
	if err := c.ShouldBindJSON(&b); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.ObjektBlock{}).Where("id = ?", id).Updates(&b).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ObjektHandler) DeleteBlock(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	tx.Where("block_id = ?", id).Delete(&models.ObjektApartment{})
	if err := tx.Delete(&models.ObjektBlock{}, id).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =================== Apartments ===================

// ListApartments — GET /api/objekt/projects/:id/apartments
func (h *ObjektHandler) ListApartments(c *gin.Context) {
	pid, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var blocks []models.ObjektBlock
	if err := h.DB.Where("project_id = ?", pid).Find(&blocks).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(blocks) == 0 {
		c.JSON(http.StatusOK, []models.ObjektApartment{})
		return
	}
	ids := make([]uint, 0, len(blocks))
	for _, b := range blocks {
		ids = append(ids, b.ID)
	}
	var apts []models.ObjektApartment
	if err := h.DB.Where("block_id IN ?", ids).Order("block_id, floor").
		Find(&apts).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, apts)
}

// UpdateApartment — PUT /api/objekt/apartments/:id
func (h *ObjektHandler) UpdateApartment(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var a models.ObjektApartment
	if err := c.ShouldBindJSON(&a); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if a.Area > 0 && a.PricePerM2 > 0 && a.TotalPrice == 0 {
		a.TotalPrice = a.Area * a.PricePerM2
	}
	if err := h.DB.Model(&models.ObjektApartment{}).Where("id = ?", id).Updates(&a).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdateApartmentStatus — PATCH /api/objekt/apartments/:id/status
func (h *ObjektHandler) UpdateApartmentStatus(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.ObjektApartment{}).Where("id = ?", id).
		Update("status", body.Status).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =================== Helpers ===================

func (h *ObjektHandler) autoCreateApartments(b models.ObjektBlock) {
	now := time.Now()
	for floor := b.FloorFrom; floor <= b.FloorTo; floor++ {
		var existing models.ObjektApartment
		err := h.DB.Where("block_id = ? AND floor = ?", b.ID, floor).First(&existing).Error
		if err == nil {
			continue
		}
		apt := models.ObjektApartment{
			BlockID:    b.ID,
			Floor:      floor,
			Area:       b.DefaultArea,
			Rooms:      b.DefaultRooms,
			Windows:    b.DefaultWindows,
			PricePerM2: b.DefaultPricePerM2,
			TotalPrice: b.DefaultArea * b.DefaultPricePerM2,
			Status:     "free",
			Balcony:    b.DefaultBalcony,
			Bathroom:   b.DefaultBathroom,
			PlanImage:  b.DefaultPlanImage,
			CreatedAt:  now,
		}
		_ = h.DB.Create(&apt).Error
	}
}
