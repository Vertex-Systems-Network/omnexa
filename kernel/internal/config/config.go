// Package config provides the P01.02 typed, deterministic process-configuration foundation.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Environment is a governed process environment identity. It is configuration,
// not authorization, tenancy, or entitlement state.
type Environment string

const (
	EnvironmentLocal      Environment = "local"
	EnvironmentCI         Environment = "ci"
	EnvironmentPreview    Environment = "preview"
	EnvironmentTest       Environment = "test"
	EnvironmentStaging    Environment = "staging"
	EnvironmentProduction Environment = "production"
)

// Kind defines the typed representation of a configuration value.
type Kind string

const (
	KindString      Kind = "string"
	KindBool        Kind = "bool"
	KindInt         Kind = "int"
	KindDuration    Kind = "duration"
	KindEnvironment Kind = "environment"
)

// Source records the winning source after deterministic precedence resolution.
type Source string

const (
	SourceDefault     Source = "default"
	SourceFile        Source = "file"
	SourceEnvironment Source = "environment"
	SourceOverride    Source = "override"
)

const redacted = "<redacted>"

var keyPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Definition declares the type, requirement, default, and sensitivity for one
// stable lowercase_snake_case configuration key.
type Definition struct {
	Key       string
	Env       string
	Kind      Kind
	Required  bool
	Default   string
	Sensitive bool
}

// Options supplies explicit configuration sources. FilePath is optional and is
// never auto-discovered. Environ should be a captured environment map; nil means
// no environment variables. Overrides are intended for narrowly scoped tests or
// explicit caller-controlled construction, never global mutable runtime state.
type Options struct {
	FilePath  string
	Environ   map[string]string
	Overrides map[string]string
	Strict    bool
}

// Provenance is safe diagnostic metadata. It intentionally contains no value.
type Provenance struct {
	Key       string
	Source    Source
	Sensitive bool
}

type resolvedValue struct {
	value     any
	source    Source
	sensitive bool
}

// Config is an immutable resolved configuration snapshot.
type Config struct {
	values map[string]resolvedValue
	order  []string
}

// ApplicationSchema is the bounded P01.02 process schema. Later packages add
// their own keys only through governed changes; this package does not pre-create
// database, cache, telemetry, tenancy, or business configuration.
func ApplicationSchema() []Definition {
	return []Definition{
		{
			Key:     "environment",
			Env:     "OMNEXA_ENVIRONMENT",
			Kind:    KindEnvironment,
			Default: string(EnvironmentLocal),
		},
	}
}

// LoadApplication resolves the current application schema from explicit options.
func LoadApplication(options Options) (Config, error) {
	return Load(ApplicationSchema(), options)
}

// OSOptions snapshots the current process environment for the application
// boundary. It performs no mutation and no secret logging.
func OSOptions(environ []string) Options {
	values := EnvironMap(environ)
	return Options{
		FilePath: values["OMNEXA_CONFIG_FILE"],
		Environ:  values,
		Strict:   true,
	}
}

// EnvironMap converts an os.Environ-style list into an isolated map.
func EnvironMap(environ []string) map[string]string {
	values := make(map[string]string, len(environ))
	for _, item := range environ {
		key, value, ok := strings.Cut(item, "=")
		if !ok || key == "" {
			continue
		}
		values[key] = value
	}
	return values
}

// Load resolves definitions with deterministic precedence:
// default -> explicit JSON file -> environment -> explicit override.
func Load(definitions []Definition, options Options) (Config, error) {
	if len(definitions) == 0 {
		return Config{}, errors.New("configuration schema is empty")
	}

	defs := make(map[string]Definition, len(definitions))
	order := make([]string, 0, len(definitions))
	envToKey := make(map[string]string, len(definitions))
	for _, definition := range definitions {
		normalized, err := normalizeDefinition(definition)
		if err != nil {
			return Config{}, err
		}
		if _, exists := defs[normalized.Key]; exists {
			return Config{}, fmt.Errorf("duplicate configuration key %q", normalized.Key)
		}
		if existing, exists := envToKey[normalized.Env]; exists {
			return Config{}, fmt.Errorf("configuration environment variable %q is shared by keys %q and %q", normalized.Env, existing, normalized.Key)
		}
		defs[normalized.Key] = normalized
		envToKey[normalized.Env] = normalized.Key
		order = append(order, normalized.Key)
	}

	raw := make(map[string]rawValue, len(definitions))
	for _, key := range order {
		definition := defs[key]
		if definition.Default != "" {
			raw[key] = rawValue{text: definition.Default, source: SourceDefault}
		}
	}

	if options.FilePath != "" {
		fileValues, err := readJSONFile(options.FilePath)
		if err != nil {
			return Config{}, err
		}
		if err := applyNamedValues(raw, fileValues, defs, SourceFile, options.Strict); err != nil {
			return Config{}, err
		}
	}

	if options.Strict {
		for envKey := range options.Environ {
			if !strings.HasPrefix(envKey, "OMNEXA_") || envKey == "OMNEXA_CONFIG_FILE" {
				continue
			}
			if _, known := envToKey[envKey]; !known {
				return Config{}, fmt.Errorf("unknown Omnexa environment configuration key %q", envKey)
			}
		}
	}
	for envKey, key := range envToKey {
		if value, exists := options.Environ[envKey]; exists {
			raw[key] = rawValue{text: value, source: SourceEnvironment}
		}
	}

	if err := applyNamedValues(raw, options.Overrides, defs, SourceOverride, true); err != nil {
		return Config{}, err
	}

	resolved := make(map[string]resolvedValue, len(definitions))
	for _, key := range order {
		definition := defs[key]
		candidate, exists := raw[key]
		if !exists || strings.TrimSpace(candidate.text) == "" {
			if definition.Required {
				return Config{}, fmt.Errorf("required configuration key %q is missing", key)
			}
			continue
		}
		value, err := parseValue(definition, candidate.text)
		if err != nil {
			return Config{}, err
		}
		resolved[key] = resolvedValue{value: value, source: candidate.source, sensitive: definition.Sensitive}
	}

	return Config{values: resolved, order: order}, nil
}

type rawValue struct {
	text   string
	source Source
}

func normalizeDefinition(definition Definition) (Definition, error) {
	if !keyPattern.MatchString(definition.Key) {
		return Definition{}, fmt.Errorf("configuration key %q must be lowercase snake_case", definition.Key)
	}
	if definition.Env == "" {
		definition.Env = "OMNEXA_" + strings.ToUpper(definition.Key)
	}
	if !strings.HasPrefix(definition.Env, "OMNEXA_") {
		return Definition{}, fmt.Errorf("environment variable for configuration key %q must use OMNEXA_ prefix", definition.Key)
	}
	switch definition.Kind {
	case KindString, KindBool, KindInt, KindDuration, KindEnvironment:
	default:
		return Definition{}, fmt.Errorf("configuration key %q has unsupported type %q", definition.Key, definition.Kind)
	}
	return definition, nil
}

func applyNamedValues(target map[string]rawValue, values map[string]string, defs map[string]Definition, source Source, strict bool) error {
	for key, value := range values {
		if _, known := defs[key]; !known {
			if strict {
				return fmt.Errorf("unknown configuration key %q from %s", key, source)
			}
			continue
		}
		target[key] = rawValue{text: value, source: source}
	}
	return nil
}

func readJSONFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read configuration file: %w", err)
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, fmt.Errorf("parse configuration file: %w", err)
	}
	values := make(map[string]string, len(object))
	for key, raw := range object {
		var text string
		if err := json.Unmarshal(raw, &text); err == nil {
			values[key] = text
			continue
		}
		var boolean bool
		if err := json.Unmarshal(raw, &boolean); err == nil {
			values[key] = strconv.FormatBool(boolean)
			continue
		}
		var number json.Number
		decoder := json.NewDecoder(strings.NewReader(string(raw)))
		decoder.UseNumber()
		if err := decoder.Decode(&number); err == nil {
			values[key] = number.String()
			continue
		}
		return nil, fmt.Errorf("configuration file key %q must be a string, boolean, or number", key)
	}
	return values, nil
}

func parseValue(definition Definition, raw string) (any, error) {
	switch definition.Kind {
	case KindString:
		return raw, nil
	case KindBool:
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, invalidValueError(definition)
		}
		return value, nil
	case KindInt:
		value, err := strconv.Atoi(raw)
		if err != nil {
			return nil, invalidValueError(definition)
		}
		return value, nil
	case KindDuration:
		value, err := time.ParseDuration(raw)
		if err != nil {
			return nil, invalidValueError(definition)
		}
		return value, nil
	case KindEnvironment:
		value := Environment(raw)
		if !value.Valid() {
			return nil, invalidValueError(definition)
		}
		return value, nil
	default:
		return nil, fmt.Errorf("configuration key %q has unsupported type", definition.Key)
	}
}

func invalidValueError(definition Definition) error {
	return fmt.Errorf("configuration key %q has an invalid %s value", definition.Key, definition.Kind)
}

// Valid reports whether the environment identity is governed by the P01.02 contract.
func (environment Environment) Valid() bool {
	switch environment {
	case EnvironmentLocal, EnvironmentCI, EnvironmentPreview, EnvironmentTest, EnvironmentStaging, EnvironmentProduction:
		return true
	default:
		return false
	}
}

// Environment returns a typed environment value.
func (config Config) Environment(key string) (Environment, bool) {
	value, ok := config.values[key]
	if !ok {
		return "", false
	}
	typed, ok := value.value.(Environment)
	return typed, ok
}

// String returns a typed string value.
func (config Config) String(key string) (string, bool) {
	value, ok := config.values[key]
	if !ok {
		return "", false
	}
	typed, ok := value.value.(string)
	return typed, ok
}

// Bool returns a typed boolean value.
func (config Config) Bool(key string) (bool, bool) {
	value, ok := config.values[key]
	if !ok {
		return false, false
	}
	typed, ok := value.value.(bool)
	return typed, ok
}

// Int returns a typed integer value.
func (config Config) Int(key string) (int, bool) {
	value, ok := config.values[key]
	if !ok {
		return 0, false
	}
	typed, ok := value.value.(int)
	return typed, ok
}

// Duration returns a typed duration value.
func (config Config) Duration(key string) (time.Duration, bool) {
	value, ok := config.values[key]
	if !ok {
		return 0, false
	}
	typed, ok := value.value.(time.Duration)
	return typed, ok
}

// Provenance returns safe source metadata in schema order without values.
func (config Config) Provenance() []Provenance {
	result := make([]Provenance, 0, len(config.order))
	for _, key := range config.order {
		value, exists := config.values[key]
		if !exists {
			continue
		}
		result = append(result, Provenance{Key: key, Source: value.source, Sensitive: value.sensitive})
	}
	return result
}

// Redacted returns diagnostic display values. Sensitive values are never exposed.
func (config Config) Redacted() map[string]string {
	result := make(map[string]string, len(config.values))
	for key, value := range config.values {
		if value.sensitive {
			result[key] = redacted
			continue
		}
		result[key] = fmt.Sprint(value.value)
	}
	return result
}
