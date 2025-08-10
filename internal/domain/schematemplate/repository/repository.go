package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"progressive/internal/domain/schematemplate"

	"github.com/jmoiron/sqlx"
)

// ReadRepository defines read operations for schema templates
type ReadRepository interface {
	FindAll(ctx context.Context) ([]*schematemplate.SchemaTemplate, error)
	FindByID(ctx context.Context, id string) (*schematemplate.SchemaTemplate, error)
	FindByCategory(ctx context.Context, category string) ([]*schematemplate.SchemaTemplate, error)
	Exists(ctx context.Context, id string) (bool, error)
}

// WriteRepository defines write operations for schema templates
type WriteRepository interface {
	Create(ctx context.Context, template *schematemplate.SchemaTemplate) error
	Update(ctx context.Context, template *schematemplate.SchemaTemplate) error
	Delete(ctx context.Context, id string) error
}

// SchemaTemplateRepository combines both ReadRepository and WriteRepository interfaces
type SchemaTemplateRepository interface {
	ReadRepository
	WriteRepository
}

// schemaTemplateDB represents the database model
type schemaTemplateDB struct {
	ID          string          `db:"id"`
	Name        string          `db:"name"`
	Description string          `db:"description"`
	Category    string          `db:"category"`
	Icon        string          `db:"icon"`
	Schema      json.RawMessage `db:"schema"`
	SampleData  json.RawMessage `db:"sample_data"`
	CreatedAt   time.Time       `db:"created_at"`
}

// PostgresSchemaTemplateRepository implements Repository using PostgreSQL
type PostgresSchemaTemplateRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sqlx.DB) *PostgresSchemaTemplateRepository {
	return &PostgresSchemaTemplateRepository{db: db}
}

// NewPostgresRepositoryWithDefaults creates a new PostgreSQL repository and initializes default templates
func NewPostgresRepositoryWithDefaults(ctx context.Context, db *sqlx.DB) (*PostgresSchemaTemplateRepository, error) {
	repo := &PostgresSchemaTemplateRepository{db: db}
	
	// Initialize default templates
	if err := repo.initializeDefaultTemplates(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize default templates: %w", err)
	}
	
	return repo, nil
}

// initializeDefaultTemplates initializes default templates using type-safe definitions
func (r *PostgresSchemaTemplateRepository) initializeDefaultTemplates(ctx context.Context) error {
	templateDefs := GetDefaultTemplateDefinitions()
	
	for _, def := range templateDefs {
		// Check if template already exists
		exists, err := r.Exists(ctx, def.ID)
		if err != nil {
			return fmt.Errorf("failed to check template existence for %s: %w", def.ID, err)
		}
		
		if exists {
			continue // Skip if already exists
		}
		
		// Convert to domain model
		template, err := def.toDomainModel()
		if err != nil {
			return fmt.Errorf("failed to convert template definition %s: %w", def.ID, err)
		}
		
		// Create template
		if err := r.Create(ctx, template); err != nil {
			return fmt.Errorf("failed to create template %s: %w", def.ID, err)
		}
	}
	
	return nil
}

// FindAll retrieves all schema templates
func (r *PostgresSchemaTemplateRepository) FindAll(ctx context.Context) ([]*schematemplate.SchemaTemplate, error) {
	query := `
		SELECT id, name, description, category, icon, schema, sample_data, created_at 
		FROM templates 
		ORDER BY category, name
	`

	var templates []schemaTemplateDB
	if err := r.db.SelectContext(ctx, &templates, query); err != nil {
		return nil, fmt.Errorf("failed to find all templates: %w", err)
	}

	return r.toDomainList(templates), nil
}

// FindByID retrieves a schema template by ID
func (r *PostgresSchemaTemplateRepository) FindByID(ctx context.Context, id string) (*schematemplate.SchemaTemplate, error) {
	query := `
		SELECT id, name, description, category, icon, schema, sample_data, created_at 
		FROM templates 
		WHERE id = $1
	`

	var template schemaTemplateDB
	if err := r.db.GetContext(ctx, &template, query, id); err != nil {
		return nil, fmt.Errorf("failed to find template by id %s: %w", id, err)
	}

	return r.toDomain(&template), nil
}

// FindByCategory retrieves schema templates by category
func (r *PostgresSchemaTemplateRepository) FindByCategory(ctx context.Context, category string) ([]*schematemplate.SchemaTemplate, error) {
	query := `
		SELECT id, name, description, category, icon, schema, sample_data, created_at 
		FROM templates 
		WHERE category = $1
		ORDER BY name
	`

	var templates []schemaTemplateDB
	if err := r.db.SelectContext(ctx, &templates, query, category); err != nil {
		return nil, fmt.Errorf("failed to find templates by category %s: %w", category, err)
	}

	return r.toDomainList(templates), nil
}

// Exists checks if a template exists by ID
func (r *PostgresSchemaTemplateRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := `SELECT COUNT(*) FROM templates WHERE id = $1`

	var count int
	if err := r.db.GetContext(ctx, &count, query, id); err != nil {
		return false, fmt.Errorf("failed to check template existence: %w", err)
	}

	return count > 0, nil
}

// Create creates a new schema template
func (r *PostgresSchemaTemplateRepository) Create(ctx context.Context, template *schematemplate.SchemaTemplate) error {
	query := `
		INSERT INTO templates (id, name, description, category, icon, schema, sample_data, created_at)
		VALUES (:id, :name, :description, :category, :icon, :schema, :sample_data, :created_at)
	`

	dbModel := r.toDBModel(template)
	if _, err := r.db.NamedExecContext(ctx, query, dbModel); err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	return nil
}

// Update updates an existing schema template
func (r *PostgresSchemaTemplateRepository) Update(ctx context.Context, template *schematemplate.SchemaTemplate) error {
	query := `
		UPDATE templates 
		SET name = :name,
		    description = :description,
		    category = :category,
		    icon = :icon,
		    schema = :schema,
		    sample_data = :sample_data
		WHERE id = :id
	`

	dbModel := r.toDBModel(template)
	result, err := r.db.NamedExecContext(ctx, query, dbModel)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("template with id %s not found", template.ID)
	}

	return nil
}

// Delete deletes a schema template by ID
func (r *PostgresSchemaTemplateRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM templates WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("template with id %s not found", id)
	}

	return nil
}

// Helper methods for conversion between domain and database models

func (r *PostgresSchemaTemplateRepository) toDomain(db *schemaTemplateDB) *schematemplate.SchemaTemplate {
	return &schematemplate.SchemaTemplate{
		ID:          db.ID,
		Name:        db.Name,
		Description: db.Description,
		Category:    db.Category,
		Icon:        db.Icon,
		Schema:      db.Schema,
		SampleData:  db.SampleData,
		CreatedAt:   db.CreatedAt,
	}
}

func (r *PostgresSchemaTemplateRepository) toDomainList(dbList []schemaTemplateDB) []*schematemplate.SchemaTemplate {
	templates := make([]*schematemplate.SchemaTemplate, len(dbList))
	for i, db := range dbList {
		templates[i] = r.toDomain(&db)
	}
	return templates
}

func (r *PostgresSchemaTemplateRepository) toDBModel(domain *schematemplate.SchemaTemplate) *schemaTemplateDB {
	return &schemaTemplateDB{
		ID:          domain.ID,
		Name:        domain.Name,
		Description: domain.Description,
		Category:    domain.Category,
		Icon:        domain.Icon,
		Schema:      domain.Schema,
		SampleData:  domain.SampleData,
		CreatedAt:   domain.CreatedAt,
	}
}
