package command

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func (r *Cli) ListFolder(path string, f func(data files.IsMetadata) error) error {
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
				Recursive:                   true,
				IncludeDeleted:              true,
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

func (r *Cli) Download(path string) (io.ReadCloser, error) {
	_, body, err := r.fileClient.Download(&files.DownloadArg{
		Path: path,
	})
	return body, err
}

type UploadResult struct {
	Exist bool
}

func (r *Cli) Upload(localFile, remotePath string, blockSuccessCallback func(idx int)) (*UploadResult, error) {
	res := new(UploadResult)
	info, err := os.Stat(localFile)
	if err != nil {
		return res, err
	}

	if info.IsDir() {
		return res, nil
	}

	localContent, err := ioutil.ReadFile(localFile)
	if err != nil {
		return nil, err
	}
	localContentLength := uint64(len(localContent))

	if fileMeta, _ := r.fileClient.GetMetadata(&files.GetMetadataArg{
		Path:                            remotePath,
		IncludeMediaInfo:                true,
		IncludeDeleted:                  false,
		IncludeHasExplicitSharedMembers: false,
		IncludePropertyGroups:           nil,
	}); fileMeta != nil {
		meta, _ := fileMeta.(*files.FileMetadata)
		if meta != nil {
			if r.TryCheckLocalContentHash([]byte(localContent), meta.ContentHash) {
				res.Exist = true
				return res, nil
			}
		}
	}

	commitInfo := files.CommitInfo{
		Path:       remotePath,
		Mode:       &files.WriteMode{Tagged: dropbox.Tagged{Tag: files.WriteModeOverwrite}},
		Autorename: false,
		// ClientModified: ptr.Time(info.ModTime()),
		Mute:           true,
		PropertyGroups: nil,
		StrictConflict: false,
	}

	// 小于 20M，直接上传
	if info.Size()/1024/1024 <= 20 {
		_, err = r.fileClient.Upload(&files.UploadArg{
			CommitInfo: commitInfo,
		}, bytes.NewReader(localContent))
		return res, err
	}

	// 否则分片上传
	blockSize := uint64(4 * 1024 * 1024) // 4M
	resp, err := r.fileClient.UploadSessionStart(&files.UploadSessionStartArg{
		Close:       false,
		SessionType: &files.UploadSessionType{Tagged: dropbox.Tagged{Tag: files.UploadSessionTypeSequential}},
	}, bytes.NewReader(localContent[:blockSize]))
	if err != nil {
		return nil, err
	}
	blockSuccessCallback(1)

	// 第一个 4M 已经使用了，所以跳过第一个 4M
	idx := 1
	for i := uint64(blockSize); i < localContentLength; i += blockSize {
		idx++
		var reader io.Reader
		isClose := false
		if i+blockSize >= localContentLength {
			reader = bytes.NewReader(localContent[i:])
			isClose = true
		} else {
			reader = bytes.NewReader(localContent[i : i+blockSize])
		}

		// 这里不能使用并发, 必须一个接着一个的上传
		err = r.fileClient.UploadSessionAppendV2(&files.UploadSessionAppendArg{
			Cursor: &files.UploadSessionCursor{
				SessionId: resp.SessionId,
				Offset:    i,
			},
			Close: isClose,
		}, reader)
		if err != nil {
			return nil, err
		}
		blockSuccessCallback(idx)
	}

	_, err = r.fileClient.UploadSessionFinish(&files.UploadSessionFinishArg{
		Cursor: &files.UploadSessionCursor{
			SessionId: resp.SessionId,
			Offset:    localContentLength,
		},
		Commit: &commitInfo,
	}, nil)

	return res, err
}
