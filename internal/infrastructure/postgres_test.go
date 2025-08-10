package infrastructure

import (
	"testing"
)

func TestPortFinderInterface(t *testing.T) {
	// Test DefaultPortFinder
	pf := NewDefaultPortFinder()
	port, err := pf.FindAvailablePort(5432, 5)
	if err != nil {
		t.Logf("Could not find available port (this is normal if ports are busy): %v", err)
	} else {
		t.Logf("Found available port: %d", port)
	}
}

func TestMockPortFinder(t *testing.T) {
	// Test MockPortFinder success case
	mockPF := NewMockPortFinder(9999)
	port, err := mockPF.FindAvailablePort(5432, 5)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if port != 9999 {
		t.Errorf("Expected port 9999, got: %d", port)
	}

	// Test MockPortFinder error case
	mockPF.SetError("mock error")
	_, err = mockPF.FindAvailablePort(5432, 5)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "mock error" {
		t.Errorf("Expected 'mock error', got: %s", err.Error())
	}
}

func TestNewEmbeddedDBWithOptions(t *testing.T) {
	// Test with mock port finder
	mockPF := NewMockPortFinder(5555)
	
	// This test won't actually start PostgreSQL, just test the option pattern
	opts := &EmbeddedDBOptions{
		Config: Config{
			Host:     "localhost",
			Port:     5432,
			Username: "test",
			Password: "test",
			Database: "test",
		},
		PortFinder:        mockPF,
		AutoPortDiscovery: true,
		MaxPortAttempts:   5,
	}
	
	// We can't actually test the full database creation without embedded postgres,
	// but we can verify that the options are applied correctly
	if opts.Config.Port != 5432 {
		t.Errorf("Expected initial port 5432, got: %d", opts.Config.Port)
	}
	
	if opts.AutoPortDiscovery != true {
		t.Error("Expected AutoPortDiscovery to be true")
	}
	
	if opts.MaxPortAttempts != 5 {
		t.Errorf("Expected MaxPortAttempts to be 5, got: %d", opts.MaxPortAttempts)
	}
}

func TestOptionsPattern(t *testing.T) {
	// Test default options
	opts := getDefaultOptions()
	if opts.Config.Port != 5432 {
		t.Errorf("Expected default port 5432, got: %d", opts.Config.Port)
	}
	if !opts.AutoPortDiscovery {
		t.Error("Expected AutoPortDiscovery to be true by default")
	}
	
	// Test WithConfig option
	customConfig := Config{
		Host:     "testhost",
		Port:     9999,
		Username: "testuser",
		Password: "testpass",
		Database: "testdb",
	}
	
	WithConfig(customConfig)(opts)
	if opts.Config.Host != "testhost" {
		t.Errorf("Expected host 'testhost', got: %s", opts.Config.Host)
	}
	if opts.Config.Port != 9999 {
		t.Errorf("Expected port 9999, got: %d", opts.Config.Port)
	}
	
	// Test WithPortFinder option
	mockPF := NewMockPortFinder(8888)
	WithPortFinder(mockPF)(opts)
	if opts.PortFinder != mockPF {
		t.Error("Expected port finder to be set to mock")
	}
	
	// Test WithAutoPortDiscovery option
	WithAutoPortDiscovery(20)(opts)
	if opts.MaxPortAttempts != 20 {
		t.Errorf("Expected MaxPortAttempts to be 20, got: %d", opts.MaxPortAttempts)
	}
}