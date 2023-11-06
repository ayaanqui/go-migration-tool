# Go Migration Tool
A simple database migration tool that allows making sequential migrations easy.

## Installation
```
$ go get github.com/ayaanqui/go-migration-tool
```

## Usage
This repository comes with 2 different sets of tooling. The first and probably the most important is the migration-tool. This is module allows the user to create a migration table with a list of all successful migrations. And when new migrations are detected, it will go ahead and create these migrations.

The second module is a simple CLI, which allows the user to create these database migration files. These files are just regular `.sql` files with a specific migration name and a UNIX timestamp.

### migration-tool
```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/ayaanqui/go-migration-tool/migration_tool"
)

func main() {
    db, err := get_db_connection() // Arbitrary function that returns an pointer to sql.DB
    migration := migration_tool.New(db, &migration_tool.Config{
        Directory: "./migrations", // Directory which will contain all migraiton files
        TableName: "migrations", // Name of the table that will hold all successful migrations
    })
    migration.RunMigration()

    // Basic net/http server
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Homepage")
    })
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### CLI
```
go run github.com/ayaanqui/go-migration-tool --directory "./migrations" create-migration MyNewMigration
```
Running this command will create a new file inside the `./migrations` directory with the file name `[timestamp]_MyNewMigration.sql`.

Assuming that the server was setup in a way, such that the `RunMigration()` method is called before starting the server, the method should take the contents of the generated SQL file and execute it, while also creating a new row inside the migration tabel with the UNIX timpstamp and the name of the migration.
