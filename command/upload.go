package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Upload() *cli.Command {
	return &cli.Command{
		Name: "upload",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token"},
		},
		Action: func(c *cli.Context) error {
			r := New(c.String("token"))

			localPath := r.formatLocalPath(c.Args().First())
			remotePath := r.formatRemotePath(c.Args().Get(1))

			fmt.Printf("> start upload %q to %q.\n", localPath, remotePath)

			res, err := r.Upload(localPath, remotePath, func(idx int) {
				fmt.Printf("> upload %q to %q block[%d] success.\n", localPath, remotePath, idx)
			})
			if err != nil {
				fmt.Printf("> upload %q to %q fail: %s.\n", localPath, remotePath, err)
				return err
			} else if res.Exist {
				fmt.Printf("> upload %q to %q exist, skip.\n", localPath, remotePath)
				return nil
			}
			fmt.Printf("> upload %q to %q success.\n", localPath, remotePath)
			return nil
		},
	}
}
