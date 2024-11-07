package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, color.RedString("usage: %s [url] [filename]\n", os.Args[0]))
		os.Exit(1)
	}

	url := os.Args[1]
	filename := os.Args[2]

	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)

	defer check(resp.Body.Close)

	f, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)

	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	io.Copy(io.MultiWriter(f, bar), resp.Body)
}

func check(f func() error) {
	if err := f(); err != nil {
		fmt.Fprintf(os.Stderr, "received error: %v\n", err)
	}
}
