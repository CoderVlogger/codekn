package crawler

import "database/sql"

func ExecuteMigration1(db *sql.DB) error {
	var query = `
		CREATE TABLE IF NOT EXISTS profx.sources
		(
			url     VARCHAR(1000)                             NOT NULL,
			rule    VARCHAR(1000)                             NOT NULL,
			created TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS profx.links
		(
			hash    VARCHAR(64) PRIMARY KEY                   NOT NULL,
			url     VARCHAR(1000)                             NOT NULL,
			created TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS profx.resources
		(
			hash        VARCHAR(64) PRIMARY KEY                   NOT NULL,
			kind        VARCHAR(64)                               NOT NULL,
			url         VARCHAR(1000)                             NOT NULL,
			title       VARCHAR(1000),
			description VARCHAR(1000),
			created     TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
		);
		CREATE INDEX idx_resources_kind ON profx.resources (kind);
		CREATE TABLE IF NOT EXISTS profx.sys_logs
		(
			message VARCHAR(1000)                             NOT NULL,
			created TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
		);
	`
	_, err := db.Exec(query)

	return err
}

func RevertMigration1(db *sql.DB) error {
	var query = `
		DROP TABLE IF EXISTS profx.sources;
		DROP TABLE IF EXISTS profx.links;
		DROP TABLE IF EXISTS profx.resources;
		DROP TABLE IF EXISTS profx.sys_logs;
	`
	_, err := db.Exec(query)

	return err
}
