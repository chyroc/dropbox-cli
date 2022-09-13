package command

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
)

func SaveURL() *cli.Command {
	return &cli.Command{
		Name:      "save-url",
		UsageText: "dropbox-cli save-url <http-url> [<remote-path>]",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token"},
		},
		Action: func(c *cli.Context) error {
			var err error
			r := New(c.String("token"), false)

			if c.Args().Len() < 1 {
				return cli.ShowCommandHelp(c, "save-url")
			}

			uri := c.Args().Get(0)
			path := ""
			if c.Args().Len() > 1 {
				path = c.Args().Get(1)
			} else {
				path, err = urlToFilepath(uri)
				if err != nil {
					return err
				}
			}
			path = "/" + path

			fmt.Printf("> start save url %q to %q.\n", uri, path)

			jobID, err := r.SaveURL(path, uri)
			if err != nil {
				return err
			}

			for {
				status, err := r.CheckSaveURLJob(jobID)
				if err != nil {
					return err
				}
				fmt.Println(status)
				if status == "failed" || status == "complete" {
					break
				}
			}
			return nil
		},
	}
}

func urlToFilepath(uri string) (string, error) {
	uriParsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/%s", uriParsed.Host, uriParsed.Path)
	for _, v := range []string{":", "/", "+", "-", " "} {
		path = strings.ReplaceAll(path, v, "_")
	}
	return path, nil
}
