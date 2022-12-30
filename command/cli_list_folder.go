package command

import (
	"fmt"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func (r *Cli) ListFolder(path string, recursive, includeDeleted bool, f func(data files.IsMetadata) error) error {
	if path == "/" {
		path = ""
	}

	meta, err := r.fileClient.GetMetadata(&files.GetMetadataArg{Path: path})
	if err != nil {
		if !strings.Contains(err.Error(), "The root folder is unsupported.") {
			return err
		}
		// continue to list folder
	} else {
		if file, ok := meta.(*files.FileMetadata); ok {
			return f(file)
		}
	}

	cursor := ""
	if !r.disableCursorCache {
		cursor = r.getCursor()
		fmt.Printf("> [meta] get cursor=%s\n", cursor)
	}
	for {
		var resp *files.ListFolderResult
		var err error
		if cursor == "" {
			resp, err = r.fileClient.ListFolder(&files.ListFolderArg{
				Path:                        path,
				Recursive:                   recursive,
				IncludeDeleted:              includeDeleted,
				Limit:                       100,
				IncludeNonDownloadableFiles: false,
			})
		} else {
			resp, err = r.fileClient.ListFolderContinue(&files.ListFolderContinueArg{
				Cursor: cursor,
			})
		}
		if !r.disableCursorCache && resp != nil && resp.Cursor != "" {
			r.setCursor(cursor)
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
