package crawler

import "database/sql"

func ExecuteMigration3(db *sql.DB) error {
	var query = `
		RENAME TABLE profx.source_rules TO profx.rules;
		ALTER TABLE profx.links
			ADD COLUMN source   VARCHAR(64) AFTER url,
			ADD COLUMN from_url VARCHAR(1000) AFTER source;
		ALTER TABLE profx.resources
			ADD COLUMN source   VARCHAR(64) AFTER kind,
			ADD COLUMN from_url VARCHAR(1000) AFTER source;

		CREATE INDEX idx_links_source ON profx.links (source);
		CREATE INDEX idx_resources_source ON profx.resources (source);
	`
	_, err := db.Exec(query)

	return err
}

func RevertMigration3(db *sql.DB) error {
	var query = `
		ALTER TABLE profx.links
			DROP COLUMN source,
			DROP COLUMN from_url;
		ALTER TABLE profx.resources
			DROP COLUMN source,
			DROP COLUMN from_url;
		RENAME TABLE profx.rules TO profx.source_rules;
	`
	_, err := db.Exec(query)

	return err
}
