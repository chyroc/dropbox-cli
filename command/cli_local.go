package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (r *Cli) formatPath(path string) string {
	path = "/" + strings.Trim(strings.TrimSpace(path), "/")
	return path
}

func (r *Cli) generateLocalName(remotePath string) string {
	// /GitHub/100days/meta.json -> 100days/meta.json
	if strings.ToLower(remotePath)[0:len(r.RemoteRootPath)] != strings.ToLower(r.RemoteRootPath) {
		panic(fmt.Sprintf("%q not start with %q", remotePath, r.RemoteRootPath))
	}
	remotePath = strings.Trim(remotePath[len(r.RemoteRootPath):], "/")
	return fmt.Sprintf("%s/%s", r.LocalRootPath, remotePath)
}

func (r *Cli) SetRootPath(path string) {
	localDir := filepath.Base(strings.TrimSpace(path))
	if path == "/" {
		localDir = "dropbox-local"
	}

	r.LocalRootPath = localDir                                          // 100days
	r.RemoteRootPath = "/" + strings.Trim(strings.TrimSpace(path), "/") // /GitHub/100days
}

func (r *Cli) WriteToLocal(filename string, body io.Reader) error {
	// defer body.Close()

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, body)
	return err
}
