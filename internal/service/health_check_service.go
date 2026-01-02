package service

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/repository"
)

const (
	StatusOK         = "OK"
	RedisPingTimeout = "redis: client is closed"
)

// healthCheckService implements the HealthCheck interface and provides methods for performing health checks.
type healthCheckService struct {
	serviceName     string
	instanceID      string
	healthCheckRepo repository.HealthCheck
}

// HealthCheck defines the interface for health check service.
// It provides a method to perform health checks and retrieve service status.
//
//go:generate mockery --name HealthCheck --filename health_check_service.go --output ./mocks
type HealthCheck interface {
	// Check performs a health check and returns the status, service name, and instance ID.
	//
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellations
	//
	// Returns:
	//   - string: The health status ("OK")
	//   - string: The name of the service
	//   - string: The unique instance ID of the service
	//   - error: An error if the health check fails, nil otherwise
	Check(ctx context.Context) (string, string, string, error)
}

// NewHealthCheck creates a new instance of HealthCheck service.
//
// Parameters:
//   - serviceName: The name of the service
//   - instanceID: The unique instance ID of the service
//   - healthCheckRepo: The repository used for performing health checks
//
// Returns:
//   - HealthCheck: The initialized HealthCheck service instance
func NewHealthCheck(serviceName, instanceID string, healthCheckRepo repository.HealthCheck) HealthCheck {
	return &healthCheckService{
		serviceName:     serviceName,
		instanceID:      instanceID,
		healthCheckRepo: healthCheckRepo,
	}
}

// Check performs a health check and returns the status, service name, and instance ID.
// It always returns "OK" as the status for a healthy service.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//
// Returns:
//   - string: The health status ("OK")
//   - string: The name of the service
//   - string: The unique instance ID of the service
//   - error: An error if the health check fails, nil otherwise
func (s *healthCheckService) Check(ctx context.Context) (string, string, string, error) {
	if err := s.healthCheckRepo.Ping(ctx); err != nil {
		return RedisPingTimeout, s.serviceName, s.instanceID, err
	}

	return StatusOK, s.serviceName, s.instanceID, nil
}
