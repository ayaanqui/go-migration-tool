package main

import (
	"database/sql"
	"fmt"
	"os"
)

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

func New(db_conn *sql.DB, config Config) MigrationTool {
	if db_conn == nil {
		panic("database connection is not defined")
	}
	if config.TableName == "" {
		config.TableName = "gorm_migrations"
	}

	create_migration_table(db_conn, config.TableName)
	return MigrationTool{
		DbConn: db_conn,
		Config: config,
	}
}

func (c *MigrationTool) RunMigration() {
	// retrieve all rows from config.MigrationTable table
	rows, err := c.DbConn.Query(fmt.Sprintf(`
		SELECT id, name
		FROM "%s";
	`, c.Config.TableName))
	if err != nil {
		panic(fmt.Sprintf("could not select from %s table", c.Config.TableName))
	}

	db_migrations := []GormMigrationTable{}
	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		
		db_migrations = append(db_migrations, GormMigrationTable{
			Id: id,
			Name: name,
		})
	}
	rows.Close()

	// get all migration files from config.MigrationDirectory directory
	migration_files, err := os.ReadDir(c.Config.Directory)
	if err != nil {
		panic(err)
	}

	file_migrations := []ParsedFileName{}
	for i, file := range migration_files {
		file_name := file.Name()
		parsed_val, err := parse_file_name(file_name)
		if err != nil || parsed_val.FileExtension != ".sql" {
			continue
		}

		id := parsed_val.Id
		if id == db_migrations[i].Id {
			continue
		}
		file_migrations = append(file_migrations, parsed_val)
	}

	for _, val := range file_migrations {
		data, err := os.ReadFile(fmt.Sprintf("%s/%s", c.Config.Directory, val.Raw))
		if err != nil {
			panic(err)
		}

		tx, err := c.DbConn.Begin()
		if err != nil {
			panic(err)
		}
		tx.Exec(string(data))
		tx.Exec(fmt.Sprintf(`
			INSERT INTO %s (id, name) VALUES(%s, %s);
		`, c.Config.TableName, val.Id, val.MigrationName))
		tx.Commit()
	}
}