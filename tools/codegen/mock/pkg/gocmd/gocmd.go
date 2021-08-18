// Package gocmd provides access to some go commands functions.
package gocmd

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type command string

func (c command) exec(ctx context.Context, args ...string) *exec.Cmd {
	args = append([]string{string(c)}, args...)

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Stderr = os.Stderr

	return cmd
}

// =============================================================================
// mod
// =====================================================================================

const mod command = "mod"

// DownloadModules downloads all modules in the given dir.
// If dir is empty, the current working directory will be used.
func DownloadModules(ctx context.Context, dir string) error {
	cmd := mod.exec(ctx, "download")
	cmd.Dir = dir

	return errors.Wrap(cmd.Run(), "gocmd")
}

type GoMod struct {
	Module  Module
	Go      string
	Require []Require
	Exclude []Module
	Replace []Replace
	Retract []Retract
}

type Module struct {
	Path    string
	Version string
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

type Replace struct {
	Old Module
	New Module
}

type Retract struct {
	Low       string
	High      string
	Rationale string
}

// ModFile returns an ast of the go.mod file found in the given dir.
// The passed dir doesn't have to be the root directory of the module, but may
// also be a child directory of the module.
// If dir is empty, the current working directory will be used.
func ModFile(dir string) (*GoMod, error) {
	cmd := mod.exec(context.Background(), "edit", "-json")
	cmd.Dir = dir

	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var m *GoMod
	return m, errors.Wrap(json.Unmarshal(data, &m), "gocmd")
}

// =============================================================================
// env
// =====================================================================================

const env command = "env"

type GoEnv string

const (
	ModCacheEnv GoEnv = "GOMODCACHE"
)

// Env retrieves the given GoEnv or returns its default.
func Env(name GoEnv) (string, error) {
	cmd := env.exec(context.Background(), string(name))

	data, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "gocmd")
	}

	return strings.TrimSpace(string(data)), nil
}
