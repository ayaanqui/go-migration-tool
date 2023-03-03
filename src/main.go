package main

import (
	"log"
	"os"

	"github.com/ayaanqui/go-migration-tool/src/tool"
	"github.com/urfave/cli"
)

func main() {
	config := tool.Config{}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "directory",
				Usage: "specifies the directory that results are stored",
				Destination: &config.Directory,
				Required: true,
			},
		},
		Commands: []cli.Command{
			{
				Name: "create-migration",
				ShortName: "c",
				Usage: "creates a new migration file",
				Action: func(c *cli.Context) error {
					migration_tool := tool.MigrationTool{
						Config: &config,
					}
					migration_tool.Config.Directory = tool.StringTrailingSlash(migration_tool.Config.Directory)
                    migration_name := c.Args().First()
					return migration_tool.CreateMigrationFile(migration_name)
                },
			},
		},
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}