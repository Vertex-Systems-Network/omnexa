package developer

import (
	"reflect"
	"testing"
)

func TestVerificationEnvironmentDropsHostCredentials(t *testing.T) {
	environment := []string{
		"PATH=/usr/bin",
		"HOME=/home/runner",
		"RUNNER_TEMP=/tmp/runner",
		"P04_04_TEST_DATABASE_URL=synthetic-database",
		"OPENAI_API_KEY=synthetic-openai-secret",
		"ANTHROPIC_API_KEY=synthetic-anthropic-secret",
		"GITHUB_TOKEN=synthetic-github-secret",
		"AWS_ACCESS_KEY_ID=synthetic-aws-access",
		"AWS_SECRET_ACCESS_KEY=synthetic-aws-secret",
		"SSH_AUTH_SOCK=/tmp/agent.sock",
		"DATABASE_URL=postgres://production.example.invalid/db",
		"OMNEXA_DATABASE_URL=postgres://restricted.example.invalid/db",
		"P04_04_SECRET=must-not-pass",
	}

	got := verificationEnvironment(environment)
	want := []string{
		"PATH=/usr/bin",
		"HOME=/home/runner",
		"RUNNER_TEMP=/tmp/runner",
		"P04_04_TEST_DATABASE_URL=synthetic-database",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("verificationEnvironment() = %q, want %q", got, want)
	}
}

func TestVerificationEnvironmentAllowsOnlyStructuredSyntheticFixtures(t *testing.T) {
	cases := map[string]bool{
		"P01_04_TEST_DATABASE_URL": true,
		"P99_99_TEST_PROVIDER_URL": true,
		"P01_TEST_DATABASE_URL":    false,
		"P01_XX_TEST_DATABASE_URL": false,
		"P01_04_SECRET":            false,
		"OPENAI_API_KEY":           false,
		"ANTHROPIC_API_KEY":        false,
		"GITHUB_TOKEN":             false,
	}
	for key, want := range cases {
		if got := allowedVerificationEnvironmentKey(key); got != want {
			t.Errorf("allowedVerificationEnvironmentKey(%q) = %t, want %t", key, got, want)
		}
	}
}
