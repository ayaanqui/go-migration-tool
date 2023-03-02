package main

import (
	"database/sql"
	"time"
)

type MigrationTool struct {
	DbConn *sql.DB
	Config Config
}

type Config struct {
	Directory string // the path to the migration directory. Ex: ./src/migration
	TableName string // defaults to "gorm_migrations"
}

type GormMigrationTable struct {
	Id string
	Name string
	MigrationDate time.Time
}

type ParsedFileName struct {
	Id string
	MigrationName string
	FileExtension string
	Raw string
}