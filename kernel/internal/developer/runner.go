// Package developer implements the bounded P01.12 repository-owned developer CLI.
package developer

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"strings"
)

var errUnsupportedExecutable = errors.New("developer command executable or arguments are not approved")

// CommandRunner executes one allowlisted repository command without invoking a shell.
// Implementations must not add hidden command expansion or mutate repository state.
type CommandRunner interface {
	Run(ctx context.Context, directory string, environment []string, stdout io.Writer, stderr io.Writer, executable string, arguments ...string) error
}

// OSCommandRunner executes the bounded set of repository tooling used by P01.12.
type OSCommandRunner struct{}

// Run executes an exact approved command directly. Both executable and arguments
// are matched against a fixed allowlist before a literal command is constructed.
func (OSCommandRunner) Run(ctx context.Context, directory string, environment []string, stdout io.Writer, stderr io.Writer, executable string, arguments ...string) error {
	command, err := approvedCommand(ctx, executable, arguments)
	if err != nil {
		return err
	}
	command.Dir = directory
	command.Env = environment
	command.Stdout = stdout
	command.Stderr = stderr
	return command.Run()
}

func approvedCommand(ctx context.Context, executable string, arguments []string) (*exec.Cmd, error) {
	key := executable + "\x00" + strings.Join(arguments, "\x00")
	switch key {
	case "python\x00scripts/validate_governance.py":
		return exec.CommandContext(ctx, "python", "scripts/validate_governance.py"), nil
	case "python\x00scripts/validate_development_spec.py":
		return exec.CommandContext(ctx, "python", "scripts/validate_development_spec.py"), nil
	case "python\x00scripts/validate_operations_spec.py":
		return exec.CommandContext(ctx, "python", "scripts/validate_operations_spec.py"), nil
	case "python\x00scripts/validate_freeze_review.py":
		return exec.CommandContext(ctx, "python", "scripts/validate_freeze_review.py"), nil
	case "python\x00scripts/validate_p01_preparation.py":
		return exec.CommandContext(ctx, "python", "scripts/validate_p01_preparation.py"), nil
	case "python\x00scripts/validate_p01_package_specs.py":
		return exec.CommandContext(ctx, "python", "scripts/validate_p01_package_specs.py"), nil
	case "bash\x00scripts/verify_go_quality.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_go_quality.sh"), nil
	case "bash\x00scripts/verify_p01_01.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_01.sh"), nil
	case "bash\x00scripts/verify_p01_02.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_02.sh"), nil
	case "bash\x00scripts/verify_p01_03.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_03.sh"), nil
	case "bash\x00scripts/verify_p01_04.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_04.sh"), nil
	case "bash\x00scripts/verify_p01_05.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_05.sh"), nil
	case "bash\x00scripts/verify_p01_06.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_06.sh"), nil
	case "bash\x00scripts/verify_p01_07.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_07.sh"), nil
	case "bash\x00scripts/verify_p01_08.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_08.sh"), nil
	case "bash\x00scripts/verify_p01_09.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_09.sh"), nil
	case "bash\x00scripts/verify_p01_10.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_10.sh"), nil
	case "bash\x00scripts/verify_p01_11.sh":
		return exec.CommandContext(ctx, "bash", "scripts/verify_p01_11.sh"), nil
	case "go\x00test\x00./kernel/...":
		return exec.CommandContext(ctx, "go", "test", "./kernel/..."), nil
	case "go\x00build\x00./kernel/...":
		return exec.CommandContext(ctx, "go", "build", "./kernel/..."), nil
	case "go\x00mod\x00verify":
		return exec.CommandContext(ctx, "go", "mod", "verify"), nil
	default:
		return nil, errUnsupportedExecutable
	}
}
