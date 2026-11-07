package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Initialize(dbDir string) error {
	var dsn string
	if dbDir == ":memory:" {
		dsn = "file:memdb1?mode=memory&cache=shared&_pragma=foreign_keys(1)&_pragma=recursive_triggers(1)"
	} else {
		dbPath := filepath.Join(dbDir, "app.db")
		dsn = fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=recursive_triggers(1)", dbPath)
	}
	
	var err error
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("could not open database: %w", err)
	}

	return createTables()
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS ideas (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			icon_path TEXT,
			is_important BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS memories (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			memory_date DATETIME,
			is_important BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS interlinks (
			id TEXT PRIMARY KEY,
			source_id TEXT NOT NULL,
			source_type TEXT NOT NULL,
			target_id TEXT NOT NULL,
			target_type TEXT NOT NULL,
			name TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS metaData (
			id TEXT PRIMARY KEY,
			parent_id TEXT NOT NULL,
			parent_type TEXT NOT NULL,
			level INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS attachments (
			id TEXT PRIMARY KEY,
			metadata_id TEXT NOT NULL,
			file_path TEXT NOT NULL,
			is_photo BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(metadata_id) REFERENCES metaData(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_metaData_parent ON metaData(parent_id);`,
		`CREATE INDEX IF NOT EXISTS idx_interlinks_source ON interlinks(source_id);`,
		`CREATE INDEX IF NOT EXISTS idx_interlinks_target ON interlinks(target_id);`,

		// Create Triggers for Polymorphic Cascading Deletes
		`CREATE TRIGGER IF NOT EXISTS cascade_idea_delete
		AFTER DELETE ON ideas
		BEGIN
			DELETE FROM metaData WHERE parent_id = old.id AND parent_type = 'Idea';
		END;`,
		`CREATE TRIGGER IF NOT EXISTS cascade_memory_delete
		AFTER DELETE ON memories
		BEGIN
			DELETE FROM metaData WHERE parent_id = old.id AND parent_type = 'Memory';
		END;`,
		`CREATE TRIGGER IF NOT EXISTS cascade_interlink_delete
		AFTER DELETE ON interlinks
		BEGIN
			DELETE FROM metaData WHERE parent_id = old.id AND parent_type = 'InterLink';
		END;`,
		`CREATE TRIGGER IF NOT EXISTS cascade_metadata_delete
		AFTER DELETE ON metaData
		BEGIN
			DELETE FROM metaData WHERE parent_id = old.id AND parent_type = 'MetaData';
		END;`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("error creating table: %w initially executing %s", err, query)
		}
	}
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
