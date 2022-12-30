package command

import (
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func (r *Cli) GetMetadata(path string) (files.IsMetadata, error) {
	if path == "/" {
		path = ""
	}

	meta, err := r.fileClient.GetMetadata(&files.GetMetadataArg{Path: path})
	if err != nil {
		if !strings.Contains(err.Error(), "The root folder is unsupported.") {
			return nil, err
		}
		return nil, nil
	} else {
		return meta, nil
	}
}
