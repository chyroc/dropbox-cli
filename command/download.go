package command

import (
	"fmt"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/urfave/cli/v2"
)

func Download() *cli.Command {
	return &cli.Command{
		Name: "download",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token"},
		},
		Action: func(c *cli.Context) error {
			r := New(c.String("token"))

			path := r.formatRemotePath(c.Args().First())
			r.SetRootPath(path)

			fmt.Printf("> start download %q to %q.\n", r.RemoteRootPath, r.LocalRootPath)

			err := r.ListFolder(path, func(data files.IsMetadata) error {
				switch v := data.(type) {
				case *files.FileMetadata:
					local := r.generateLocalName(v.PathDisplay)
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
					local := r.generateLocalName(v.PathDisplay)
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
				}
				return nil
			})
			if err != nil {
				return err
			}

			fmt.Printf("> download %q to %q success.\n", r.RemoteRootPath, r.LocalRootPath)
			return nil
		},
	}
}
