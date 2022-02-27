package command

import (
	"fmt"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/ivpusic/grpool"
	"github.com/urfave/cli/v2"
)

func Download() *cli.Command {
	return &cli.Command{
		Name:      "download",
		UsageText: "dropbox-cli download <remote-path> <local-path>[/]",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token"},
			&cli.BoolFlag{Name: "disable-cursor-cache"},
		},
		Action: func(c *cli.Context) error {
			r := New(c.String("token"), c.Bool("disable-cursor-cache"))

			remoteRootPath := toRemotePath(formatPath(c.Args().Get(0)))                     // left slash, right no slash
			localRootPath := toLocalPath(formatPathByRev(c.Args().Get(0), c.Args().Get(1))) // left, right both no slash
			r.setLocalRootPath(localRootPath)

			if remoteRootPath == "" {
				fmt.Printf("> download fail: empty remote path.\n")
				return fmt.Errorf("empty remote path")
			}
			if localRootPath == "" {
				fmt.Printf("> download %q fail: empty local path.\n", remoteRootPath)
				return fmt.Errorf("empty local path")
			}

			fmt.Printf("> start download %q to %q.\n", remoteRootPath, localRootPath)

			pool := grpool.NewPool(20, 100)
			defer pool.Release()

			err := r.ListFolder(remoteRootPath, func(data files.IsMetadata) error {
				v := data
				pool.JobQueue <- func() {
					_ = r.download(localRootPath, remoteRootPath, v)
				}
				return nil
			})
			if err != nil {
				return err
			}

			pool.WaitAll()

			fmt.Printf("> download %q to %q success.\n", remoteRootPath, localRootPath)
			return nil
		},
	}
}

func (r *Cli) download(localRootPath, remoteRootPath string, data files.IsMetadata) error {
	switch v := data.(type) {
	case *files.FileMetadata:
		return r.downloadFile(localRootPath, remoteRootPath, v)
	case *files.FolderMetadata:
		return r.downloadDir(localRootPath, remoteRootPath, v)
	case *files.DeletedMetadata:
		return r.downloadDel(localRootPath, remoteRootPath, v)
	default:
		return nil
	}
}

func (r *Cli) downloadFile(localRootPath, remoteRootPath string, v *files.FileMetadata) error {
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
	_ = os.Chtimes(local, v.ClientModified, v.ClientModified)
	fmt.Printf("> download file %q to %q success.\n", v.PathDisplay, local)
	return nil
}

func (r *Cli) downloadDir(localRootPath, remoteRootPath string, v *files.FolderMetadata) error {
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
	return nil
}

func (r *Cli) downloadDel(localRootPath, remoteRootPath string, v *files.DeletedMetadata) error {
	local := formatRelatePath(localRootPath, remoteRootPath, v.PathDisplay)
	if fileToDelete, _ := os.Stat(local); fileToDelete != nil {
		_ = os.RemoveAll(local)
		fmt.Printf("> remove file %q success.\n", local)
	}
	return nil
}
