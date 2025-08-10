package schematemplate

import (
	"encoding/json"
	"time"
)

// SchemaTemplate represents a reusable table template
type SchemaTemplate struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Icon        string          `json:"icon"`
	Schema      json.RawMessage `json:"schema"`
	SampleData  json.RawMessage `json:"sample_data,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

// Category constants
const (
	CategoryBusiness = "business"
	CategoryGame     = "game"
	CategoryCustom   = "custom"
)

// NewSchemaTemplate creates a new schema template
func NewSchemaTemplate(id, name, description, category, icon string, schema json.RawMessage) *SchemaTemplate {
	return &SchemaTemplate{
		ID:          id,
		Name:        name,
		Description: description,
		Category:    category,
		Icon:        icon,
		Schema:      schema,
		CreatedAt:   time.Now(),
	}
}

// Validate validates the schema template
func (st *SchemaTemplate) Validate() error {
	// Add validation logic here
	return nil
}