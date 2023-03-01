package main

import (
	"database/sql"
	"time"
)

type MigrationTool struct {
	DbConn *sql.DB
	Config MigrationConfig
}

type MigrationConfig struct {
	MigrationDirectory string // the path to the migration directory. Ex: ./src/migration
	MigrationTable string // defaults to "gorm_migrations"
}

type GormMigration struct {
	Id string
	Name string
	MigrationDate time.Time
}