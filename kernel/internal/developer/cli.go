package developer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/buildinfo"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/database"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/operations"
)

const (
	exitOK              = 0
	exitIOFailure       = 1
	exitUsage           = 2
	exitOperationFailed = 3
)

const helpText = `Omnexa developer CLI
usage: omnexa <command>

commands:
  help                      show this help
  version                   show deterministic build identity
  health                    show safe JSON health/readiness diagnostics
  db migrate                run the bounded non-production kernel migration boundary
  verify <target>           run governed verification target

verify targets:
  governance format lint static unit contracts integration migrations
  security module-lifecycle build release all
`

var errRepositoryRootNotFound = errors.New("Omnexa repository root could not be found")

// Options supplies explicit process inputs to the P01.12 developer CLI.
type Options struct {
	Arguments   []string
	Environment []string
	Stdout      io.Writer
	Stderr      io.Writer
	WorkDir     string
	Runner      CommandRunner
}

// Run executes one bounded developer CLI command and returns a stable process exit code.
func Run(ctx context.Context, options Options) int {
	if ctx == nil {
		ctx = context.Background()
	}
	options = withDefaults(options)
	if len(options.Arguments) == 0 {
		return writeFailure(options.Stderr, "omnexa: command is required", exitUsage)
	}

	switch options.Arguments[0] {
	case "help", "--help", "-h":
		return runHelp(options.Stdout, options.Stderr)
	case "version", "--version":
		return runVersion(options.Stdout, options.Stderr)
	case "health":
		return runHealth(ctx, options)
	case "db":
		return runDatabase(ctx, options)
	case "verify":
		return runVerify(ctx, options)
	default:
		return writeFailure(options.Stderr, "omnexa: unknown command", exitUsage)
	}
}

func withDefaults(options Options) Options {
	if options.Environment == nil {
		options.Environment = os.Environ()
	}
	if options.Stdout == nil {
		options.Stdout = io.Discard
	}
	if options.Stderr == nil {
		options.Stderr = io.Discard
	}
	if options.Runner == nil {
		options.Runner = OSCommandRunner{}
	}
	return options
}

func runHelp(stdout, stderr io.Writer) int {
	if _, err := io.WriteString(stdout, helpText); err != nil {
		return writeFailure(stderr, "omnexa: output failed", exitIOFailure)
	}
	return exitOK
}

func runVersion(stdout, stderr io.Writer) int {
	identity := buildinfo.Current()
	if _, err := fmt.Fprintf(stdout, "omnexa version=%s commit=%s\n", identity.Version, identity.Commit); err != nil {
		return writeFailure(stderr, "omnexa: output failed", exitIOFailure)
	}
	return exitOK
}

func runHealth(ctx context.Context, options Options) int {
	if len(options.Arguments) != 1 {
		return writeFailure(options.Stderr, "omnexa health: unexpected arguments", exitUsage)
	}
	if _, err := database.LoadConfiguration(config.OSOptions(options.Environment)); err != nil {
		return writeFailure(options.Stderr, "omnexa health: configuration is invalid", exitOperationFailed)
	}

	manager := operations.NewManager(nil)
	if !manager.MarkReady() {
		return writeFailure(options.Stderr, "omnexa health: lifecycle transition failed", exitOperationFailed)
	}
	report := manager.Evaluate(ctx)
	encoder := json.NewEncoder(options.Stdout)
	encoder.SetEscapeHTML(true)
	if err := encoder.Encode(report); err != nil {
		return writeFailure(options.Stderr, "omnexa health: output failed", exitIOFailure)
	}
	return exitOK
}

func runDatabase(ctx context.Context, options Options) int {
	if len(options.Arguments) != 2 || options.Arguments[1] != "migrate" {
		return writeFailure(options.Stderr, "omnexa db: supported command is migrate", exitUsage)
	}
	resolved, err := database.LoadConfiguration(config.OSOptions(options.Environment))
	if err != nil {
		return writeFailure(options.Stderr, "omnexa db migrate: configuration is invalid", exitOperationFailed)
	}
	environment, ok := resolved.Environment("environment")
	if !ok || !explicitEnvironment(resolved) || !developerMigrationEnvironment(environment) {
		return writeFailure(options.Stderr, "omnexa db migrate: explicit non-production environment is required", exitOperationFailed)
	}

	pool, err := database.NewPool(ctx, resolved)
	if err != nil {
		return writeFailure(options.Stderr, "omnexa db migrate: database is unavailable", exitOperationFailed)
	}
	defer pool.Close()

	migrator, err := database.NewMigrator(pool, "kernel.foundation", nil, 5*time.Second)
	if err != nil {
		return writeFailure(options.Stderr, "omnexa db migrate: migration boundary is invalid", exitOperationFailed)
	}
	if err := migrator.Run(ctx); err != nil {
		return writeFailure(options.Stderr, "omnexa db migrate: migration failed", exitOperationFailed)
	}
	if _, err := fmt.Fprintf(options.Stdout, "omnexa db migrate: PASS environment=%s\n", environment); err != nil {
		return writeFailure(options.Stderr, "omnexa db migrate: output failed", exitIOFailure)
	}
	return exitOK
}

func runVerify(ctx context.Context, options Options) int {
	if len(options.Arguments) != 2 {
		return writeFailure(options.Stderr, "omnexa verify: target is required", exitUsage)
	}
	root, err := findRepositoryRoot(options.WorkDir)
	if err != nil {
		return writeFailure(options.Stderr, "omnexa verify: repository root not found", exitOperationFailed)
	}
	return runVerification(ctx, root, options.Environment, options.Stdout, options.Stderr, options.Runner, options.Arguments[1])
}

func explicitEnvironment(resolved config.Config) bool {
	for _, provenance := range resolved.Provenance() {
		if provenance.Key == "environment" {
			return provenance.Source != config.SourceDefault
		}
	}
	return false
}

func developerMigrationEnvironment(environment config.Environment) bool {
	switch environment {
	case config.EnvironmentLocal, config.EnvironmentCI, config.EnvironmentPreview, config.EnvironmentTest:
		return true
	case config.EnvironmentStaging, config.EnvironmentProduction:
		return false
	default:
		return false
	}
}

func findRepositoryRoot(start string) (string, error) {
	if start == "" {
		current, err := os.Getwd()
		if err != nil {
			return "", errRepositoryRootNotFound
		}
		start = current
	}
	absolute, err := filepath.Abs(start)
	if err != nil {
		return "", errRepositoryRootNotFound
	}
	for directory := absolute; ; directory = filepath.Dir(directory) {
		if repositoryMarkersPresent(directory) {
			return directory, nil
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			break
		}
	}
	return "", errRepositoryRootNotFound
}

func repositoryMarkersPresent(directory string) bool {
	markers := []string{
		filepath.Join(directory, ".go-version"),
		filepath.Join(directory, "go.work"),
		filepath.Join(directory, "scripts", "verify_go_quality.sh"),
	}
	for _, marker := range markers {
		info, err := os.Stat(marker)
		if err != nil || info.IsDir() {
			return false
		}
	}
	return true
}

func writeFailure(stderr io.Writer, message string, code int) int {
	if _, err := fmt.Fprintln(stderr, message); err != nil {
		return exitIOFailure
	}
	return code
}
