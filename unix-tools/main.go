/*
Copyright © 2023 Julian Easterling julian@julianscorner.com

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
	"os"
	"path/filepath"
	"strings"

	"github.com/dcjulian29/go-toolbox/docker"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
)

func main() {
	binary, interactive, entrypoint, image, tag, content := determineValues()
	args := filesystem.EnsureUnixPathArguments()
	data, work := docker.HostContainerVolume()
	volumes := []string{data}

	if len(content) > 0 {
		if err := filesystem.EnsureFileExist(entrypoint, content); err != nil {
			fmt.Println(textformat.Fatal(err.Error()))
			os.Exit(1)
		}

		volumes = append(volumes, fmt.Sprintf("%s:/docker-entrypoint.sh", strings.ReplaceAll(entrypoint, "/", "\\")))
		entrypoint = "/docker-entrypoint.sh"
		binary = ""
	}

	opts := docker.ContainerOptions{
		AdditionalArgs:   strings.Join(args, " "),
		Command:          binary,
		EntryPoint:       entrypoint,
		Image:            image,
		Interactive:      interactive,
		Tag:              tag,
		Volumes:          volumes,
		WorkingDirectory: work,
	}

	if _, err := docker.Run(opts); err != nil {
		fmt.Println(textformat.Fatal(err.Error()))
		os.Exit(1)
	}
}

func determineValues() (string, bool, string, string, string, []byte) {
	temp, _ := os.LookupEnv("TEMP")
	prefix := "/usr/bin"
	binary := strings.ReplaceAll(filepath.Base(os.Args[0]), ".exe", "")
	interactive := false
	entrypoint := ""
	custompoint := strings.ReplaceAll(fmt.Sprintf("%s\\docker-entrypoint.sh", temp), "\\", "/")
	image := "alpine"
	tag := "latest"
	var content []byte

	switch binary {
	case "alpine":
		prefix = "/bin"
		binary = "sh"
		interactive = true

	case "base64":
		prefix = "/bin"

	case "cat":
		prefix = "/bin"

	case "curl":
		image = "curlimages/curl"
		entrypoint = custompoint
		content = []byte(`#!/bin/sh

/usr/bin/curl $@
`)

	case "debian":
		prefix = "/bin"
		binary = "bash"
		image = "debian"
		tag = "stable"
		interactive = true

	case "dog":
		entrypoint = custompoint
		content = []byte(`#!/bin/sh

/sbin/apk add dog > /dev/null
/usr/bin/dog $@
`)

	case "grep":
		prefix = "/bin"

	case "gunzip":
		prefix = "/bin"

	case "gzip":
		prefix = "/bin"

	case "http":
		entrypoint = custompoint
		content = []byte(`#!/bin/sh

/sbin/apk add httpie > /dev/null
/usr/bin/http $@
`)

	case "jq":
		entrypoint = custompoint
		content = []byte(`#!/bin/sh

/sbin/apk add jq > /dev/null
/usr/bin/jq $@
`)
		interactive = true

	case "nano":
		entrypoint = custompoint
		content = []byte(`#!/bin/sh

/sbin/apk add nano > /dev/null
/usr/bin/nano $@
`)
		interactive = true

	case "sed":
		prefix = "/bin"

	case "tar":
		prefix = "/bin"

	case "yamllint":
		entrypoint = custompoint
		content = []byte(`#!/bin/sh

/sbin/apk add yamllint > /dev/null
/usr/bin/yamllint $@
`)

	case "yq":
		entrypoint = custompoint
		content = []byte(`#!/bin/sh

/sbin/apk add yq > /dev/null
/usr/bin/yq $@
`)
		interactive = true

	case "zcat":
		prefix = "/bin"
	}

	binary = fmt.Sprintf("%s/%s", prefix, binary)

	return binary, interactive, entrypoint, image, tag, content
}
