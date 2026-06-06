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

type TasksHandler struct {
	DB       *gorm.DB
	Notifier *services.Notifier
}

// List — GET /api/tasks
func (h *TasksHandler) List(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	role := middleware.CurrentRole(c)

	var tasks []models.Task
	q := h.DB.Order("id DESC")
	if role != "admin" {
		q = q.Where("author_id = ? OR executor_id = ?", uid, uid)
	}
	if err := q.Find(&tasks).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tasks)
}

type taskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ExecutorID  uint   `json:"executor_id"`
	Photo       string `json:"photo"`
	Status      string `json:"status"`
}

// Create — POST /api/tasks
func (h *TasksHandler) Create(c *gin.Context) {
	var data taskInput
	if err := c.ShouldBindJSON(&data); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid body")
		return
	}
	uid := middleware.CurrentUserID(c)
	t := models.Task{
		Title:       data.Title,
		Description: data.Description,
		AuthorID:    uid,
		ExecutorID:  data.ExecutorID,
		Photo:       data.Photo,
		Status:      firstNonEmpty(data.Status, "new"),
		CreatedAt:   time.Now(),
	}
	if err := h.DB.Create(&t).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	if data.ExecutorID != 0 && data.ExecutorID != uid && h.Notifier != nil {
		_ = h.Notifier.Notify(services.Payload{
			UserID:     data.ExecutorID,
			Title:      "Вазифаи нав",
			Body:       t.Title,
			Type:       "task",
			EntityType: "task",
			EntityID:   t.ID,
		})
		h.Notifier.Broadcast(ws.EventTaskAssigned, t)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": t.ID, "task": t})
}

// Update — PUT /api/tasks/:id
func (h *TasksHandler) Update(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	var t models.Task
	if err := h.DB.First(&t, id).Error; err != nil {
		JSONError(c, http.StatusNotFound, "task not found")
		return
	}
	var data taskInput
	_ = c.ShouldBindJSON(&data)
	patch := map[string]any{}
	if data.Title != "" {
		patch["title"] = data.Title
	}
	if data.Description != "" {
		patch["description"] = data.Description
	}
	if data.ExecutorID != 0 {
		patch["executor_id"] = data.ExecutorID
	}
	if data.Photo != "" {
		patch["photo"] = data.Photo
	}
	if data.Status != "" {
		patch["status"] = data.Status
	}
	if len(patch) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	if err := h.DB.Model(&t).Updates(patch).Error; err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Delete — DELETE /api/tasks/:id
func (h *TasksHandler) Delete(c *gin.Context) {
	id, ok := ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.Task{}, id).Error; err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
