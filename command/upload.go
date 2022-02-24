package command

import (
	"fmt"
	"io/fs"

	"github.com/urfave/cli/v2"
)

func Upload() *cli.Command {
	return &cli.Command{
		Name:      "upload",
		UsageText: "dropbox-cli upload <local-path> <remote-path>[/]",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token"},
		},
		Action: func(c *cli.Context) error {
			r := New(c.String("token"))

			localRootPath := toLocalPath(formatPath(c.Args().Get(0)))                         // left, right both no slash
			remoteRootPath := toRemotePath(formatPathByRev(c.Args().Get(0), c.Args().Get(1))) // left slash, right no slash
			if localRootPath == "" {
				fmt.Printf("> upload fail: empty local path.\n")
				return fmt.Errorf("empty local path")
			}
			if remoteRootPath == "" {
				fmt.Printf("> upload %q fail: empty remote path.\n", localRootPath)
				return fmt.Errorf("empty remote path")
			}

			fmt.Printf("> start upload %q to %q.\n", localRootPath, remoteRootPath)

			err := r.ListLocal(localRootPath, func(localPath string, info fs.FileInfo) error {
				if info.IsDir() {
					return nil
				}

				remotePath := formatRelatePath(remoteRootPath, localRootPath, localPath)

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
			})
			if err != nil {
				return err
			}

			// fmt.Printf("> upload %q to %q success.\n", localPath, remotePath)
			return nil
		},
	}
}
