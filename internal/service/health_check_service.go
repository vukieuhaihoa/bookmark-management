package service

const (
	StatusOK = "OK"
)

// healthCheckService implements the HealthCheck interface and provides methods for performing health checks.
type healthCheckService struct {
	serviceName string
	instanceID  string
}

// HealthCheck defines the interface for health check service.
// It provides a method to perform health checks and retrieve service status.
//
//go:generate mockery --name HealthCheck --filename health_check_service.go --output ./mocks
type HealthCheck interface {
	Check() (string, string, string)
}

// NewHealthCheck creates a new instance of HealthCheck service.
//
// Parameters:
//   - serviceName: The name of the service
//   - instanceID: The unique instance ID of the service
//
// Returns:
//   - HealthCheck: The initialized HealthCheck service instance
func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

// Check performs a health check and returns the status, service name, and instance ID.
// It always returns "OK" as the status for a healthy service.
//
// Returns:
//   - string: The health status ("OK")
//   - string: The name of the service
//   - string: The unique instance ID of the service
func (s *healthCheckService) Check() (string, string, string) {
	return StatusOK, s.serviceName, s.instanceID
}
