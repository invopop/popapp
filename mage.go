//go:build mage

package main

import (
	"os"

	"github.com/invopop/tasks"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"golang.org/x/term"
)

// ABOUT: This file contains the mage build system for the popapp service.
// It provides a set of tasks for building, running, and releasing the service.
// For local development you may want to use air to automatically rebuild and restart the service.
// air will automatically call `mage serve` when changes are detected.

const (
	name     = "popapp"
	runImage = "golang:1.25.0-alpine"
)

// Build builds the service
func Build() error {
	if err := sh.RunV("templ", "generate"); err != nil {
		return err
	}
	changed, err := target.Dir("./"+name, ".")
	if os.IsNotExist(err) || (err == nil && changed) {
		args := []string{
			"GOOS=linux",
			"GOARCH=amd64",
			"GO111MODULE=on",
			"CGO_ENABLED=0",
			"go", "build", ".",
		}
		return sh.RunV("env", args...)
	}
	return nil
}

// Serve starts the service
func Serve() error {
	mg.Deps(Build)
	return dockerRunCmd(name, "8080", "/src/"+name, "serve")
}

// Shell runs an interactive shell within a docker container.
func Shell() error {
	return dockerRunCmd(name+"-shell", "", "sh")
}

// Release tries to tag the current branch with a timestamp
// so it will be built and released.
func Release() error {
	return tasks.Release()
}

func dockerRunCmd(name, publicPort string, cmd ...string) error {
	call, args := dockerCmdPrep(name, publicPort, cmd...)
	return sh.RunV(call, args...)
}

func dockerCmdPrep(name, publicPort string, cmd ...string) (string, []string) {
	args := []string{
		"run",
		"--rm",
		"--name", name,
		"--network", "invopop-local",
		"-v", "$PWD:/src",
		"-w", "/src",
	}
	if term.IsTerminal(int(os.Stdin.Fd())) {
		args = append(args, "-it")
	}
	if publicPort != "" {
		args = append(args,
			"--label", "traefik.enable=true",
			"--label", "traefik.http.routers."+name+".rule=Host(`"+name+".invopop.dev`)",
			"--label", "traefik.http.routers."+name+".tls=true",
			"--expose", publicPort,
		)
	}
	args = append(args, runImage)
	args = append(args, cmd...)
	return "docker", args
}
