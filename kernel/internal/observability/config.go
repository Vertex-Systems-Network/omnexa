package observability

import (
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/buildinfo"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	observabilityEnabledKey          = "observability_enabled"
	observabilityServiceNameKey      = "observability_service_name"
	observabilityLogLevelKey         = "observability_log_level"
	observabilityExportIntervalKey   = "observability_export_interval"
	observabilityExportTimeoutKey    = "observability_export_timeout"
	observabilityShutdownTimeoutKey  = "observability_shutdown_timeout"
	defaultObservabilityServiceName  = "omnexa-kernel"
	defaultObservabilityExportPeriod = 30 * time.Second
)

var serviceNamePattern = regexp.MustCompile(`^[a-z][a-z0-9._-]{1,62}$`)

// Settings is the bounded P01.07 observability configuration. It contains no
// exporter credentials or backend-specific configuration.
type Settings struct {
	Enabled         bool
	ServiceName     string
	ServiceVersion  string
	Environment     config.Environment
	LogLevel        slog.Level
	ExportInterval  time.Duration
	ExportTimeout   time.Duration
	ShutdownTimeout time.Duration
}

// LoadConfiguration composes P01.07 settings over the completed P01.02 schema.
func LoadConfiguration(options config.Options) (config.Config, error) {
	definitions := append(config.ApplicationSchema(), observabilityDefinitions()...)
	resolved, err := config.Load(definitions, options)
	if err != nil {
		return config.Config{}, safeWrappedFailure(
			err,
			codeConfigurationInvalid,
			failure.CategoryValidation,
			"observability configuration is invalid",
		)
	}
	return resolved, nil
}

func observabilityDefinitions() []config.Definition {
	return []config.Definition{
		{Key: observabilityEnabledKey, Env: "OMNEXA_OBSERVABILITY_ENABLED", Kind: config.KindBool, Default: "true"},
		{Key: observabilityServiceNameKey, Env: "OMNEXA_OBSERVABILITY_SERVICE_NAME", Kind: config.KindString, Default: defaultObservabilityServiceName},
		{Key: observabilityLogLevelKey, Env: "OMNEXA_OBSERVABILITY_LOG_LEVEL", Kind: config.KindString, Default: "auto"},
		{Key: observabilityExportIntervalKey, Env: "OMNEXA_OBSERVABILITY_EXPORT_INTERVAL", Kind: config.KindDuration, Default: defaultObservabilityExportPeriod.String()},
		{Key: observabilityExportTimeoutKey, Env: "OMNEXA_OBSERVABILITY_EXPORT_TIMEOUT", Kind: config.KindDuration, Default: "3s"},
		{Key: observabilityShutdownTimeoutKey, Env: "OMNEXA_OBSERVABILITY_SHUTDOWN_TIMEOUT", Kind: config.KindDuration, Default: "5s"},
	}
}

// SettingsFromConfig validates the process-wide P01.07 configuration.
func SettingsFromConfig(resolved config.Config) (Settings, error) {
	environment, ok := resolved.Environment("environment")
	if !ok || !environment.Valid() {
		return Settings{}, invalidSetting("environment")
	}
	enabled, ok := resolved.Bool(observabilityEnabledKey)
	if !ok {
		return Settings{}, invalidSetting(observabilityEnabledKey)
	}
	serviceName, ok := resolved.String(observabilityServiceNameKey)
	if !ok || !serviceNamePattern.MatchString(serviceName) {
		return Settings{}, invalidSetting(observabilityServiceNameKey)
	}
	levelText, ok := resolved.String(observabilityLogLevelKey)
	if !ok {
		return Settings{}, invalidSetting(observabilityLogLevelKey)
	}
	level, ok := resolveLogLevel(levelText, environment)
	if !ok {
		return Settings{}, invalidSetting(observabilityLogLevelKey)
	}
	exportInterval, ok := resolved.Duration(observabilityExportIntervalKey)
	if !ok || exportInterval < time.Second || exportInterval > 10*time.Minute {
		return Settings{}, invalidSetting(observabilityExportIntervalKey)
	}
	exportTimeout, ok := resolved.Duration(observabilityExportTimeoutKey)
	if !ok || exportTimeout < 10*time.Millisecond || exportTimeout > 30*time.Second {
		return Settings{}, invalidSetting(observabilityExportTimeoutKey)
	}
	shutdownTimeout, ok := resolved.Duration(observabilityShutdownTimeoutKey)
	if !ok || shutdownTimeout < 10*time.Millisecond || shutdownTimeout > 30*time.Second {
		return Settings{}, invalidSetting(observabilityShutdownTimeoutKey)
	}

	return Settings{
		Enabled:         enabled,
		ServiceName:     serviceName,
		ServiceVersion:  buildinfo.Current().Version,
		Environment:     environment,
		LogLevel:        level,
		ExportInterval:  exportInterval,
		ExportTimeout:   exportTimeout,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}

func invalidSetting(key string) error {
	return safeFailure(
		codeConfigurationInvalid,
		failure.CategoryValidation,
		"observability configuration is invalid",
		failure.WithDetail(key+" is missing or outside the supported P01.07 bounds"),
	)
}

func resolveLogLevel(raw string, environment config.Environment) (slog.Level, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "auto":
		if environment == config.EnvironmentLocal || environment == config.EnvironmentTest {
			return slog.LevelDebug, true
		}
		return slog.LevelInfo, true
	case "debug":
		return slog.LevelDebug, true
	case "info":
		return slog.LevelInfo, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return 0, false
	}
}
