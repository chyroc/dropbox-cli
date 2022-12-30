package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

type promptClient struct {
	cli             *Cli
	suggestions     []prompt.Suggest
	currentDir      string
	currentChildren map[string]string
	promptTitle     string
}

func newPromptClient(cli *Cli, promptTitle string) *promptClient {
	return &promptClient{
		cli:             cli,
		suggestions:     nil,
		currentDir:      "/",
		currentChildren: map[string]string{},
		promptTitle:     promptTitle,
	}
}

func (r *promptClient) Run() error {
	p := prompt.New(
		r.executor,
		r.completer,
		// prompt.OptionHistory(readHistory()),
		prompt.OptionPrefix(r.prefix()),
		prompt.OptionPrefixTextColor(prompt.White),
		prompt.OptionLivePrefix(r.livePrefix),
		prompt.OptionTitle(r.promptTitle),
		prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 0x62},
			Fn:        prompt.GoLeftWord,
		}, prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 0x66},
			Fn:        prompt.GoRightWord,
		}, prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 0x7f},
			Fn:        prompt.DeleteWord,
		}, prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 0x08},
			Fn:        prompt.DeleteWord,
		}),
	)
	p.Run()

	return nil
}

func (r *promptClient) executor(in string) {
	in = strings.TrimSpace(in)
	if in == "ls" || in == "l" {
		r.listFiles(r.currentDir)
	} else if strings.HasPrefix(in, "cd ") {
		r.cdDir(strings.SplitN(in, " ", 2)[1])
	} else if in == "cd" {
		r.cdDir("/")
	}
}

func (r *promptClient) completer(in prompt.Document) []prompt.Suggest {
	return nil
}

func (r *promptClient) prefix() string {
	if r.currentDir == "" {
		return "dropbox:/> "
	}

	return fmt.Sprintf("dropbox:%s> ", r.currentDir)
}

func (r *promptClient) livePrefix() (string, bool) {
	return r.prefix(), true
}

func (r *promptClient) listFiles(dir string) {
	header := []string{"Name", "Size", "ClientModified"}
	datas := [][]string{}
	err := r.cli.ListFolder(dir, false, false, func(data files.IsMetadata) error {
		switch v := data.(type) {
		case *files.FileMetadata:
			datas = append(datas, []string{v.Name, formatSize(v.Size), v.ClientModified.Format("2006-01-02 15:04:05")})
		case *files.FolderMetadata:
			datas = append(datas, []string{formatDir(v.Name), "", ""})
		default:
			return nil
		}
		return nil
	})
	if err != nil {
		printError(err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(datas)
	table.Render()
}

func (r *promptClient) cdDir(dir string) {
	dir = joinPath(r.currentDir, dir)
	meta, err := r.cli.GetMetadata(dir)
	if err != nil {
		printError(err)
		return
	}
	if meta == nil {
		// root dir
		r.currentDir = dir
		return
	}
	if _, ok := meta.(*files.FolderMetadata); !ok {
		printError(fmt.Errorf("'%s' is not a directory", dir))
		return
	}
	r.currentDir = dir
}

func joinPath(cur string, path string) string {
	if strings.HasPrefix(path, "/") {
		// /路径/ -> /路径
		return strings.TrimRight(path, "/")
	}
	if cur == "" {
		cur = "/"
	}
	cur = strings.TrimRight(filepath.Join(cur, path), "/")
	if !strings.HasPrefix(cur, "/") {
		cur = "/" + cur
	}
	return cur
}

var (
	dirColor = color.New(color.FgCyan).Add(color.Bold)
	redColor = color.New(color.FgRed).Add(color.Bold)
)

func printError(err error) {
	if err != nil {
		_, _ = redColor.Println(err.Error())
	}
}

func formatDir(v string, args ...interface{}) string {
	if len(args) == 0 {
		return dirColor.Sprint(v)
	}
	return dirColor.Sprintf(v, args...)
}

func formatSize(size uint64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%d KB", size/1024)
	}
	if size < 1024*1024*1024 {
		return fmt.Sprintf("%d MB", size/1024/1024)
	}
	return fmt.Sprintf("%d GB", size/1024/1024/1024)
}
