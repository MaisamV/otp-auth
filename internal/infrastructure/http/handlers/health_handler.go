package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}

// Health handles the health check request
// @Summary Health Check
// @Description Check the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Version:   "1.0.0", // TODO: Get from build info or config
		Services: map[string]string{
			"database": "unknown", // TODO: Check database connection
			"redis":    "unknown", // TODO: Check Redis connection
		},
	}

	c.JSON(http.StatusOK, response)
}

// Ready handles the readiness check request
// @Summary Readiness Check
// @Description Check if the service is ready to serve requests
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	// TODO: Check if all dependencies are ready (database, Redis, etc.)
	// For now, assume ready if service is running
	uptime := time.Since(h.startTime)

	response := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Version:   "1.0.0",
		Services: map[string]string{
			"database": "ready", // TODO: Actual health check
			"redis":    "ready", // TODO: Actual health check
		},
	}

	c.JSON(http.StatusOK, response)
}

// Live handles the liveness check request
// @Summary Liveness Check
// @Description Check if the service is alive
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := HealthResponse{
		Status:    "alive",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Version:   "1.0.0",
		Services:  map[string]string{},
	}

	c.JSON(http.StatusOK, response)
}