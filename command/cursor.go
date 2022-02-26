package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Meta struct {
	Cursor string `json:"cursor"`
}

func (r *Cli) setLocalRootPath(path string) {
	r.localRootPath = path
}

func (r *Cli) getCursor() string {
	meta := r.loadMeta()
	return meta.Cursor
}

func (r *Cli) setCursor(cursor string) {
	file := fmt.Sprintf("%s/meta.json", r.localRootPath)
	meta := r.loadMeta()
	meta.Cursor = cursor
	if bs, _ := json.Marshal(meta); len(bs) > 0 {
		err := ioutil.WriteFile(file, bs, 0644)
		if err != nil {
			fmt.Printf("> [meta] save cursor fail: %s\n", err.Error())
		}
	}
}

func (r *Cli) loadMeta() Meta {
	file := fmt.Sprintf("%s/meta.json", r.localRootPath)
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return Meta{}
	}
	var meta Meta
	_ = json.Unmarshal(bs, &meta)
	return meta
}
