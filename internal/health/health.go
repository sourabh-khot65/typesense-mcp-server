// Package health provides health check functionality.
package health

import (
	"context"
	"fmt"
	"time"

	"typesense-mcp-server/internal/constants"
	"typesense-mcp-server/internal/errors"
	"typesense-mcp-server/internal/logger"
	"typesense-mcp-server/internal/services"
)

// Status represents the health status
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// Check represents a single health check
type Check struct {
	Name     string        `json:"name"`
	Status   Status        `json:"status"`
	Duration time.Duration `json:"duration_ms"`
	Error    string        `json:"error,omitempty"`
	Details  interface{}   `json:"details,omitempty"`
}

// Result represents the overall health check result
type Result struct {
	Status    Status           `json:"status"`
	Timestamp time.Time        `json:"timestamp"`
	Duration  time.Duration    `json:"duration_ms"`
	Checks    map[string]Check `json:"checks"`
	Version   string           `json:"version"`
}

// Checker provides health check functionality
type Checker struct {
	typesenseService services.TypesenseService
	logger           *logger.Logger
	version          string
}

// New creates a new health checker
func New(typesenseService services.TypesenseService, version string, log *logger.Logger) *Checker {
	if log == nil {
		log = logger.NewDefault()
	}

	return &Checker{
		typesenseService: typesenseService,
		logger:           log,
		version:          version,
	}
}

// Check performs all health checks and returns the result
func (c *Checker) Check(ctx context.Context) *Result {
	start := time.Now()
	
	entry := c.logger.WithComponent(constants.ComponentHealth).
		WithField(constants.LogFieldOperation, "health_check")

	result := &Result{
		Timestamp: start,
		Checks:    make(map[string]Check),
		Version:   c.version,
	}

	// Perform individual checks
	c.checkTypesense(ctx, result)
	c.checkBasic(ctx, result)

	// Determine overall status
	result.Status = c.determineOverallStatus(result.Checks)
	result.Duration = time.Since(start)

	entry.WithField(constants.LogFieldDuration, result.Duration.Milliseconds()).
		WithField("overall_status", result.Status).
		WithField("check_count", len(result.Checks)).
		Info("health check completed")

	return result
}

// checkTypesense performs Typesense service health check
func (c *Checker) checkTypesense(ctx context.Context, result *Result) {
	start := time.Now()
	check := Check{
		Name: "typesense",
	}

	if c.typesenseService == nil {
		check.Status = StatusUnhealthy
		check.Error = "typesense service not initialized"
	} else {
		err := c.typesenseService.Health(ctx)
		if err != nil {
			check.Status = StatusUnhealthy
			check.Error = err.Error()
			
			c.logger.WithComponent(constants.ComponentHealth).
				WithError(err).
				Error("typesense health check failed")
		} else {
			check.Status = StatusHealthy
			check.Details = map[string]interface{}{
				"service": "available",
			}
		}
	}

	check.Duration = time.Since(start)
	result.Checks["typesense"] = check
}

// checkBasic performs basic application health checks
func (c *Checker) checkBasic(ctx context.Context, result *Result) {
	start := time.Now()
	check := Check{
		Name:   "basic",
		Status: StatusHealthy,
		Details: map[string]interface{}{
			"version":   c.version,
			"timestamp": time.Now().UTC(),
		},
	}

	check.Duration = time.Since(start)
	result.Checks["basic"] = check
}

// determineOverallStatus determines the overall health status based on individual checks
func (c *Checker) determineOverallStatus(checks map[string]Check) Status {
	if len(checks) == 0 {
		return StatusUnhealthy
	}

	hasUnhealthy := false
	hasDegraded := false

	for _, check := range checks {
		switch check.Status {
		case StatusUnhealthy:
			hasUnhealthy = true
		case StatusDegraded:
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return StatusUnhealthy
	}
	if hasDegraded {
		return StatusDegraded
	}
	return StatusHealthy
}

// IsHealthy returns true if the overall status is healthy
func (r *Result) IsHealthy() bool {
	return r.Status == StatusHealthy
}

// GetFailedChecks returns a list of failed health checks
func (r *Result) GetFailedChecks() []Check {
	var failed []Check
	for _, check := range r.Checks {
		if check.Status != StatusHealthy {
			failed = append(failed, check)
		}
	}
	return failed
}

// String returns a string representation of the health result
func (r *Result) String() string {
	failedCount := len(r.GetFailedChecks())
	totalCount := len(r.Checks)
	
	if failedCount == 0 {
		return fmt.Sprintf("Health: %s (%d/%d checks passed)", r.Status, totalCount, totalCount)
	}
	
	return fmt.Sprintf("Health: %s (%d/%d checks failed)", r.Status, failedCount, totalCount)
}

// HealthError creates an error from a health check result
func (r *Result) HealthError() error {
	if r.IsHealthy() {
		return nil
	}

	failed := r.GetFailedChecks()
	if len(failed) == 0 {
		return errors.NewServiceError("health check failed")
	}

	// Create a detailed error with all failed checks
	err := errors.NewServiceError("health check failed").
		WithDetail("overall_status", r.Status).
		WithDetail("failed_checks", len(failed)).
		WithDetail("total_checks", len(r.Checks))

	for i, check := range failed {
		err.WithDetail(fmt.Sprintf("failed_check_%d", i), map[string]interface{}{
			"name":     check.Name,
			"status":   check.Status,
			"error":    check.Error,
			"duration": check.Duration.Milliseconds(),
		})
	}

	return err
}
