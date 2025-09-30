//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = BuildAll

var Aliases = map[string]interface{}{
	"i":   InstallAll,
	"ia":  InstallAPI,
	"ib":  InstallBot,
	"ic":  InstallCLI,
	"iad": InstallAllDeps,
	"b":   BuildAll,
	"ba":  BuildAPI,
	"bb":  BuildBot,
	"bc":  BuildCLI,
	"bf":  BuildFrontend,
}

var (
	toolsDir     = "tools"
	rustupHome   = filepath.Join(toolsDir, "rustup")
	cargoHome    = filepath.Join(toolsDir, "cargo")
	cargoBinDir  = filepath.Join(cargoHome, "bin")
	buildDir     = "build"
	webStaticDir = filepath.Join(buildDir, "static")
	binDir       = filepath.Join(buildDir, "bin")
)

func CheckRoot() error {
	var CurrentUser, _ = user.Current()
	if CurrentUser.Username == "root" {
		println("It's recommended to compile as non-root user.")
	}
	return nil
}

func BuildAPI() error {
	CheckRoot()
	fmt.Println("Building API...")
	cmd := exec.Command("go", "build", "-o", filepath.Join(binDir, "xyter-api"), "./cmd/xyter-api");
	cmd.Env = append(os.Environ(),
	  "CGO_ENABLED=0",
        )
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func BuildCLI() error {
	CheckRoot()
	fmt.Println("Building CLI...")
	cmd := exec.Command("go", "build", "-o", filepath.Join(binDir, "xyter-cli"), "./cmd/xyter-cli");
	cmd.Env = append(os.Environ(),
	  "CGO_ENABLED=0",
        )
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func BuildBot() error {
	CheckRoot()
	fmt.Println("Building Discord Bot...")
	cmd := exec.Command("go", "build", "-o", filepath.Join(binDir, "xyter-bot"), "./cmd/xyter-bot");
	cmd.Env = append(os.Environ(),
	  "CGO_ENABLED=0",
        )
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func BuildFrontend() error {
	CheckRoot()
	mg.SerialDeps(InstallRust, InstallFrontendDeps)
	cmd := exec.Command("trunk", "build", "--release")
	cmd.Env = RustEnv()
	cmd.Dir = "frontend"
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// A build step that requires additional params, or platform specific steps for example
func BuildAll() error {
	if err := os.MkdirAll("bin", 0755); err != nil {
		return err
	}

	fmt.Println("Building everything...")
	mg.Deps(InstallAllDeps)
	mg.Deps(BuildAPI)
	mg.Deps(BuildBot)
	mg.Deps(BuildCLI)
	mg.Deps(BuildFrontend)

	return nil
}

func InstallAPI() error {
	mg.Deps(BuildAPI)
	if err := os.Rename("./bin/xyter-api", "/usr/bin/xyter-api"); err != nil {
		return err
	}
	return nil
}

func InstallBot() error {
	mg.Deps(BuildBot)
	if err := os.Rename("./bin/xyter-cli", "/usr/bin/xyter-cli"); err != nil {
		return err
	}
	return nil
}

func InstallCLI() error {
	mg.Deps(BuildCLI)
	if err := os.Rename("./bin/xyter-bot", "/usr/bin/xyter-bot"); err != nil {
		return err
	}
	return nil
}

func InstallFrontend() error {
	mg.Deps(BuildFrontend)
	if err := os.Rename("./frontend/dist", "./build/static"); err != nil {
		return err
	}
	return nil
}

// A custom install step if you need your bin someplace other than go/bin
func InstallAll() error {
	mg.Deps(BuildAll)
	fmt.Println("Installing...")
	mg.Deps(InstallAPI, InstallBot, InstallCLI)
	return nil
}

// Manage your deps, or running package managers.
func InstallAllDeps() error {
	fmt.Println("Installing Deps...")
	return nil
}
func InstallFrontendDeps() error {
	fmt.Println("Installing Deps for frontend...")
	cmd := exec.Command("cargo", "install", "trunk")
	cmd.Env = RustEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InstallRust installs rustup locally if not installed already.
func InstallRust() error {
	fmt.Println("Installing local Rust toolchain with rustup...")

	// Download rustup-init script or executable URL
	var rustupURL string
	switch runtime.GOOS {
	case "linux":
		rustupURL = "https://sh.rustup.rs"
	case "darwin":
		rustupURL = "https://sh.rustup.rs"
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	// Create tools directory
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		return err
	}

	// Run rustup-init with env vars to install locally
	cmd := exec.Command("sh", "-c", "curl --proto '=https' --tlsv1.2 -sSf "+rustupURL+" | sh -s -- -y --no-modify-path --default-toolchain stable")
	cmd.Env = RustEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RustEnv returns environment variables to use the local Rust toolchain
func RustEnv() []string {
	env := os.Environ()
	env = append(env,
		"RUSTUP_HOME="+rustupHome,
		"CARGO_HOME="+cargoHome,
		"PATH="+cargoBinDir+string(os.PathListSeparator)+os.Getenv("PATH"),
	)
	return env
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("build")
	os.RemoveAll("tools")
}
