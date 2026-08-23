package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/buildinfo"
)

func TestRunWithDefaultConfigurationPreservesKernelOutput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(&stdout, &stderr, nil)
	if code != 0 {
		t.Fatalf("run() code = %d, stderr = %q", code, stderr.String())
	}
	want := fmt.Sprintf("omnexa-kernel %s\n", buildinfo.Current())
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunCommandWithoutArgumentsPreservesKernelStartup(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCommand(context.Background(), nil, &stdout, &stderr, nil)
	if code != 0 {
		t.Fatalf("runCommand() code = %d, stderr = %q", code, stderr.String())
	}
	want := fmt.Sprintf("omnexa-kernel %s\n", buildinfo.Current())
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunCommandExposesHelpAndVersion(t *testing.T) {
	for _, arguments := range [][]string{{"help"}, {"version"}} {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := runCommand(context.Background(), arguments, &stdout, &stderr, nil)
		if code != 0 {
			t.Fatalf("arguments %v code = %d, stderr = %q", arguments, code, stderr.String())
		}
		if stdout.Len() == 0 {
			t.Fatalf("arguments %v produced empty stdout", arguments)
		}
		if stderr.Len() != 0 {
			t.Fatalf("arguments %v stderr = %q", arguments, stderr.String())
		}
	}
}

func TestRunCommandUnknownCommandFailsClosed(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCommand(context.Background(), []string{"admin"}, &stdout, &stderr, nil)
	if code == 0 {
		t.Fatal("runCommand() code = 0, want failure")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunUsesExplicitConfigFileThenEnvironmentOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"environment":"ci"}`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run(&stdout, &stderr, []string{
		"OMNEXA_CONFIG_FILE=" + path,
		"OMNEXA_ENVIRONMENT=test",
	})
	if code != 0 {
		t.Fatalf("run() code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.HasPrefix(stdout.String(), "omnexa-kernel ") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunFailsClosedForInvalidEnvironmentWithoutEchoingValue(t *testing.T) {
	const raw = "production-secret-like-value"
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(&stdout, &stderr, []string{"OMNEXA_ENVIRONMENT=" + raw})
	if code == 0 {
		t.Fatal("run() code = 0, want failure")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "omnexa configuration error") {
		t.Fatalf("stderr = %q, want configuration error", stderr.String())
	}
	if strings.Contains(stderr.String(), raw) {
		t.Fatalf("stderr leaked raw invalid value: %q", stderr.String())
	}
}

func TestRunFailsClosedForUnknownOmnexaVariable(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(&stdout, &stderr, []string{"OMNEXA_ENVIRONMNET=production"})
	if code == 0 {
		t.Fatal("run() code = 0, want unknown configuration failure")
	}
	if !strings.Contains(stderr.String(), "unknown Omnexa environment configuration key") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}
