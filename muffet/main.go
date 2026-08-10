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
	"os"

	"github.com/dcjulian29/go-toolbox/docker"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/dcjulian29/go-toolbox/textformat"
)

func main() {
	args := filesystem.EnsureUnixPathArguments()

	data, work, err := docker.HostContainerVolume()
	if err != nil {
		fmt.Println(textformat.Fatal(err.Error()))
		os.Exit(1)
	}

	opts := docker.ContainerOptions{
		AdditionalArgs:   args,
		Image:            "raviqqe/muffet",
		Interactive:      true,
		Tag:              "latest",
		Volumes:          []string{data},
		WorkingDirectory: work,
	}

	if _, err := docker.Run(opts); err != nil {
		fmt.Println(textformat.Fatal(err.Error()))
		os.Exit(1)
	}
}
