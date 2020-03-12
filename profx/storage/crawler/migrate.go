package crawler

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Migrate(cfg Config) error {
	db, err := sql.Open(
		"mysql",
		CreateConnectionString(cfg.Username, cfg.Password, cfg.WriteHost, cfg.WritePort, cfg.Schema),
	)

	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	err = SetupDatabase(db)
	if err != nil {
		return err
	}

	err = SetupMigrations(db)
	if err != nil {
		return err
	}

	migration, err := GetMigration(db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Migration 1 without data
	if *migration < 1 {
		mn := 1
		err = ExecuteMigration1(db)
		if err != nil {
			log.Printf("migration %d failed: %v\n", mn, err)
			log.Printf("reverting...\n")

			mErr := err
			err = RevertMigration1(db)
			if err != nil {
				errS := fmt.Sprintf("revert of migration %d failed: %v\n", mn, err)
				log.Printf(errS)
				return fmt.Errorf(errS)
			}
			return fmt.Errorf("failed to migrate, reverted; error during migration: %v", mErr)
		} else {
			err = SetMigration(db, mn)
			if err != nil {
				return fmt.Errorf("failed to set migration %d number: %v\n", mn, err)
			}
		}
	}

	// Migration 2 with data
	if *migration < 2 {
		mn := 2
		err = ExecuteMigration2(db)
		if err != nil {
			log.Printf("migration %d failed: %v\n", mn, err)
			log.Printf("reverting...\n")

			mErr := err
			err = RevertMigration2(db)
			if err != nil {
				errS := fmt.Sprintf("revert of migration %d failed: %v\n", mn, err)
				log.Printf(errS)
				return fmt.Errorf(errS)
			}
			return fmt.Errorf("failed to migrate, reverted; error during migration: %v", mErr)
		} else {
			err = DataMigration2(db)
			if err != nil {
				errS := fmt.Sprintf("importing data for migration %d failed: %v\n", mn, err)
				log.Printf(errS)
				return fmt.Errorf(errS)
			}

			err = SetMigration(db, mn)
			if err != nil {
				return fmt.Errorf("failed to set migration %d number: %v\n", mn, err)
			}
		}
	}

	// Migration 3 without data
	if *migration < 3 {
		mn := 3
		err = ExecuteMigration3(db)
		if err != nil {
			log.Printf("migration %d failed: %v\n", mn, err)
			log.Printf("reverting...\n")

			mErr := err
			err = RevertMigration3(db)
			if err != nil {
				errS := fmt.Sprintf("revert of migration %d failed: %v\n", mn, err)
				log.Printf(errS)
				return fmt.Errorf(errS)
			}
			return fmt.Errorf("failed to migrate, reverted; error during migration: %v", mErr)
		} else {
			err = SetMigration(db, mn)
			if err != nil {
				return fmt.Errorf("failed to set migration %d number: %v\n", mn, err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func SetupDatabase(db *sql.DB) error {
	var query = `
		CREATE DATABASE IF NOT EXISTS profx CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
	`
	_, err := db.Exec(query)

	return err
}

func SetupMigrations(db *sql.DB) (err error) {
	var query = `
		CREATE TABLE IF NOT EXISTS profx.migrations
		(
			id SMALLINT UNSIGNED PRIMARY KEY,
			migration SMALLINT UNSIGNED NOT NULL
		);
		INSERT IGNORE INTO profx.migrations (id, migration) VALUES (1, 0);
	`
	_, err = db.Exec(query)
	return
}

func GetMigration(db *sql.DB) (*int, error) {
	var res int
	var query = `SELECT migration FROM profx.migrations WHERE id = 1`
	err := db.QueryRow(query).Scan(&res)
	if err != nil {
		return &res, fmt.Errorf("can't get migration: %w", err)
	}

	return &res, nil
}

func SetMigration(db *sql.DB, m int) error {
	var query = `
		UPDATE profx.migrations
		SET migration = ?
		WHERE id = 1
	`
	res, err := db.Exec(query, m)

	if err != nil {
		return fmt.Errorf("can't set migration: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("can't set migration: %w", err)
	}

	return nil
}
