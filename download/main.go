/*
Copyright © 2026 Julian Easterling

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
