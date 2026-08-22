package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestApplicationDefaultsToLocal(t *testing.T) {
	cfg, err := LoadApplication(Options{Strict: true})
	if err != nil {
		t.Fatalf("LoadApplication() error = %v", err)
	}
	environment, ok := cfg.Environment("environment")
	if !ok {
		t.Fatal("environment missing or wrong type")
	}
	if environment != EnvironmentLocal {
		t.Fatalf("environment = %q, want %q", environment, EnvironmentLocal)
	}
	provenance := cfg.Provenance()
	want := []Provenance{{Key: "environment", Source: SourceDefault, Sensitive: false}}
	if !reflect.DeepEqual(provenance, want) {
		t.Fatalf("provenance = %#v, want %#v", provenance, want)
	}
}

func TestDeterministicPrecedenceDefaultFileEnvironmentOverride(t *testing.T) {
	path := writeConfigFile(t, `{"environment":"ci","workers":2,"enabled":false,"timeout":"3s"}`)
	definitions := []Definition{
		{Key: "environment", Kind: KindEnvironment, Default: "local"},
		{Key: "workers", Kind: KindInt, Default: "1"},
		{Key: "enabled", Kind: KindBool, Default: "true"},
		{Key: "timeout", Kind: KindDuration, Default: "1s"},
	}

	cfg, err := Load(definitions, Options{
		FilePath: path,
		Environ: map[string]string{
			"OMNEXA_ENVIRONMENT": "staging",
			"OMNEXA_WORKERS":     "4",
		},
		Overrides: map[string]string{
			"environment": "test",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got, _ := cfg.Environment("environment"); got != EnvironmentTest {
		t.Fatalf("environment = %q, want test", got)
	}
	if got, _ := cfg.Int("workers"); got != 4 {
		t.Fatalf("workers = %d, want 4", got)
	}
	if got, _ := cfg.Bool("enabled"); got != false {
		t.Fatalf("enabled = %v, want false", got)
	}
	if got, _ := cfg.Duration("timeout"); got != 3*time.Second {
		t.Fatalf("timeout = %v, want 3s", got)
	}

	provenance := cfg.Provenance()
	want := []Provenance{
		{Key: "environment", Source: SourceOverride},
		{Key: "workers", Source: SourceEnvironment},
		{Key: "enabled", Source: SourceFile},
		{Key: "timeout", Source: SourceFile},
	}
	if !reflect.DeepEqual(provenance, want) {
		t.Fatalf("provenance = %#v, want %#v", provenance, want)
	}
}

func TestRequiredSettingFailsClosedWithoutValueLeak(t *testing.T) {
	definitions := []Definition{{Key: "api_token", Kind: KindString, Required: true, Sensitive: true}}
	_, err := Load(definitions, Options{Strict: true})
	if err == nil {
		t.Fatal("Load() error = nil, want required-key failure")
	}
	if !strings.Contains(err.Error(), `required configuration key "api_token" is missing`) {
		t.Fatalf("error = %q, want safe missing-key message", err)
	}
}

func TestInvalidSensitiveValueDoesNotLeakRawValue(t *testing.T) {
	const secret = "super-secret-token-value"
	definitions := []Definition{{Key: "secret_timeout", Kind: KindDuration, Required: true, Sensitive: true}}
	_, err := Load(definitions, Options{
		Environ: map[string]string{"OMNEXA_SECRET_TIMEOUT": secret},
		Strict:  true,
	})
	if err == nil {
		t.Fatal("Load() error = nil, want invalid-duration failure")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error leaked secret value: %q", err)
	}
	if !strings.Contains(err.Error(), `configuration key "secret_timeout" has an invalid duration value`) {
		t.Fatalf("error = %q, want safe typed failure", err)
	}
}

func TestSecretRedactionAndProvenanceNeverExposeValue(t *testing.T) {
	const secret = "do-not-print-this-token"
	definitions := []Definition{
		{Key: "api_token", Kind: KindString, Required: true, Sensitive: true},
		{Key: "region", Kind: KindString, Default: "local"},
	}
	cfg, err := Load(definitions, Options{
		Environ: map[string]string{"OMNEXA_API_TOKEN": secret},
		Strict:  true,
	})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	redactedValues := cfg.Redacted()
	if redactedValues["api_token"] != redacted {
		t.Fatalf("api_token diagnostic = %q, want redacted", redactedValues["api_token"])
	}
	if redactedValues["region"] != "local" {
		t.Fatalf("region diagnostic = %q, want local", redactedValues["region"])
	}
	if strings.Contains(strings.Join(mapValues(redactedValues), " "), secret) {
		t.Fatal("redacted diagnostics leaked secret")
	}
	for _, item := range cfg.Provenance() {
		if strings.Contains(item.Key+string(item.Source), secret) {
			t.Fatal("provenance leaked secret")
		}
	}
}

func TestUnknownOmnexaEnvironmentKeyFailsStrictly(t *testing.T) {
	_, err := LoadApplication(Options{
		Environ: map[string]string{"OMNEXA_ENVIRONMNET": "production"},
		Strict:  true,
	})
	if err == nil {
		t.Fatal("LoadApplication() error = nil, want unknown-key failure")
	}
	if !strings.Contains(err.Error(), "unknown Omnexa environment configuration key") {
		t.Fatalf("error = %q", err)
	}
}

func TestUnknownFileKeyFailsStrictly(t *testing.T) {
	path := writeConfigFile(t, `{"environment":"ci","environmnet":"production"}`)
	_, err := LoadApplication(Options{FilePath: path, Strict: true})
	if err == nil {
		t.Fatal("LoadApplication() error = nil, want unknown file key failure")
	}
	if !strings.Contains(err.Error(), `unknown configuration key "environmnet" from file`) {
		t.Fatalf("error = %q", err)
	}
}

func TestInvalidEnvironmentFailsWithoutEchoingRawValue(t *testing.T) {
	const raw = "production-secret-like-value"
	_, err := LoadApplication(Options{
		Environ: map[string]string{"OMNEXA_ENVIRONMENT": raw},
		Strict:  true,
	})
	if err == nil {
		t.Fatal("LoadApplication() error = nil, want invalid environment failure")
	}
	if strings.Contains(err.Error(), raw) {
		t.Fatalf("error leaked raw invalid value: %q", err)
	}
}

func TestLoaderInstancesDoNotShareOverrides(t *testing.T) {
	first, err := LoadApplication(Options{Overrides: map[string]string{"environment": "test"}, Strict: true})
	if err != nil {
		t.Fatalf("first LoadApplication() error = %v", err)
	}
	second, err := LoadApplication(Options{Overrides: map[string]string{"environment": "staging"}, Strict: true})
	if err != nil {
		t.Fatalf("second LoadApplication() error = %v", err)
	}
	third, err := LoadApplication(Options{Strict: true})
	if err != nil {
		t.Fatalf("third LoadApplication() error = %v", err)
	}

	if got, _ := first.Environment("environment"); got != EnvironmentTest {
		t.Fatalf("first environment = %q", got)
	}
	if got, _ := second.Environment("environment"); got != EnvironmentStaging {
		t.Fatalf("second environment = %q", got)
	}
	if got, _ := third.Environment("environment"); got != EnvironmentLocal {
		t.Fatalf("third environment = %q, want local with no leaked override", got)
	}
}

func TestOSOptionsUsesExplicitConfigPathAndCapturedEnvironment(t *testing.T) {
	options := OSOptions([]string{
		"PATH=/usr/bin",
		"OMNEXA_CONFIG_FILE=/tmp/omnexa-safe-config.json",
		"OMNEXA_ENVIRONMENT=ci",
	})
	if options.FilePath != "/tmp/omnexa-safe-config.json" {
		t.Fatalf("FilePath = %q", options.FilePath)
	}
	if options.Environ["OMNEXA_ENVIRONMENT"] != "ci" {
		t.Fatalf("captured environment = %#v", options.Environ)
	}
	if !options.Strict {
		t.Fatal("OSOptions Strict = false, want true")
	}
}

func TestDefinitionValidationRejectsUnsafeSchema(t *testing.T) {
	tests := []struct {
		name        string
		definitions []Definition
	}{
		{name: "bad key", definitions: []Definition{{Key: "Bad-Key", Kind: KindString}}},
		{name: "bad env prefix", definitions: []Definition{{Key: "region", Env: "REGION", Kind: KindString}}},
		{name: "unsupported kind", definitions: []Definition{{Key: "region", Kind: Kind("object")}}},
		{name: "duplicate key", definitions: []Definition{{Key: "region", Kind: KindString}, {Key: "region", Kind: KindString}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := Load(test.definitions, Options{}); err == nil {
				t.Fatal("Load() error = nil, want schema validation failure")
			}
		})
	}
}

func writeConfigFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func mapValues(values map[string]string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}
