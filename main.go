package main

import (
	"log"
	"os"

	"github.com/chyroc/dropbox-cli/command"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "dropbox-cli",
		Commands: []*cli.Command{
			command.Download(),
			command.Upload(),
			command.SaveURL(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
