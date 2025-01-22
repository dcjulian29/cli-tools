package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

func main() {
	action := ""

	if len(os.Args) > 1 {
		action = os.Args[1]
		action = strings.ReplaceAll(action, "-", "")
	} else {
		if fileExists("ansible.cfg") {
			action = "ansible"
		}

		if fileExists("dockerfile") {
			action = "docker"
		}

		if fileExists("go.mod") {
			action = "go"
		}

		if fileExists(".goreleaser.yml") || fileExists(".goreleaser.yaml") {
			action = "goreleaser"
		}

		if fileExists("build.sh") && runtime.GOOS != "windows" {
			action = "sh"
		}

		if fileExists("build.bat") && runtime.GOOS == "windows" {
			action = "bat"
		}

		if fileExists("build.cmd") && runtime.GOOS == "windows" {
			action = "cmd"
		}

		if fileExists("build.ps1") {
			action = "ps"
		}

		if fileExists("build.cake") {
			action = "cake"
		}
	}

	switch action {
	case "archive":
		archive()
		os.Exit(0)
	case "cake":
		buildCake()
		os.Exit(0)
	case "ps":
	case "pshell":
	case "powershell":
		buildPowershell()
		os.Exit(0)
	case "bat":
		buildDos(true)
		os.Exit(0)
	case "cmd":
		buildDos(false)
		os.Exit(0)
	case "sh":
		buildBash()
		os.Exit(0)
	case "goreleaser":
		buildGoReleaser()
		os.Exit(0)
	case "go":
		buildGo()
		os.Exit(0)
	case "ansible":
		buildAnsible()
		os.Exit(0)
	case "docker":
		buildDocker()
		os.Exit(0)
	case "":
	default:
		fmt.Println(color.RedString("Nothing found to build!"))
		os.Exit(1)
	}
}

func archive() {
	pwd, _ := os.Getwd()
	name := filepath.Base(pwd)

	dst := fmt.Sprintf("../%s.7z", name)

	fmt.Printf("Archiving '%s'...\n", pwd)

	run("7z", []string{"a", "-t7z", "-mx9", "-y", "-r", dst, "."})
}

func buildAnsible() {
	run("ansible-lint", []string{"."})
}

func buildBash() {
	run("bash", []string{"build.sh"})
}

func buildCake() {
	cmd := exec.Command("dotnet", "tool", "list")
	tools, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(color.RedString("dotnet SDK is not present!s"))
		os.Exit(2)
	}

	if !strings.Contains(string(tools), "cake.tool") {
		if !fileExists(".config/dotnet-tools.json") {
			cmd := exec.Command("dotnet", "new", "tool-manifest")
			_, err := cmd.CombinedOutput()

			if err != nil {
				fmt.Println(color.RedString("Installing Cake.Tool: %s", err))
				os.Exit(3)
			}

			err = run("dotnet", []string{"tool", "install", "Cake.Tool"})

			if err != nil {
				fmt.Println(color.RedString("Cake.Tool is not present and could not be installed!"))
				os.Exit(4)
			}
		}
	}

	var params []string

	if len(os.Args) > 0 {
		if os.Args[0] == "cake" {
			if !strings.Contains("-", os.Args[1]) {
				params = []string{"--target=" + os.Args[1]}
			} else {
				params = os.Args[1:]
			}
		} else {
			params = os.Args
		}
	}

	params = append([]string{"cake"}, params...)
	run("dotnet", params)
}

func buildDocker() {
	run("docker", []string{"build", "."})
}

func buildDos(batchFile bool) {
	if batchFile {
		run("cmd.exe", []string{"/C", "build.bat"})
	} else {
		run("cmd.exe", []string{"/C", "build.cmd"})
	}
}

func buildGoReleaser() {
	run("goreleaser", []string{"release", "--snapshot", "--clean"})
}

func buildGo() {
	run("go", []string{"mod", "tidy"})
	run("go", []string{"build", "-a", "."})
}

func buildPowershell() {
	run("pwsh", []string{"-f", "build.ps1"})
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func run(binary string, params []string) error {
	cmd := exec.Command(binary, params...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
