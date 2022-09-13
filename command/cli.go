package command

import (
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type Cli struct {
	config     dropbox.Config
	fileClient files.Client

	disableCursorCache bool
	localRootPath      string
}

func New(token string, disableCursorCache bool) *Cli {
	if token == "" {
		token = os.Getenv("DROPBOX_TOKEN")
	}
	r := new(Cli)
	r.config = dropbox.Config{
		Token: token,
	}
	r.disableCursorCache = disableCursorCache

	r.fileClient = files.New(r.config)

	return r
}
