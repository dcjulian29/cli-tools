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
	"github.com/theckman/yacspin"
)

type RepoInfo struct {
	folder string
	data   string
}

func main() {
	pwd, _ := os.Getwd()

	if len(os.Args) > 1 {
		pwd = os.Args[1]
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

			results <- process(d)
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

	spinner.Stop()

	for _, info := range processed {
		if len(info.data) > 0 {
			fmt.Printf("\n%s\n%s", color.HiYellowString(info.folder), info.data)
		}
	}
}

func process(path string) RepoInfo {
	_ = executeGit(path, "fetch")

	dirty := executeGit(path, "diff", "--stat")
	local := executeGit(path, "rev-parse", "@")
	remote := executeGit(path, "rev-parse", "@{u}")
	base := executeGit(path, "merge-base", "@", "@{u}")
	data := ""
	pull := false
	push := false

	if local != remote {
		if local == base {
			pull = true
		} else {
			if remote == base {
				push = true
			} else {
				pull = true
				push = true
			}
		}
	}

	if len(dirty) > 0 {
		data = color.CyanString("Skipping as it isn't clean...\n")
	} else {
		if pull {
			data = executeGit(path, "pull", "--rebase", "--prune", "--recurse-submodules=yes")
		}

		if push {
			if len(data) > 0 {
				data = fmt.Sprintf("%s\n\n%s", data, executeGit(path, "push"))
			} else {
				data = executeGit(path, "push")
			}
		}
	}

	return RepoInfo{
		folder: path,
		data:   data,
	}
}

func executeGit(path string, params ...string) string {
	cmd := exec.Command("git", params...)
	cmd.Dir = path

	out, err := cmd.CombinedOutput()

	if err != nil {
		return err.Error()
	}

	return string(out)
}
