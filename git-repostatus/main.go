package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/theckman/yacspin"
)

type RepoInfo struct {
	folder      string
	dirty       bool
	push_needed bool
	pull_needed bool
	diverged    bool
	untracked   bool
}

var (
	Only_Dirty     bool
	Only_Push      bool
	Only_Pull      bool
	Only_Diverged  bool
	Only_Untracked bool
)

func main() {
	pwd, _ := os.Getwd()

	if len(os.Args) > 1 {
		pwd = os.Args[1]
	}

	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "dirty":
			Only_Dirty = true
		case "push":
			Only_Push = true
		case "pull":
			Only_Pull = true
		case "diverged":
			Only_Diverged = true
		case "untracked":
			Only_Untracked = true
		}
	}

	spinner, _ := yacspin.New(yacspin.Config{
		Frequency: 100 * time.Millisecond,
		Colors:    []string{"fgYellow"},
		CharSet:   yacspin.CharSets[69],
	})

	spinner.Start()

	var wg sync.WaitGroup
	var directories []string

	results := make(chan RepoInfo)
	processed := make([]RepoInfo, 0)

	err := filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if info.IsDir() && info.Name() == ".git" {
				directories = append(directories, filepath.Dir(path))
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, d := range directories {
		wg.Add(1)

		go func(d string) {
			defer wg.Done()

			info := process(d)

			if Only_Dirty {
				if info.dirty {
					results <- info
				}
			} else if Only_Pull {
				if info.pull_needed {
					results <- info
				}
			} else if Only_Push {
				if info.push_needed {
					results <- info
				}
			} else if Only_Diverged {
				if info.diverged {
					results <- info
				}
			} else if Only_Untracked {
				if info.untracked {
					results <- info
				}
			} else {
				results <- info
			}
		}(d)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for item := range results {
		processed = append(processed, item)
	}

	sort.Slice(processed, func(i, j int) bool {
		return strings.Compare(processed[i].folder, processed[j].folder) < 0
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PATH", "DIRTY", "PUSH", "PULL", "DIVERGED", "UNTRACKED"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, info := range processed {
		table.Append([]string{
			output_path(info),
			output_yellow(info.dirty),
			output_red(info.push_needed),
			output_red(info.pull_needed),
			output_red(info.diverged),
			output_yellow(info.untracked),
		})
	}

	spinner.Stop()

	table.Render()
}

func output_path(info RepoInfo) string {
	if info.dirty {
		return color.YellowString(info.folder)
	}

	if info.pull_needed || info.push_needed {
		return color.RedString(info.folder)
	}

	return color.GreenString(info.folder)
}

func output_red(f bool) string {
	if f {
		return color.RedString("yes")
	}

	return color.GreenString("no")

}

func output_yellow(f bool) string {
	if f {
		return color.YellowString("yes")
	}

	return color.GreenString("no")

}

func process(path string) RepoInfo {
	_ = executeExternalProgramCapture("git", path, "fetch")
	dirty := executeExternalProgramCapture("git", path, "diff", "--stat")
	untracked := executeExternalProgramCapture("git", path, "ls-files", "--others", "--exclude-standard")
	local := executeExternalProgramCapture("git", path, "rev-parse", "@")
	remote := executeExternalProgramCapture("git", path, "rev-parse", "@{u}")
	base := executeExternalProgramCapture("git", path, "merge-base", "@", "@{u}")
	pull := false
	push := false
	diverged := false

	if local != remote {
		if local == base {
			pull = true
		} else {
			if remote == base {
				push = true
			} else {
				diverged = true
			}
		}
	}

	return RepoInfo{
		folder:      path,
		dirty:       len(dirty) > 0,
		push_needed: push,
		pull_needed: pull,
		diverged:    diverged,
		untracked:   len(untracked) > 0,
	}
}

func executeExternalProgramCapture(program string, path string, params ...string) string {
	cmd := exec.Command(program, params...)
	cmd.Stdin = os.Stdin
	cmd.Dir = path
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
	}

	return string(out)
}
