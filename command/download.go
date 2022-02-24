package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Download() *cli.Command {
	return &cli.Command{
		Name: "download",
		Action: func(c *cli.Context) error {
			fmt.Println("download")

			return nil
		},
	}
}
