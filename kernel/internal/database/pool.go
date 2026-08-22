package database

import (
	"context"
	"errors"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	databaseURLKey                   = "database_url"
	databaseConnectTimeoutKey        = "database_connect_timeout"
	databaseMaxConnectionsKey        = "database_max_connections"
	databaseMinConnectionsKey        = "database_min_connections"
	databaseMaxConnectionLifetimeKey = "database_max_connection_lifetime"
	databaseMaxConnectionIdleTimeKey = "database_max_connection_idle_time"
)

// PoolSettings is the bounded PostgreSQL pool configuration consumed by P01.04.
// URL is RESTRICTED configuration and must never be rendered in diagnostics.
type PoolSettings struct {
	URL                   string
	ConnectTimeout        time.Duration
	MaxConnections        int32
	MinConnections        int32
	MaxConnectionLifetime time.Duration
	MaxConnectionIdleTime time.Duration
}

// LoadConfiguration composes the completed P01.02 application schema with the
// P01.04 PostgreSQL settings without making the existing kernel startup path
// database-dependent. Provider secrets remain marked sensitive in Config.
func LoadConfiguration(options config.Options) (config.Config, error) {
	definitions := append(config.ApplicationSchema(), databaseDefinitions()...)
	resolved, err := config.Load(definitions, options)
	if err != nil {
		return config.Config{}, safeWrappedFailure(
			err,
			codeConfigurationInvalid,
			failure.CategoryValidation,
			"database configuration is invalid",
		)
	}
	return resolved, nil
}

func databaseDefinitions() []config.Definition {
	return []config.Definition{
		{Key: databaseURLKey, Env: "OMNEXA_DATABASE_URL", Kind: config.KindString, Sensitive: true},
		{Key: databaseConnectTimeoutKey, Env: "OMNEXA_DATABASE_CONNECT_TIMEOUT", Kind: config.KindDuration, Default: "5s"},
		{Key: databaseMaxConnectionsKey, Env: "OMNEXA_DATABASE_MAX_CONNECTIONS", Kind: config.KindInt, Default: "10"},
		{Key: databaseMinConnectionsKey, Env: "OMNEXA_DATABASE_MIN_CONNECTIONS", Kind: config.KindInt, Default: "0"},
		{Key: databaseMaxConnectionLifetimeKey, Env: "OMNEXA_DATABASE_MAX_CONNECTION_LIFETIME", Kind: config.KindDuration, Default: "30m"},
		{Key: databaseMaxConnectionIdleTimeKey, Env: "OMNEXA_DATABASE_MAX_CONNECTION_IDLE_TIME", Kind: config.KindDuration, Default: "5m"},
	}
}

// SettingsFromConfig validates the P01.04 PostgreSQL settings without exposing
// the connection URL in returned errors.
func SettingsFromConfig(resolved config.Config) (PoolSettings, error) {
	url, ok := resolved.String(databaseURLKey)
	if !ok || url == "" {
		return PoolSettings{}, safeFailure(
			codeConfigurationInvalid,
			failure.CategoryValidation,
			"database configuration is invalid",
			failure.WithDetail("database_url is required when PostgreSQL is enabled"),
		)
	}

	connectTimeout, ok := resolved.Duration(databaseConnectTimeoutKey)
	if !ok || connectTimeout <= 0 {
		return PoolSettings{}, invalidSetting(databaseConnectTimeoutKey)
	}
	maxConnections, ok := resolved.Int(databaseMaxConnectionsKey)
	if !ok || maxConnections < 1 || maxConnections > 1000 {
		return PoolSettings{}, invalidSetting(databaseMaxConnectionsKey)
	}
	minConnections, ok := resolved.Int(databaseMinConnectionsKey)
	if !ok || minConnections < 0 || minConnections > maxConnections {
		return PoolSettings{}, invalidSetting(databaseMinConnectionsKey)
	}
	maxLifetime, ok := resolved.Duration(databaseMaxConnectionLifetimeKey)
	if !ok || maxLifetime <= 0 {
		return PoolSettings{}, invalidSetting(databaseMaxConnectionLifetimeKey)
	}
	maxIdleTime, ok := resolved.Duration(databaseMaxConnectionIdleTimeKey)
	if !ok || maxIdleTime <= 0 {
		return PoolSettings{}, invalidSetting(databaseMaxConnectionIdleTimeKey)
	}

	return PoolSettings{
		URL:                   url,
		ConnectTimeout:        connectTimeout,
		MaxConnections:        int32(maxConnections),
		MinConnections:        int32(minConnections),
		MaxConnectionLifetime: maxLifetime,
		MaxConnectionIdleTime: maxIdleTime,
	}, nil
}

func invalidSetting(key string) error {
	return safeFailure(
		codeConfigurationInvalid,
		failure.CategoryValidation,
		"database configuration is invalid",
		failure.WithDetail(key+" is outside the supported P01.04 bounds"),
	)
}

// NewPool constructs and verifies a PostgreSQL pool from P01.02 configuration.
// Connection establishment and the initial ping are bounded by ConnectTimeout.
func NewPool(ctx context.Context, resolved config.Config) (*pgxpool.Pool, error) {
	settings, err := SettingsFromConfig(resolved)
	if err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(settings.URL)
	if err != nil {
		return nil, safeWrappedFailure(
			err,
			codeConfigurationInvalid,
			failure.CategoryValidation,
			"database configuration is invalid",
		)
	}
	poolConfig.ConnConfig.ConnectTimeout = settings.ConnectTimeout
	poolConfig.MaxConns = settings.MaxConnections
	poolConfig.MinConns = settings.MinConnections
	poolConfig.MaxConnLifetime = settings.MaxConnectionLifetime
	poolConfig.MaxConnIdleTime = settings.MaxConnectionIdleTime

	bounded, cancel := context.WithTimeout(ctx, settings.ConnectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(bounded, poolConfig)
	if err != nil {
		return nil, classifyConnectionFailure(err)
	}
	if err := pool.Ping(bounded); err != nil {
		pool.Close()
		return nil, classifyConnectionFailure(err)
	}
	return pool, nil
}

func classifyConnectionFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) || errors.Is(cause, context.Canceled) {
		return safeWrappedFailure(
			cause,
			codeConnectionUnavailable,
			failure.CategoryTimeout,
			"database connection timed out",
			failure.WithRetryable(true),
		)
	}
	return safeWrappedFailure(
		cause,
		codeConnectionUnavailable,
		failure.CategoryUnavailable,
		"database connection is unavailable",
		failure.WithRetryable(true),
	)
}
