package command

import (
	"fmt"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/urfave/cli/v2"
)

func Download() *cli.Command {
	return &cli.Command{
		Name:      "download",
		UsageText: "dropbox-cli download <remote-path> <local-path>[/]",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token"},
		},
		Action: func(c *cli.Context) error {
			r := New(c.String("token"))

			remoteRootPath := toRemotePath(formatPath(c.Args().Get(0)))                     // left slash, right no slash
			localRootPath := toLocalPath(formatPathByRev(c.Args().Get(0), c.Args().Get(1))) // left, right both no slash
			if remoteRootPath == "" {
				fmt.Printf("> download fail: empty remote path.\n")
				return fmt.Errorf("empty remote path")
			}
			if localRootPath == "" {
				fmt.Printf("> download %q fail: empty local path.\n", remoteRootPath)
				return fmt.Errorf("empty local path")
			}

			fmt.Printf("> start download %q to %q.\n", remoteRootPath, localRootPath)

			err := r.ListFolder(remoteRootPath, func(data files.IsMetadata) error {
				switch v := data.(type) {
				case *files.FileMetadata:
					local := formatRelatePath(localRootPath, remoteRootPath, v.PathDisplay)
					if r.TryCheckLocalContentHash(local, v.ContentHash) {
						fmt.Printf("> download file %q to %q exist, skip.\n", v.PathDisplay, local)
						return nil
					}
					body, err := r.Download(v.Id)
					if err != nil {
						fmt.Printf("> download file %q to %q fail: %s.\n", v.PathDisplay, local, err)
						return err
					}
					defer body.Close()
					err = r.WriteToLocal(local, body)
					if err != nil {
						fmt.Printf("> download file %q to %q fail: %s.\n", v.PathDisplay, local, err)
						return err
					}
					fmt.Printf("> download file %q to %q success.\n", v.PathDisplay, local)
				case *files.FolderMetadata:
					local := formatRelatePath(localRootPath, remoteRootPath, v.PathDisplay)
					localInfo, _ := os.Stat(local)
					if localInfo != nil {
						fmt.Printf("> create dir %q to %q exist, skip.\n", v.PathDisplay, local)
						return nil
					}

					if err := os.MkdirAll(local, os.ModePerm); err != nil {
						fmt.Printf("> create dir %q to %q fail: %s.\n", v.PathDisplay, local, err)
						return err
					}
					fmt.Printf("> create file %q to %q success.\n", v.PathDisplay, local)
				case *files.DeletedMetadata:
					local := formatRelatePath(localRootPath, remoteRootPath, v.PathDisplay)
					if fileToDelete, _ := os.Stat(local); fileToDelete != nil {
						_ = os.RemoveAll(local)
						fmt.Printf("> remove file %q success.\n", local)
					}
				}
				return nil
			})
			if err != nil {
				return err
			}

			fmt.Printf("> download %q to %q success.\n", remoteRootPath, localRootPath)
			return nil
		},
	}
}
