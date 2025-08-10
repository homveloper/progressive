package infrastructure

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// Migration represents a database migration
type Migration struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	AppliedAt   time.Time `db:"applied_at"`
}

// RunMigrations executes all database migrations
func RunMigrations(db *sqlx.DB) error {
	log.Println("🔄 Running database migrations...")

	// Create migrations table if not exists
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Run each migration
	migrations := getMigrations()
	for _, migration := range migrations {
		if err := runMigration(db, migration.name, migration.query); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.name, err)
		}
	}

	log.Println("✅ All migrations completed successfully")
	return nil
}

func createMigrationsTable(db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	return err
}

func runMigration(db *sqlx.DB, name string, query string) error {
	// Check if migration already applied
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM migrations WHERE name = $1", name)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("⏭️  Migration '%s' already applied, skipping", name)
		return nil
	}

	// Run migration in transaction
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration query
	if _, err := tx.Exec(query); err != nil {
		return err
	}

	// Record migration
	if _, err := tx.Exec("INSERT INTO migrations (name) VALUES ($1)", name); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("✅ Applied migration: %s", name)
	return nil
}

type migrationDef struct {
	name  string
	query string
}

func getMigrations() []migrationDef {
	return []migrationDef{
		{
			name: "001_create_tables",
			query: `
				-- Tables table stores JSON Schema-based table definitions
				CREATE TABLE IF NOT EXISTS tables (
					id VARCHAR(255) PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					description TEXT,
					schema JSONB NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMP NOT NULL DEFAULT NOW()
				);

				-- Create index on name for faster lookups
				CREATE INDEX IF NOT EXISTS idx_tables_name ON tables(name);

				-- Records table stores the actual data for each table
				CREATE TABLE IF NOT EXISTS records (
					id SERIAL PRIMARY KEY,
					table_id VARCHAR(255) NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
					data JSONB NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMP NOT NULL DEFAULT NOW()
				);

				-- Create index on table_id for faster queries
				CREATE INDEX IF NOT EXISTS idx_records_table_id ON records(table_id);
				
				-- Create GIN index on data for efficient JSONB queries
				CREATE INDEX IF NOT EXISTS idx_records_data_gin ON records USING GIN (data);
			`,
		},
		{
			name: "002_add_templates",
			query: `
				-- Templates table stores reusable table templates
				CREATE TABLE IF NOT EXISTS templates (
					id VARCHAR(255) PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					description TEXT,
					category VARCHAR(100) NOT NULL,
					icon VARCHAR(100),
					schema JSONB NOT NULL,
					sample_data JSONB,
					created_at TIMESTAMP NOT NULL DEFAULT NOW()
				);

				-- Create index on category for filtering
				CREATE INDEX IF NOT EXISTS idx_templates_category ON templates(category);
			`,
		},
		{
			name: "003_add_metadata",
			query: `
				-- Add metadata columns to tables
				ALTER TABLE tables 
				ADD COLUMN IF NOT EXISTS record_count INTEGER DEFAULT 0,
				ADD COLUMN IF NOT EXISTS last_accessed TIMESTAMP,
				ADD COLUMN IF NOT EXISTS tags TEXT[];

				-- Create index on tags for searching
				CREATE INDEX IF NOT EXISTS idx_tables_tags_gin ON tables USING GIN (tags);

				-- Create function to update record count
				CREATE OR REPLACE FUNCTION update_table_record_count() 
				RETURNS TRIGGER AS $$
				BEGIN
					IF TG_OP = 'INSERT' THEN
						UPDATE tables 
						SET record_count = record_count + 1,
						    last_accessed = NOW()
						WHERE id = NEW.table_id;
					ELSIF TG_OP = 'DELETE' THEN
						UPDATE tables 
						SET record_count = record_count - 1,
						    last_accessed = NOW()
						WHERE id = OLD.table_id;
					END IF;
					RETURN NULL;
				END;
				$$ LANGUAGE plpgsql;

				-- Create triggers for record count
				DROP TRIGGER IF EXISTS update_record_count_on_insert ON records;
				CREATE TRIGGER update_record_count_on_insert
					AFTER INSERT ON records
					FOR EACH ROW
					EXECUTE FUNCTION update_table_record_count();

				DROP TRIGGER IF EXISTS update_record_count_on_delete ON records;
				CREATE TRIGGER update_record_count_on_delete
					AFTER DELETE ON records
					FOR EACH ROW
					EXECUTE FUNCTION update_table_record_count();

				-- Update existing record counts
				UPDATE tables t
				SET record_count = (
					SELECT COUNT(*) 
					FROM records r 
					WHERE r.table_id = t.id
				);
			`,
		},
	}
}