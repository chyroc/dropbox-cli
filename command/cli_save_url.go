package command

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/async"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func (r *Cli) SaveURL(filename, fileURL string) (jobID string, err error) {
	resp, err := r.fileClient.SaveUrl(&files.SaveUrlArg{
		Path: filename,
		Url:  fileURL,
	})
	if err != nil {
		return "", err
	}
	return resp.AsyncJobId, nil
}

func (r *Cli) CheckSaveURLJob(jobID string) (status string, err error) {
	res, err := r.fileClient.SaveUrlCheckJobStatus(&async.PollArg{
		AsyncJobId: jobID,
	})
	if err != nil {
		return "", err
	}
	return res.Tag, nil
}
