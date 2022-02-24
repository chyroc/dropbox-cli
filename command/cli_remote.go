package command

import (
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func (r *Cli) ListFolder(path string, f func(data files.IsMetadata) error) error {
	if path == "/" {
		path = ""
	}

	cursor := ""
	for {
		var resp *files.ListFolderResult
		var err error
		if cursor == "" {
			resp, err = r.fileClient.ListFolder(&files.ListFolderArg{
				Path:                        path,
				Recursive:                   true,
				IncludeDeleted:              false,
				Limit:                       100,
				IncludeNonDownloadableFiles: false,
			})
		} else {
			resp, err = r.fileClient.ListFolderContinue(&files.ListFolderContinueArg{
				Cursor: cursor,
			})
		}
		if err != nil {
			return err
		}
		for _, v := range resp.Entries {
			err = f(v)
			if err != nil {
				return err
			}
		}
		if !resp.HasMore {
			break
		}
		cursor = resp.Cursor
	}

	return nil
}

func (r *Cli) Download(path string) (io.ReadCloser, error) {
	_, body, err := r.fileClient.Download(&files.DownloadArg{
		Path: path,
	})
	return body, err
}
