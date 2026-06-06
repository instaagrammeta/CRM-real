package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"gorm.io/gorm"
)

type PostsHandler struct {
	DB *gorm.DB
}

func (h *PostsHandler) List(c *gin.Context) {
	q := h.DB.Order("id DESC")
	if v := c.Query("user_id"); v != "" {
		q = q.Where("user_id = ?", v)
	}
	if v := c.Query("category"); v != "" {
		q = q.Where("category = ?", v)
	}
	if v := c.Query("project"); v != "" {
		q = q.Where("project = ?", v)
	}
	limit := QueryInt(c, "limit", 200)
	if limit > 1000 {
		limit = 1000
	}
	var rows []models.Post
	if err := q.Limit(limit).Find(&rows).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *PostsHandler) Create(c *gin.Context) {
	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	p.UserID = middleware.CurrentUserID(c)
	p.UserName = middleware.CurrentUserName(c)
	p.CreatedAt = time.Now()
	if err := h.DB.Create(&p).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": p.ID, "post": p})
}

func (h *PostsHandler) Update(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.DB.Model(&models.Post{}).Where("id = ?", id).Updates(&p).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *PostsHandler) Delete(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.Post{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
