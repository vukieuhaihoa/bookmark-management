package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
)

// healthCheckResponse represents the JSON structure returned by the health check endpoint.
type healthCheckResponse struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

// healthCheckHandler handles health check HTTP requests.
// It uses the HealthCheck service to retrieve health status information.
type healthCheckHandler struct {
	svc service.HealthCheck
}

// HealthCheck defines the interface for health check HTTP handlers.
// It provides a method for handling health check requests using the Gin framework.
type HealthCheck interface {
	Check(ctx *gin.Context)
}

// NewHealthCheck creates a new instance of the HealthCheck handler.
// It accepts a health check service implementation and returns a handler
// that can process HTTP requests for health checks.
//
// Parameters:
//   - svc: The health check service used for retrieving health status
//
// Returns:h
//   - HealthCheck: A new health check handler instance
func NewHealthCheck(svc service.HealthCheck) HealthCheck {
	return &healthCheckHandler{svc: svc}
}

// Check is a Gin framework handler that performs a health check.
// It calls the health check service to retrieve the current health status
// and returns it as a JSON response.
//
// Parameters:
//   - c: The Gin context containing the HTTP request and response
//
// Response:h
//   - 200 OK: Returns the health status as a JSON object
func (h *healthCheckHandler) Check(c *gin.Context) {
	message, serviceName, instanceID := h.svc.Check()
	c.JSON(http.StatusOK, healthCheckResponse{
		Message:     message,
		ServiceName: serviceName,
		InstanceID:  instanceID,
	})
}
