package command

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (r *Cli) ListLocal(dir string, f func(path string, info os.FileInfo) error) error {
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
		path := toLocalPath(formatPath(fmt.Sprintf("%s/%s", dir, info.Name())))
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

// left no /
// right may have /
func formatPath(path string) string {
	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "./") {
		path = path[1:]
	}
	// like /a/b/c/ or /a/b/c

	path = strings.TrimLeft(path, "/")
	// like a/b/c/ or a/b/c
	return path
}

// return like a or a/b or a/b/c
// left or right both has no /
func formatPathByRev(base, path string) string {
	// like some
	base = filepath.Base(strings.Trim(strings.TrimSpace(base), "/"))

	// like a/b/c/ or a/b/c
	path = formatPath(path)

	// path = a/b/c/, base=some, => a/b/c/some
	if strings.HasSuffix(path, "/") {
		path = strings.Trim(path, "/") + "/" + base
	}

	return path
}

func toRemotePath(path string) string {
	return "/" + strings.Trim(path, "/")
}

func toLocalPath(path string) string {
	return strings.Trim(path, "/")
}

func formatRelatePath(anotherRoot, sideRoot, sideFull string) string {
	if sideRoot == sideFull {
		return anotherRoot
	}
	// 读取的 sideFull path 前几位，必须和 sideRoot 相同
	if strings.ToLower(sideFull)[0:len(sideRoot)] != strings.ToLower(sideRoot) {
		panic(fmt.Sprintf("%q not start with %q", sideFull, sideRoot))
	}
	// 去掉 sideRoot，得到后面的相对路径
	base := strings.Trim(sideFull[len(sideRoot):], "/")

	// 将相对路径和 anotherRoot path 拼接，得到最终的相对路径
	return fmt.Sprintf("%s/%s", strings.TrimRight(anotherRoot, "/"), base)
}
