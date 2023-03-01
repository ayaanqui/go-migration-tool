package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
)

type MigrationConfig struct {
	MigrationDirectory string // the path to the migration directory. Ex: ./src/migration
	MigrationTable string // defaults to "gorm_migrations"
}

type GormMigration struct {
	Id string
	Name string
	MigrationDate time.Time
}

// creates migration table if it doesn't already exist
func create_migration_table(db_conn *sql.DB, table_name string) {
	_, err := db_conn.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%s" (
			id VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			migration_date TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`, table_name))
	if err != nil {
		panic("could not execute migration table creation query")
	}
}

func New(db_conn *sql.DB, config MigrationConfig) {
	if db_conn == nil {
		panic("database connection is not defined")
	}
	if config.MigrationTable == "" {
		config.MigrationTable = "gorm_migrations"
	}

	create_migration_table(db_conn, config.MigrationTable)

	// retrieve all rows from config.MigrationTable table
	rows, err := db_conn.Query(fmt.Sprintf(`
		SELECT id, name
		FROM "%s";
	`, config.MigrationTable))
	if err != nil {
		panic(fmt.Sprintf("could not select from %s table", config.MigrationTable))
	}

	db_migrations := []GormMigration{}
	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		
		db_migrations = append(db_migrations, GormMigration{
			Id: id,
			Name: name,
		})
	}
	rows.Close()

	// get all migration files from config.MigrationDirectory directory
	migration_files, err := os.ReadDir(config.MigrationDirectory)
	if err != nil {
		panic(err)
	}

	file_migrations := []GormMigration{}
	for _, file := range migration_files {
		file_name := file.Name()
		split_file_name := strings.SplitN(file_name, "_", 2)
		if len(split_file_name) != 2 {
			continue
		}
		raw_id := split_file_name[0]
		raw_migration_name := split_file_name[1]

		file_migrations = append(file_migrations, GormMigration{
			Id: raw_id,
			Name: raw_migration_name,
		})
	}
}