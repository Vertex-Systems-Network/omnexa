package database

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestLoadConfigurationAndSettingsRedactDatabaseURL(t *testing.T) {
	const dsn = "postgres://omnexa:restricted-test-secret@127.0.0.1:5432/omnexa?sslmode=disable"
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{
			databaseURLKey:                   dsn,
			databaseConnectTimeoutKey:        "750ms",
			databaseMaxConnectionsKey:        "12",
			databaseMinConnectionsKey:        "2",
			databaseMaxConnectionLifetimeKey: "20m",
			databaseMaxConnectionIdleTimeKey: "3m",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}

	settings, err := SettingsFromConfig(resolved)
	if err != nil {
		t.Fatalf("SettingsFromConfig() error = %v", err)
	}
	if settings.URL != dsn || settings.ConnectTimeout != 750*time.Millisecond || settings.MaxConnections != 12 || settings.MinConnections != 2 {
		t.Fatalf("unexpected pool settings: %#v", settings)
	}
	if settings.MaxConnectionLifetime != 20*time.Minute || settings.MaxConnectionIdleTime != 3*time.Minute {
		t.Fatalf("unexpected connection lifetime settings: %#v", settings)
	}

	redacted := resolved.Redacted()
	if redacted[databaseURLKey] != "<redacted>" {
		t.Fatalf("database_url diagnostic = %q, want redacted", redacted[databaseURLKey])
	}
	if strings.Contains(strings.Join(mapStringValues(redacted), " "), "restricted-test-secret") {
		t.Fatal("redacted configuration leaked database secret")
	}
}

func TestSettingsRequireDatabaseURLWithoutSecretDisclosure(t *testing.T) {
	resolved, err := LoadConfiguration(config.Options{Strict: true})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}
	_, err = SettingsFromConfig(resolved)
	if err == nil {
		t.Fatal("SettingsFromConfig() error = nil, want missing database URL failure")
	}
	assertFailureCode(t, err, codeConfigurationInvalid)
	if strings.Contains(err.Error(), "postgres://") {
		t.Fatalf("error exposed connection data: %q", err)
	}
}

func TestLoadConfigurationRejectsInvalidDatabaseValueWithoutEchoingRawValue(t *testing.T) {
	const raw = "restricted-secret-like-timeout"
	_, err := LoadConfiguration(config.Options{
		Environ: map[string]string{"OMNEXA_DATABASE_CONNECT_TIMEOUT": raw},
		Strict:  true,
	})
	if err == nil {
		t.Fatal("LoadConfiguration() error = nil, want invalid duration failure")
	}
	assertFailureCode(t, err, codeConfigurationInvalid)
	if strings.Contains(err.Error(), raw) {
		t.Fatalf("error leaked raw configuration value: %q", err)
	}
}

func TestSettingsRejectPoolBounds(t *testing.T) {
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{
			databaseURLKey:            "postgres://example.invalid/omnexa",
			databaseMaxConnectionsKey: "2",
			databaseMinConnectionsKey: "3",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}
	_, err = SettingsFromConfig(resolved)
	if err == nil {
		t.Fatal("SettingsFromConfig() error = nil, want invalid bounds failure")
	}
	assertFailureCode(t, err, codeConfigurationInvalid)
}

func TestNewPoolRejectsMalformedURLWithoutPublishingProviderText(t *testing.T) {
	const secret = "restricted-parse-secret"
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{
			databaseURLKey:            "postgres://user:" + secret + "@[]/bad",
			databaseConnectTimeoutKey: "100ms",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = NewPool(ctx, resolved)
	if err == nil {
		t.Fatal("NewPool() error = nil, want malformed URL failure")
	}
	assertFailureCode(t, err, codeConfigurationInvalid)
	if strings.Contains(err.Error(), secret) || strings.Contains(err.Error(), "postgres://") {
		t.Fatalf("safe failure exposed provider connection text: %q", err)
	}
}

func assertFailureCode(t *testing.T, err error, want failure.Code) {
	t.Helper()
	got, ok := failure.CodeOf(err)
	if !ok {
		t.Fatalf("error %T is not a structured failure: %v", err, err)
	}
	if got != want {
		t.Fatalf("failure code = %q, want %q", got, want)
	}
}

func mapStringValues(values map[string]string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}
