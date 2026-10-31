package handlers

import (
	"net/http"
	"strconv"

	"recur_backend/internals/repository"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	repo *repository.DashboardRepository
}

func NewDashboardHandler(repo *repository.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{repo: repo}
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	months := 12
	if raw := c.Query("months"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 24 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "months must be between 1 and 24"})
			return
		}
		months = n
	}

	resp, err := h.repo.GetDashboard(c.Request.Context(), months)
	if err != nil {
		respondDBError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *DashboardHandler) RegisterRoutes(admin *gin.RouterGroup) {
	admin.GET("/dashboard", h.GetDashboard)
}
