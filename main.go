package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/chyroc/dropbox-cli/command"
)

func main() {
	app := &cli.App{
		Name: "dropbox-cli",
		Commands: []*cli.Command{
			command.Download(),
			command.Upload(),
			command.SaveURL(),
			command.Prompt(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
