package command

import (
	"github.com/urfave/cli/v2"
)

func Prompt() *cli.Command {
	return &cli.Command{
		Name:      "prompt",
		UsageText: "dropbox-cli prompt",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token"},
		},
		Action: func(c *cli.Context) error {
			r := New(c.String("token"), true)

			return newPromptClient(r, "dropbox-cli").Run()
		},
	}
}
