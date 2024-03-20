package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func main() {
	action := "build"

	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	if strings.EqualFold(action, "archive") {
		archive()
		os.Exit(0)
	}

	if fileExists("build.cake") {
		buildCake()
		os.Exit(0)
	}

	if fileExists("build.ps1") {
		buildPowershell()
		os.Exit(0)
	}

	if fileExists("build.bat") {
		buildDos(true)
		os.Exit(0)
	}

	if fileExists("build.cmd") {
		buildDos(false)
		os.Exit(0)
	}

	if fileExists("build.sh") {
		buildBash()
		os.Exit(0)
	}

	if fileExists(".goreleaser.yml") || fileExists(".goreleaser.yaml") {
		buildGoReleaser()
		os.Exit(0)
	}

	if fileExists("go.mod") {
		buildGo()
		os.Exit(0)
	}

	if fileExists("ansible.cfg") {
		buildAnsible()
		os.Exit(0)
	}

	if fileExists("dockerfile") {
		buildDocker()
		os.Exit(0)
	}

	fmt.Println(color.RedString("Nothing found to build!"))
}

func archive() {
	// To Do
}

func buildAnsible() {
	run("ansible-lint", ".")
}

func buildBash() {
	run("bash", "build.sh")
}

func buildCake() {
	if len(os.Args) == 1 {
		run("dotnet", "cake", "--target="+os.Args[0])
	}

	var params = strings.Join(os.Args, " ")

	run("dotnet", "cake", params)
}

func buildDocker() {
	run("docker", "build", ".")
}

func buildDos(batchFile bool) {
	if batchFile {
		run("cmd.exe", "/C", "build.bat")
	} else {
		run("cmd.exe", "/C", "build.cmd")
	}
}

func buildGoReleaser() {
	run("goreleaser", "release", "--snapshot", "--clean")
}

func buildGo() {
	run("go", "mod", "tidy")
	run("go", "build", "-a", ".")
}

func buildPowershell() {
	run("pwsh", "-f", "build.ps1")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func run(binary string, params ...string) error {
	cmd := exec.Command(binary, params...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
