package command

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (r *Cli) formatRemotePath(path string) string {
	path = "/" + strings.Trim(strings.TrimSpace(path), "/")
	return path
}

// `.` or `GitHub`
func (r *Cli) formatLocalPath(path string) string {
	if strings.HasPrefix(path, "./") {
		path = path[1:]
	}
	path = strings.Trim(strings.TrimSpace(path), "/")
	if path == "" {
		path = "."
	}
	return path
}

func (r *Cli) formatRevRemotePath(localPath, path, remotePath string) string {
	// localPath  = a/b
	// path       = a/b/c/d.txt
	// remotePath = Git

	// result     = Git/c/d.txt
	return fmt.Sprintf("%s/%s", remotePath, path[len(localPath)+1:])
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

func (r *Cli) ListLocal(dir string, f func(path string, info fs.FileInfo) error) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return f(dir, info)
	}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, info := range fs {
		path := r.formatLocalPath(fmt.Sprintf("%s/%s", dir, info.Name()))
		if err := f(path, info); err != nil {
			return err
		}
		if info.IsDir() {
			if err := r.ListLocal(path, f); err != nil {
				return err
			}
		}
	}
	return nil
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

func (r *Cli) GenContentHash(localContent []byte) (string, error) {
	blocks := []byte{}
	size := 1024 * 1024 * 4
	for i := 0; i < len(localContent); i += size {
		if i+size > len(localContent) {
			d := sha256.Sum256(localContent[i:])
			blocks = append(blocks, d[:]...)
		} else {
			d := sha256.Sum256(localContent[i : i+size])
			blocks = append(blocks, d[:]...)
		}
	}
	result := sha256.Sum256(blocks)
	return fmt.Sprintf("%x", result[:]), nil
}
