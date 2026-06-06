package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	DB *gorm.DB
}

// Stats — GET /api/dashboard/stats
func (h *DashboardHandler) Stats(c *gin.Context) {
	var users, tasks, lids, requests, houses, posts, sims int64
	h.DB.Model(&models.User{}).Count(&users)
	h.DB.Model(&models.Task{}).Count(&tasks)
	h.DB.Model(&models.KanbanLead{}).Count(&lids)
	h.DB.Model(&models.RequestsItem{}).Count(&requests)
	h.DB.Model(&models.House{}).Count(&houses)
	h.DB.Model(&models.Post{}).Count(&posts)
	h.DB.Model(&models.SimCard{}).Count(&sims)
	c.JSON(http.StatusOK, gin.H{
		"users":     users,
		"tasks":     tasks,
		"leads":     lids,
		"requests":  requests,
		"houses":    houses,
		"posts":     posts,
		"sim_cards": sims,
	})
}
