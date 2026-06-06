package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"gorm.io/gorm"
)

type SimHandler struct {
	DB *gorm.DB
}

// =============== Company phones ===============

func (h *SimHandler) ListPhones(c *gin.Context) {
	var rows []models.CompanyPhone
	if err := h.DB.Order("id DESC").Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *SimHandler) CreatePhone(c *gin.Context) {
	var p models.CompanyPhone
	if err := c.ShouldBindJSON(&p); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	p.CreatedAt = time.Now()
	if err := h.DB.Create(&p).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": p.ID, "phone": p})
}

func (h *SimHandler) UpdatePhone(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var p models.CompanyPhone
	if err := c.ShouldBindJSON(&p); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.CompanyPhone{}).Where("id = ?", id).Updates(&p).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *SimHandler) DeletePhone(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.CompanyPhone{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =============== SIM cards ===============

func (h *SimHandler) ListSims(c *gin.Context) {
	var rows []models.SimCard
	if err := h.DB.Order("id DESC").Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *SimHandler) CreateSim(c *gin.Context) {
	var s models.SimCard
	if err := c.ShouldBindJSON(&s); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	s.CreatedAt = time.Now()
	if err := h.DB.Create(&s).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": s.ID, "sim": s})
}

func (h *SimHandler) UpdateSim(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var s models.SimCard
	if err := c.ShouldBindJSON(&s); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.SimCard{}).Where("id = ?", id).Updates(&s).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *SimHandler) DeleteSim(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	tx.Where("sim_id = ?", id).Delete(&models.TariffPayment{})
	tx.Where("sim_id = ?", id).Delete(&models.SimTariff{})
	if err := tx.Delete(&models.SimCard{}, id).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =============== Tariffs ===============

func (h *SimHandler) ListTariffs(c *gin.Context) {
	q := h.DB.Order("id DESC")
	if v := c.Query("sim_id"); v != "" {
		q = q.Where("sim_id = ?", v)
	}
	var rows []models.SimTariff
	if err := q.Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *SimHandler) CreateTariff(c *gin.Context) {
	var t models.SimTariff
	if err := c.ShouldBindJSON(&t); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	t.CreatedAt = time.Now()
	if err := h.DB.Create(&t).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": t.ID, "tariff": t})
}

func (h *SimHandler) UpdateTariff(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var t models.SimTariff
	if err := c.ShouldBindJSON(&t); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.SimTariff{}).Where("id = ?", id).Updates(&t).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *SimHandler) DeleteTariff(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.SimTariff{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =============== Payments ===============

func (h *SimHandler) ListPayments(c *gin.Context) {
	q := h.DB.Order("id DESC")
	if v := c.Query("sim_id"); v != "" {
		q = q.Where("sim_id = ?", v)
	}
	var rows []models.TariffPayment
	if err := q.Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *SimHandler) DeletePayment(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.TariffPayment{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =============== Stats ===============

// Stats — GET /api/sim-phones/stats
func (h *SimHandler) Stats(c *gin.Context) {
	var totalPhones, totalSims, activeSims, expiredTariffs int64
	h.DB.Model(&models.CompanyPhone{}).Count(&totalPhones)
	h.DB.Model(&models.SimCard{}).Count(&totalSims)
	h.DB.Model(&models.SimCard{}).Where("status = ?", "active").Count(&activeSims)
	h.DB.Model(&models.SimTariff{}).Where("end_date < ?", time.Now().Format("2006-01-02")).
		Count(&expiredTariffs)
	c.JSON(http.StatusOK, gin.H{
		"phones":          totalPhones,
		"sims":            totalSims,
		"active_sims":     activeSims,
		"expired_tariffs": expiredTariffs,
	})
}
