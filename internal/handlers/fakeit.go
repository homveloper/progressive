package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"progressive/internal/pages"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// FakeitPageHandler renders the fakeit page
func (h *Handlers) FakeitPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := pages.FakeitPage().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FakeitGenerateRequest represents the request for generating fake data
type FakeitGenerateRequest struct {
	Schema       map[string]interface{}     `json:"schema"`
	FieldConfigs map[string]FakeFieldConfig `json:"fieldConfigs"`
	Count        int                        `json:"count"`
}

// FakeFieldConfig represents configuration for a single field
type FakeFieldConfig struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}

// FakeitGenerateAPIHandler handles fake data generation
func (h *Handlers) FakeitGenerateAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FakeitGenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Schema == nil {
		http.Error(w, "Schema is required", http.StatusBadRequest)
		return
	}

	if req.Count <= 0 || req.Count > 10000 {
		http.Error(w, "Count must be between 1 and 10000", http.StatusBadRequest)
		return
	}

	// Generate fake data
	data, err := h.generateFakeData(req.Schema, req.FieldConfigs, req.Count)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate data: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]interface{}{
		"success": true,
		"data":    data,
		"count":   len(data),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateFakeData generates fake data based on schema and field configurations
func (h *Handlers) generateFakeData(schema map[string]interface{}, fieldConfigs map[string]FakeFieldConfig, count int) ([]map[string]interface{}, error) {
	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid schema: properties not found")
	}

	gofakeit.Seed(time.Now().UnixNano())

	var results []map[string]interface{}

	// Extract field names in order from schema
	fieldNames := make([]string, 0, len(properties))
	for fieldName := range properties {
		fieldNames = append(fieldNames, fieldName)
	}

	for i := 0; i < count; i++ {
		record := make(map[string]interface{})

		// Process fields in the order they appear in the schema
		for _, fieldName := range fieldNames {
			fieldSchema := properties[fieldName]
			fieldSchemaMap, ok := fieldSchema.(map[string]interface{})
			if !ok {
				continue
			}

			fieldType, _ := fieldSchemaMap["type"].(string)
			if fieldType == "" {
				fieldType = "string"
			}

			// Get field configuration
			config, exists := fieldConfigs[fieldName]
			if !exists {
				// Use default configuration
				config = FakeFieldConfig{Type: "name", Params: make(map[string]string)}
			}

			value, err := h.generateFieldValue(fieldType, config)
			if err != nil {
				continue // Skip fields that can't be generated
			}

			record[fieldName] = value
		}

		results = append(results, record)
	}

	return results, nil
}

// generateFieldValue generates a fake value for a specific field
func (h *Handlers) generateFieldValue(fieldType string, config FakeFieldConfig) (interface{}, error) {
	switch fieldType {
	case "string":
		return h.generateStringValue(config)
	case "integer":
		return h.generateIntegerValue(config)
	case "number":
		return h.generateNumberValue(config)
	case "boolean":
		return h.generateBooleanValue(config)
	default:
		return h.generateStringValue(config)
	}
}

// generateStringValue generates fake string values
func (h *Handlers) generateStringValue(config FakeFieldConfig) (string, error) {
	switch config.Type {
	case "name":
		return gofakeit.Name(), nil
	case "firstName":
		return gofakeit.FirstName(), nil
	case "lastName":
		return gofakeit.LastName(), nil
	case "email":
		return gofakeit.Email(), nil
	case "phone":
		return gofakeit.Phone(), nil
	case "address":
		return gofakeit.Address().Address, nil
	case "company":
		return gofakeit.Company(), nil
	case "jobTitle":
		return gofakeit.JobTitle(), nil
	case "city":
		return gofakeit.City(), nil
	case "country":
		return gofakeit.Country(), nil
	case "lorem":
		wordCount := 5
		if wc, exists := config.Params["wordCount"]; exists {
			if parsed, err := strconv.Atoi(wc); err == nil && parsed > 0 {
				wordCount = parsed
			}
		}
		return gofakeit.LoremIpsumSentence(wordCount), nil
	case "sentence":
		return gofakeit.Sentence(rand.Intn(10) + 5), nil
	case "paragraph":
		return gofakeit.Paragraph(1, 3, rand.Intn(8)+5, " "), nil
	case "uuid":
		return gofakeit.UUID(), nil
	case "url":
		return gofakeit.URL(), nil
	case "username":
		return gofakeit.Username(), nil
	case "password":
		return gofakeit.Password(true, true, true, true, false, rand.Intn(8)+8), nil
	case "color":
		return gofakeit.HexColor(), nil
	case "custom":
		if values, exists := config.Params["values"]; exists {
			valueList := strings.Split(values, ",")
			if len(valueList) > 0 {
				trimmed := make([]string, len(valueList))
				for i, v := range valueList {
					trimmed[i] = strings.TrimSpace(v)
				}
				return gofakeit.RandomString(trimmed), nil
			}
		}
		return gofakeit.Word(), nil
	default:
		return gofakeit.Word(), nil
	}
}

// generateIntegerValue generates fake integer values
func (h *Handlers) generateIntegerValue(config FakeFieldConfig) (int, error) {
	switch config.Type {
	case "age":
		return gofakeit.Number(18, 80), nil
	case "year":
		return gofakeit.Year(), nil
	case "month":
		return gofakeit.Month(), nil
	case "day":
		return gofakeit.Day(), nil
	case "price":
		return gofakeit.Number(1000, 100000), nil
	case "quantity":
		return gofakeit.Number(1, 100), nil
	case "rating":
		return gofakeit.Number(1, 5), nil
	case "custom":
		min := 1
		max := 100
		if minStr, exists := config.Params["min"]; exists {
			if parsed, err := strconv.Atoi(minStr); err == nil {
				min = parsed
			}
		}
		if maxStr, exists := config.Params["max"]; exists {
			if parsed, err := strconv.Atoi(maxStr); err == nil {
				max = parsed
			}
		}
		return gofakeit.Number(min, max), nil
	default:
		return gofakeit.Number(1, 1000), nil
	}
}

// generateNumberValue generates fake float values
func (h *Handlers) generateNumberValue(config FakeFieldConfig) (float64, error) {
	switch config.Type {
	case "price":
		return gofakeit.Price(10.0, 1000.0), nil
	case "latitude":
		return gofakeit.Latitude(), nil
	case "longitude":
		return gofakeit.Longitude(), nil
	case "percentage":
		return gofakeit.Float64Range(0.0, 100.0), nil
	case "custom":
		min := 0.0
		max := 100.0
		if minStr, exists := config.Params["min"]; exists {
			if parsed, err := strconv.ParseFloat(minStr, 64); err == nil {
				min = parsed
			}
		}
		if maxStr, exists := config.Params["max"]; exists {
			if parsed, err := strconv.ParseFloat(maxStr, 64); err == nil {
				max = parsed
			}
		}
		return gofakeit.Float64Range(min, max), nil
	default:
		return gofakeit.Float64Range(0.0, 100.0), nil
	}
}

// generateBooleanValue generates fake boolean values
func (h *Handlers) generateBooleanValue(config FakeFieldConfig) (bool, error) {
	switch config.Type {
	case "weighted":
		probability := 50 // default 50%
		if probStr, exists := config.Params["trueProbability"]; exists {
			if parsed, err := strconv.Atoi(probStr); err == nil && parsed >= 0 && parsed <= 100 {
				probability = parsed
			}
		}
		return rand.Intn(100) < probability, nil
	default:
		return gofakeit.Bool(), nil
	}
}
