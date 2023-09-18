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
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	temp, _ := os.LookupEnv("TEMP")
	prefix := "/usr/bin"
	binary := strings.ReplaceAll(filepath.Base(os.Args[0]), ".exe", "")
	interactive := "-it"
	entrypoint := ""
	var content []byte

	switch binary {
	case "alpine":
		prefix = "/bin"
		binary = "sh"
	case "base64":
		prefix = "/bin"
	case "cat":
		prefix = "/bin"
	case "curl":
		entrypoint = strings.ReplaceAll(fmt.Sprintf("%s\\docker-entrypoint.sh", temp), "\\", "/")
		content = []byte(`#!/bin/sh

/sbin/apk add curl > /dev/null
/usr/bin/curl $@
`)
	case "doq":
		entrypoint = strings.ReplaceAll(fmt.Sprintf("%s\\docker-entrypoint.sh", temp), "\\", "/")
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
		entrypoint = strings.ReplaceAll(fmt.Sprintf("%s\\docker-entrypoint.sh", temp), "\\", "/")
		content = []byte(`#!/bin/sh

/sbin/apk add httpie > /dev/null
/usr/bin/http $@
`)
	case "jq":
		interactive = "-i"
		entrypoint = strings.ReplaceAll(fmt.Sprintf("%s\\docker-entrypoint.sh", temp), "\\", "/")
		content = []byte(`#!/bin/sh

/sbin/apk add jq > /dev/null
/usr/bin/jq $@
`)
	case "nano":
		entrypoint = strings.ReplaceAll(fmt.Sprintf("%s\\docker-entrypoint.sh", temp), "\\", "/")
		content = []byte(`#!/bin/sh

/sbin/apk add nano > /dev/null
/usr/bin/nano $@
`)
	case "sed":
		prefix = "/bin"
	case "tar":
		prefix = "/bin"
	case "yamllint":
		entrypoint = strings.ReplaceAll(fmt.Sprintf("%s\\docker-entrypoint.sh", temp), "\\", "/")
		content = []byte(`#!/bin/sh

/sbin/apk add yamllint > /dev/null
/usr/bin/yamllint $@
`)
	case "yq":
		interactive = "-i"
		entrypoint = strings.ReplaceAll(fmt.Sprintf("%s\\docker-entrypoint.sh", temp), "\\", "/")
		content = []byte(`#!/bin/sh

/sbin/apk add yq > /dev/null
/usr/bin/yq $@
`)
	case "zcat":
		prefix = "/bin"
	}

	binary = fmt.Sprintf("%s/%s", prefix, binary)
	args := os.Args[1:]
	pwd, _ := os.Getwd()
	pwd = strings.ReplaceAll(strings.ReplaceAll(pwd, "\\", "/"), ":", "")
	host := fmt.Sprintf("%s:\\", string(pwd[0]))
	container := fmt.Sprintf("/%s", string(pwd[0]))

	data := pwd[2:]

	work := fmt.Sprintf("%s/%s", container, data)

	docker := []string{
		"run",
		"--rm",
		interactive,
		"-v",
		fmt.Sprintf("%s:%s", host, container),
		"-w",
		work,
	}

	if entrypoint != "" {
		os.Remove(entrypoint)

		file, err := os.Create(entrypoint)

		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()

		if _, err = file.Write(content); err != nil {
			fmt.Println(err)
			return
		}

		docker = append(docker, "-v")
		docker = append(docker, fmt.Sprintf("%s:/docker-entrypoint.sh", entrypoint))
		docker = append(docker, "--entrypoint")
		docker = append(docker, "/docker-entrypoint.sh")
	}

	docker = append(docker, "alpine:latest")

	if entrypoint == "" {
		docker = append(docker, binary)
	}

	if len(args) > 0 {
		docker = append(docker, args...)
	}

	cmd := exec.Command("docker", docker...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n\033[1;31m%s: \033[1;33m%s\033[0m\n", "An error occurred", err)
		os.Exit(1)
	}

	os.Exit(0)
}
