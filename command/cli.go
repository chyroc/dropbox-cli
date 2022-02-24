package command

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type Cli struct {
	config     dropbox.Config
	fileClient files.Client
}

func New(token string) *Cli {
	r := new(Cli)
	r.config = dropbox.Config{
		Token: token,
	}

	r.fileClient = files.New(r.config)

	return r
}
