package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/health")
	{
		api.GET("/", h.checkHealth)
	}
}

// checkHealth godoc
// @Summary Health check
// @Description Check the health of the application
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string "Health status"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /health [get]
func (h *HealthHandler) checkHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "dont worry im okay"})
}
