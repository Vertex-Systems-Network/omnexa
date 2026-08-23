// Package developer implements the bounded P01.12 repository-owned developer CLI.
package developer

import (
	"context"
	"errors"
	"io"
	"os/exec"
)

var errUnsupportedExecutable = errors.New("developer command executable is not approved")

// CommandRunner executes one allowlisted repository command without invoking a shell.
// Implementations must not add hidden command expansion or mutate repository state.
type CommandRunner interface {
	Run(ctx context.Context, directory string, environment []string, stdout io.Writer, stderr io.Writer, executable string, arguments ...string) error
}

// OSCommandRunner executes the bounded set of repository tooling used by P01.12.
type OSCommandRunner struct{}

// Run executes an approved command directly. The executable is selected from a
// fixed allowlist so repository input can never become a shell command name.
func (OSCommandRunner) Run(ctx context.Context, directory string, environment []string, stdout io.Writer, stderr io.Writer, executable string, arguments ...string) error {
	var command *exec.Cmd
	switch executable {
	case "bash":
		command = exec.CommandContext(ctx, "bash", arguments...)
	case "python":
		command = exec.CommandContext(ctx, "python", arguments...)
	case "go":
		command = exec.CommandContext(ctx, "go", arguments...)
	default:
		return errUnsupportedExecutable
	}
	command.Dir = directory
	command.Env = environment
	command.Stdout = stdout
	command.Stderr = stderr
	return command.Run()
}
