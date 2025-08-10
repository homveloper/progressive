package infrastructure

import (
	"fmt"
	"net"
	"strconv"
)

// PortFinder is an interface for finding available ports
type PortFinder interface {
	FindAvailablePort(startPort uint32, maxAttempts int) (uint32, error)
}

// DefaultPortFinder implements PortFinder using net.Listen
type DefaultPortFinder struct{}

// NewDefaultPortFinder creates a new DefaultPortFinder
func NewDefaultPortFinder() *DefaultPortFinder {
	return &DefaultPortFinder{}
}

// FindAvailablePort finds an available port starting from startPort
func (pf *DefaultPortFinder) FindAvailablePort(startPort uint32, maxAttempts int) (uint32, error) {
	for i := 0; i < maxAttempts; i++ {
		port := startPort + uint32(i)
		if pf.isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found in range %d-%d", startPort, startPort+uint32(maxAttempts)-1)
}

// isPortAvailable checks if a port is available
func (pf *DefaultPortFinder) isPortAvailable(port uint32) bool {
	address := net.JoinHostPort("localhost", strconv.FormatUint(uint64(port), 10))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

// MockPortFinder for testing purposes
type MockPortFinder struct {
	AvailablePort uint32
	ShouldError   bool
	ErrorMessage  string
}

// NewMockPortFinder creates a new MockPortFinder
func NewMockPortFinder(availablePort uint32) *MockPortFinder {
	return &MockPortFinder{
		AvailablePort: availablePort,
		ShouldError:   false,
	}
}

// FindAvailablePort returns the pre-configured available port
func (mpf *MockPortFinder) FindAvailablePort(startPort uint32, maxAttempts int) (uint32, error) {
	if mpf.ShouldError {
		return 0, fmt.Errorf(mpf.ErrorMessage)
	}
	return mpf.AvailablePort, nil
}

// SetError configures the mock to return an error
func (mpf *MockPortFinder) SetError(errorMsg string) {
	mpf.ShouldError = true
	mpf.ErrorMessage = errorMsg
}
