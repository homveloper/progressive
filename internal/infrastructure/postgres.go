package infrastructure

import (
	"fmt"
	"log"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     uint32
	Username string
	Password string
	Database string
}

// Option defines a function type for configuring EmbeddedDB
type Option func(*EmbeddedDBOptions)

// EmbeddedDBOptions holds configuration options for EmbeddedDB
type EmbeddedDBOptions struct {
	Config             Config
	PortFinder         PortFinder
	AutoPortDiscovery  bool
	MaxPortAttempts    int
}

// WithConfig sets the database configuration
func WithConfig(config Config) Option {
	return func(opts *EmbeddedDBOptions) {
		opts.Config = config
	}
}

// WithPortFinder sets a custom port finder
func WithPortFinder(pf PortFinder) Option {
	return func(opts *EmbeddedDBOptions) {
		opts.PortFinder = pf
	}
}

// WithAutoPortDiscovery enables automatic port discovery
func WithAutoPortDiscovery(maxAttempts int) Option {
	return func(opts *EmbeddedDBOptions) {
		opts.AutoPortDiscovery = true
		opts.MaxPortAttempts = maxAttempts
	}
}

// getDefaultOptions returns default options for EmbeddedDB
func getDefaultOptions() *EmbeddedDBOptions {
	return &EmbeddedDBOptions{
		Config: Config{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
			Database: "progressive",
		},
		PortFinder:        NewDefaultPortFinder(),
		AutoPortDiscovery: true,
		MaxPortAttempts:   10,
	}
}

// EmbeddedDB holds embedded PostgreSQL instance and connection
type EmbeddedDB struct {
	DB       *sqlx.DB
	embedded *embeddedpostgres.EmbeddedPostgres
	config   Config
}

// GetConfig returns the configuration used by the EmbeddedDB
func (e *EmbeddedDB) GetConfig() Config {
	return e.config
}

// NewEmbeddedDB creates and starts an embedded PostgreSQL instance with default options
func NewEmbeddedDB(options ...Option) (*EmbeddedDB, error) {
	opts := getDefaultOptions()
	for _, option := range options {
		option(opts)
	}
	
	return NewEmbeddedDBWithOptions(opts)
}

// NewEmbeddedDBWithOptions creates and starts an embedded PostgreSQL with custom options
func NewEmbeddedDBWithOptions(opts *EmbeddedDBOptions) (*EmbeddedDB, error) {
	config := opts.Config
	
	// If auto port discovery is enabled, find an available port
	if opts.AutoPortDiscovery {
		availablePort, err := opts.PortFinder.FindAvailablePort(config.Port, opts.MaxPortAttempts)
		if err != nil {
			return nil, fmt.Errorf("failed to find available port: %w", err)
		}
		config.Port = availablePort
		log.Printf("🔍 Found available port: %d", availablePort)
	}
	
	return createEmbeddedDB(config)
}

// NewEmbeddedDBWithConfig creates and starts an embedded PostgreSQL with custom config (deprecated)
func NewEmbeddedDBWithConfig(config Config) (*EmbeddedDB, error) {
	return createEmbeddedDB(config)
}

// createEmbeddedDB is the internal function that creates the embedded database
func createEmbeddedDB(config Config) (*EmbeddedDB, error) {
	// Configure embedded PostgreSQL
	embeddedConfig := embeddedpostgres.DefaultConfig().
		Username(config.Username).
		Password(config.Password).
		Database(config.Database).
		Port(config.Port).
		StartTimeout(45 * time.Second)

	embedded := embeddedpostgres.NewDatabase(embeddedConfig)

	// Start embedded PostgreSQL
	log.Println("🚀 Starting embedded PostgreSQL...")
	if err := embedded.Start(); err != nil {
		return nil, fmt.Errorf("failed to start embedded postgres: %w", err)
	}
	log.Println("✅ Embedded PostgreSQL started successfully")

	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)

	// Wait a bit for the database to be ready
	time.Sleep(2 * time.Second)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		embedded.Stop()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Database connection established")
	
	return &EmbeddedDB{
		DB:       db,
		embedded: embedded,
		config:   config,
	}, nil
}

// NewDB creates a connection to an external PostgreSQL database
func NewDB(config Config) (*sqlx.DB, error) {
	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Database connection established")
	return db, nil
}

// Close stops the embedded PostgreSQL and closes the connection
func (e *EmbeddedDB) Close() error {
	if e.DB != nil {
		if err := e.DB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}

	if e.embedded != nil {
		log.Println("🛑 Stopping embedded PostgreSQL...")
		if err := e.embedded.Stop(); err != nil {
			return fmt.Errorf("failed to stop embedded postgres: %w", err)
		}
		log.Println("✅ Embedded PostgreSQL stopped")
	}

	return nil
}
