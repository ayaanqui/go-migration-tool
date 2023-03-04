package main

import (
	"log"
	"os"

	"github.com/ayaanqui/go-migration-tool/migration_tool"
	"github.com/urfave/cli"
)

func main() {
	config := migration_tool.Config{}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "directory",
				Usage:       "specifies the directory that results are stored",
				Destination: &config.Directory,
				Required:    true,
			},
		},
		Commands: []cli.Command{
			{
				Name:      "create-migration",
				ShortName: "c",
				Usage:     "creates a new migration file",
				Action: func(c *cli.Context) error {
					m := migration_tool.MigrationTool{
						Config: &config,
					}
					m.Config.Directory = migration_tool.StripTrailingSlash(m.Config.Directory)
					migration_name := c.Args().First()
					return m.CreateMigrationFile(migration_name)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
