package developer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/buildinfo"
)

type recordedCommand struct {
	environment []string
	executable  string
	arguments   []string
}

type recordingRunner struct {
	commands []recordedCommand
	failOn   int
}

func (runner *recordingRunner) Run(_ context.Context, _ string, environment []string, _ io.Writer, _ io.Writer, executable string, arguments ...string) error {
	runner.commands = append(runner.commands, recordedCommand{
		environment: append([]string(nil), environment...),
		executable:  executable,
		arguments:   append([]string(nil), arguments...),
	})
	if runner.failOn > 0 && len(runner.commands) == runner.failOn {
		return errors.New("synthetic command failure")
	}
	return nil
}

func TestHelpAndVersionAreDeterministic(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run(context.Background(), Options{Arguments: []string{"help"}, Stdout: &stdout, Stderr: &stderr})
	if code != exitOK {
		t.Fatalf("help code = %d, stderr = %q", code, stderr.String())
	}
	if stdout.String() != helpText {
		t.Fatalf("help output = %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Run(context.Background(), Options{Arguments: []string{"version"}, Stdout: &stdout, Stderr: &stderr})
	if code != exitOK {
		t.Fatalf("version code = %d, stderr = %q", code, stderr.String())
	}
	identity := buildinfo.Current()
	want := fmt.Sprintf("omnexa version=%s commit=%s\n", identity.Version, identity.Commit)
	if stdout.String() != want {
		t.Fatalf("version output = %q, want %q", stdout.String(), want)
	}
}

func TestHealthOutputIsSafeAndStructured(t *testing.T) {
	sensitiveFixture := strings.Repeat("x", 24)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run(context.Background(), Options{
		Arguments: []string{"health"},
		Environment: []string{
			"OMNEXA_ENVIRONMENT=ci",
			"OMNEXA_DATABASE_URL=postgres://omnexa:" + sensitiveFixture + "@127.0.0.1:5432/omnexa?sslmode=disable",
		},
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if code != exitOK {
		t.Fatalf("health code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"readiness":"healthy"`) {
		t.Fatalf("health output = %q", stdout.String())
	}
	if strings.Contains(stdout.String(), sensitiveFixture) || strings.Contains(stderr.String(), sensitiveFixture) {
		t.Fatal("health output leaked restricted configuration")
	}
	if strings.Contains(stdout.String(), "postgres://") {
		t.Fatal("health output leaked database resource identity")
	}
}

func TestDatabaseMigrationRequiresExplicitDeveloperEnvironment(t *testing.T) {
	sensitiveFixture := strings.Repeat("y", 24)
	cases := []struct {
		name        string
		environment []string
	}{
		{
			name: "default environment is not explicit",
			environment: []string{
				"OMNEXA_DATABASE_URL=postgres://omnexa:" + sensitiveFixture + "@127.0.0.1:5432/omnexa?sslmode=disable",
			},
		},
		{
			name: "production is forbidden",
			environment: []string{
				"OMNEXA_ENVIRONMENT=production",
				"OMNEXA_DATABASE_URL=postgres://omnexa:" + sensitiveFixture + "@127.0.0.1:5432/omnexa?sslmode=disable",
			},
		},
		{
			name: "staging is outside developer migration authority",
			environment: []string{
				"OMNEXA_ENVIRONMENT=staging",
				"OMNEXA_DATABASE_URL=postgres://omnexa:" + sensitiveFixture + "@127.0.0.1:5432/omnexa?sslmode=disable",
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := Run(context.Background(), Options{
				Arguments:   []string{"db", "migrate"},
				Environment: testCase.environment,
				Stdout:      &stdout,
				Stderr:      &stderr,
			})
			if code != exitOperationFailed {
				t.Fatalf("db migrate code = %d, stderr = %q", code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
			if strings.Contains(stderr.String(), sensitiveFixture) || strings.Contains(stderr.String(), "postgres://") {
				t.Fatal("db migrate failure leaked restricted configuration")
			}
		})
	}
}

func TestVerifyAllUsesCanonicalOrderedOperations(t *testing.T) {
	root := createRepositoryMarkers(t)
	runner := &recordingRunner{}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run(context.Background(), Options{
		Arguments: []string{"verify", "all"},
		Environment: []string{
			"PATH=/usr/bin",
			"P01_04_TEST_DATABASE_URL=synthetic-database",
			"OMNEXA_ENVIRONMENT=ci",
			"OMNEXA_DATABASE_URL=restricted-database",
		},
		Stdout:  &stdout,
		Stderr:  &stderr,
		WorkDir: filepath.Join(root, "kernel", "cmd"),
		Runner:  runner,
	})
	if code != exitOK {
		t.Fatalf("verify all code = %d, stderr = %q", code, stderr.String())
	}
	if len(runner.commands) != 20 {
		t.Fatalf("verify all command count = %d, want 20", len(runner.commands))
	}
	assertCommand(t, runner.commands[0], "python", "scripts/validate_governance.py")
	assertCommand(t, runner.commands[6], "bash", "scripts/verify_go_quality.sh")
	assertCommand(t, runner.commands[7], "bash", "scripts/verify_p01_01.sh")
	assertCommand(t, runner.commands[17], "bash", "scripts/verify_p01_11.sh")
	assertCommand(t, runner.commands[18], "go", "mod", "verify")
	assertCommand(t, runner.commands[19], "go", "build", "./kernel/...")
	assertEnvironment(t, runner.commands[0].environment, "PATH=/usr/bin", "P01_04_TEST_DATABASE_URL=synthetic-database")
	if !strings.Contains(stdout.String(), "omnexa verify all: PASS") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestVerifyStopsAtFirstFailureAndReturnsNonZero(t *testing.T) {
	root := createRepositoryMarkers(t)
	runner := &recordingRunner{failOn: 2}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run(context.Background(), Options{
		Arguments: []string{"verify", "governance"},
		Stdout:    &stdout,
		Stderr:    &stderr,
		WorkDir:   root,
		Runner:    runner,
	})
	if code != exitOperationFailed {
		t.Fatalf("verify failure code = %d, want %d", code, exitOperationFailed)
	}
	if len(runner.commands) != 2 {
		t.Fatalf("executed %d commands after failure, want 2", len(runner.commands))
	}
	if !strings.Contains(stderr.String(), "omnexa verify governance: FAIL") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestModuleLifecycleIsExplicitlyNotApplicableDuringP01(t *testing.T) {
	root := createRepositoryMarkers(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run(context.Background(), Options{
		Arguments: []string{"verify", "module-lifecycle"},
		Stdout:    &stdout,
		Stderr:    &stderr,
		WorkDir:   root,
	})
	if code != exitOK {
		t.Fatalf("module lifecycle code = %d, stderr = %q", code, stderr.String())
	}
	if stdout.String() != "omnexa verify module-lifecycle: N/A\n" {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestUnknownCommandAndVerifyTargetFailClosed(t *testing.T) {
	root := createRepositoryMarkers(t)
	for _, arguments := range [][]string{{"admin"}, {"verify", "everything"}} {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := Run(context.Background(), Options{
			Arguments: arguments,
			Stdout:    &stdout,
			Stderr:    &stderr,
			WorkDir:   root,
		})
		if code == exitOK {
			t.Fatalf("arguments %v unexpectedly succeeded", arguments)
		}
		if stdout.Len() != 0 {
			t.Fatalf("arguments %v stdout = %q", arguments, stdout.String())
		}
	}
}

func TestOSCommandRunnerRejectsUnapprovedExecutable(t *testing.T) {
	err := (OSCommandRunner{}).Run(context.Background(), t.TempDir(), nil, io.Discard, io.Discard, "sh", "-c", "echo unsafe")
	if !errors.Is(err, errUnsupportedExecutable) {
		t.Fatalf("Run() error = %v, want unsupported executable", err)
	}
}

func createRepositoryMarkers(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o750); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "kernel", "cmd"), 0o750); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	for path, content := range map[string]string{
		filepath.Join(root, ".go-version"):                       "1.26.7\n",
		filepath.Join(root, "go.work"):                           "go 1.26.0\n",
		filepath.Join(root, "scripts", "verify_go_quality.sh"): "#!/usr/bin/env bash\n",
	} {
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("WriteFile(%q) error = %v", path, err)
		}
	}
	return root
}

func assertCommand(t *testing.T, command recordedCommand, executable string, arguments ...string) {
	t.Helper()
	if command.executable != executable {
		t.Fatalf("executable = %q, want %q", command.executable, executable)
	}
	if strings.Join(command.arguments, "\x00") != strings.Join(arguments, "\x00") {
		t.Fatalf("arguments = %q, want %q", command.arguments, arguments)
	}
}

func assertEnvironment(t *testing.T, environment []string, expected ...string) {
	t.Helper()
	if strings.Join(environment, "\x00") != strings.Join(expected, "\x00") {
		t.Fatalf("environment = %q, want %q", environment, expected)
	}
	for _, item := range environment {
		if strings.HasPrefix(item, "OMNEXA_") {
			t.Fatalf("verification environment leaked runtime configuration: %q", item)
		}
	}
}
