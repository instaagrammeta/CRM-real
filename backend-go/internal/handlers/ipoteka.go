package handlers

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"gorm.io/gorm"
)

type IpotekaHandler struct {
	DB *gorm.DB
}

// ============== Banks ==============

func (h *IpotekaHandler) ListBanks(c *gin.Context) {
	var rows []models.Bank
	q := h.DB.Order("order_index ASC, id ASC")
	if c.Query("active") == "true" {
		q = q.Where("is_active = ?", true)
	}
	if err := q.Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *IpotekaHandler) GetBank(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var b models.Bank
	if err := h.DB.First(&b, id).Error; err != nil {
		JSONError(c, http.StatusNotFound, "bank not found")
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *IpotekaHandler) GetBankBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var b models.Bank
	if err := h.DB.Where("slug = ?", slug).First(&b).Error; err != nil {
		JSONError(c, http.StatusNotFound, "bank not found")
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *IpotekaHandler) CreateBank(c *gin.Context) {
	var b models.Bank
	if err := c.ShouldBindJSON(&b); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if b.Slug == "" {
		b.Slug = makeSlug(b.Name)
	}
	b.Slug = ensureUniqueSlug(h.DB, "banks", b.Slug, 0)
	b.CreatedAt = time.Now()
	if err := h.DB.Create(&b).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": b.ID, "bank": b})
}

func (h *IpotekaHandler) UpdateBank(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var b models.Bank
	if err := c.ShouldBindJSON(&b); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if b.Slug == "" {
		b.Slug = makeSlug(b.Name)
	}
	b.Slug = ensureUniqueSlug(h.DB, "banks", b.Slug, id)
	if err := h.DB.Model(&models.Bank{}).Where("id = ?", id).Updates(&b).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *IpotekaHandler) DeleteBank(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	tx := h.DB.Begin()
	tx.Where("bank_id = ?", id).Delete(&models.MortgageCondition{})
	if err := tx.Delete(&models.Bank{}, id).Error; err != nil {
		tx.Rollback()
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ============== Mortgage Conditions ==============

func (h *IpotekaHandler) ListConditions(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var rows []models.MortgageCondition
	if err := h.DB.Where("bank_id = ?", id).Order("order_index, id").Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *IpotekaHandler) CreateCondition(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var cond models.MortgageCondition
	if err := c.ShouldBindJSON(&cond); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	cond.BankID = id
	cond.CreatedAt = time.Now()
	if err := h.DB.Create(&cond).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": cond.ID, "condition": cond})
}

func (h *IpotekaHandler) UpdateCondition(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var cond models.MortgageCondition
	if err := c.ShouldBindJSON(&cond); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.MortgageCondition{}).Where("id = ?", id).Updates(&cond).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *IpotekaHandler) DeleteCondition(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.MortgageCondition{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ============== helpers ==============

var slugRe = regexp.MustCompile(`[^a-zа-я0-9-]+`)

func makeSlug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugRe.ReplaceAllString(s, "")
	if s == "" {
		s = "x"
	}
	return s
}

func ensureUniqueSlug(db *gorm.DB, table, slug string, excludeID uint) string {
	candidate := slug
	for i := 1; i < 1000; i++ {
		var count int64
		q := db.Table(table).Where("slug = ?", candidate)
		if excludeID > 0 {
			q = q.Where("id <> ?", excludeID)
		}
		q.Count(&count)
		if count == 0 {
			return candidate
		}
		candidate = slug + "-" + itoa(i)
	}
	return candidate
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [16]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
