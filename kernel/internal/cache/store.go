package cache

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	valkey "github.com/valkey-io/valkey-go"
)

const (
	cacheAddressKey          = "cache_address"
	cacheUsernameKey         = "cache_username"
	cachePasswordKey         = "cache_password"
	cacheConnectTimeoutKey   = "cache_connect_timeout"
	cacheOperationTimeoutKey = "cache_operation_timeout"
	cacheKeyPrefixKey        = "cache_key_prefix"
	cacheMaxValueBytesKey    = "cache_max_value_bytes"
	cacheMaxTTLKey           = "cache_max_ttl"

	defaultBlockingPoolSize = 16
)

// Settings is the bounded P01.05 provider configuration. Address, Username and
// Password are sensitive infrastructure configuration and must not be rendered
// in public diagnostics or artifacts.
type Settings struct {
	Address          string
	Username         string
	Password         string
	ConnectTimeout   time.Duration
	OperationTimeout time.Duration
	KeyPrefix        string
	MaxValueBytes    int
	MaxTTL           time.Duration
}

// Entry distinguishes a cache miss from a provider failure. Value is copied so
// callers do not retain provider-owned storage.
type Entry struct {
	Value []byte
	Found bool
}

// Store is the bounded Redis-compatible P01.05 provider adapter. It is a cache,
// never an authoritative persistence or session/authorization primitive.
type Store struct {
	client   valkey.Client
	settings Settings
}

// LoadConfiguration composes P01.05 settings over the completed P01.02 schema
// without making the current kernel startup path cache-dependent.
func LoadConfiguration(options config.Options) (config.Config, error) {
	definitions := append(config.ApplicationSchema(), cacheDefinitions()...)
	resolved, err := config.Load(definitions, options)
	if err != nil {
		return config.Config{}, safeWrappedFailure(
			err,
			codeConfigurationInvalid,
			failure.CategoryValidation,
			"cache configuration is invalid",
		)
	}
	return resolved, nil
}

func cacheDefinitions() []config.Definition {
	return []config.Definition{
		{Key: cacheAddressKey, Env: "OMNEXA_CACHE_ADDRESS", Kind: config.KindString, Sensitive: true},
		{Key: cacheUsernameKey, Env: "OMNEXA_CACHE_USERNAME", Kind: config.KindString, Sensitive: true},
		{Key: cachePasswordKey, Env: "OMNEXA_CACHE_PASSWORD", Kind: config.KindString, Sensitive: true},
		{Key: cacheConnectTimeoutKey, Env: "OMNEXA_CACHE_CONNECT_TIMEOUT", Kind: config.KindDuration, Default: "3s"},
		{Key: cacheOperationTimeoutKey, Env: "OMNEXA_CACHE_OPERATION_TIMEOUT", Kind: config.KindDuration, Default: "2s"},
		{Key: cacheKeyPrefixKey, Env: "OMNEXA_CACHE_KEY_PREFIX", Kind: config.KindString, Default: "omnexa"},
		{Key: cacheMaxValueBytesKey, Env: "OMNEXA_CACHE_MAX_VALUE_BYTES", Kind: config.KindInt, Default: "1048576"},
		{Key: cacheMaxTTLKey, Env: "OMNEXA_CACHE_MAX_TTL", Kind: config.KindDuration, Default: "24h"},
	}
}

// SettingsFromConfig validates bounded P01.05 settings without exposing
// provider endpoint or credentials.
func SettingsFromConfig(resolved config.Config) (Settings, error) {
	address, ok := resolved.String(cacheAddressKey)
	if !ok || address == "" {
		return Settings{}, invalidSetting(cacheAddressKey)
	}
	username, _ := resolved.String(cacheUsernameKey)
	password, _ := resolved.String(cachePasswordKey)
	connectTimeout, ok := resolved.Duration(cacheConnectTimeoutKey)
	if !ok || connectTimeout <= 0 || connectTimeout > 30*time.Second {
		return Settings{}, invalidSetting(cacheConnectTimeoutKey)
	}
	operationTimeout, ok := resolved.Duration(cacheOperationTimeoutKey)
	if !ok || operationTimeout <= 0 || operationTimeout > 30*time.Second {
		return Settings{}, invalidSetting(cacheOperationTimeoutKey)
	}
	prefix, ok := resolved.String(cacheKeyPrefixKey)
	if !ok || !keySegmentPattern.MatchString(prefix) {
		return Settings{}, invalidSetting(cacheKeyPrefixKey)
	}
	maxValueBytes, ok := resolved.Int(cacheMaxValueBytesKey)
	if !ok || maxValueBytes < 1 || maxValueBytes > 16*1024*1024 {
		return Settings{}, invalidSetting(cacheMaxValueBytesKey)
	}
	maxTTL, ok := resolved.Duration(cacheMaxTTLKey)
	if !ok || maxTTL < time.Millisecond || maxTTL > 7*24*time.Hour {
		return Settings{}, invalidSetting(cacheMaxTTLKey)
	}

	return Settings{
		Address:          address,
		Username:         username,
		Password:         password,
		ConnectTimeout:   connectTimeout,
		OperationTimeout: operationTimeout,
		KeyPrefix:        prefix,
		MaxValueBytes:    maxValueBytes,
		MaxTTL:           maxTTL,
	}, nil
}

func invalidSetting(key string) error {
	return safeFailure(
		codeConfigurationInvalid,
		failure.CategoryValidation,
		"cache configuration is invalid",
		failure.WithDetail(key+" is missing or outside the supported P01.05 bounds"),
	)
}

// NewStore constructs and verifies one standalone Redis-compatible provider
// endpoint. Cluster/sentinel/HA topology is not pulled forward by P01.05.
func NewStore(ctx context.Context, resolved config.Config) (*Store, error) {
	settings, err := SettingsFromConfig(resolved)
	if err != nil {
		return nil, err
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress:           []string{settings.Address},
		Username:              settings.Username,
		Password:              settings.Password,
		Dialer:                net.Dialer{Timeout: settings.ConnectTimeout, KeepAlive: time.Second},
		ConnWriteTimeout:      settings.OperationTimeout,
		DisableRetry:          true,
		DisableCache:          true,
		DisableAutoPipelining: true,
		ForceSingleClient:     true,
		BlockingPoolSize:      defaultBlockingPoolSize,
		BlockingPoolCleanup:   time.Minute,
		BlockingPoolMinSize:   0,
		ClientSetInfo:         valkey.DisableClientSetInfo,
	})
	if err != nil {
		return nil, classifyConnectionFailure(err)
	}

	bounded, cancel := context.WithTimeout(ctx, settings.ConnectTimeout)
	defer cancel()
	if err := client.Do(bounded, client.B().Ping().Build()).Error(); err != nil {
		client.Close()
		return nil, classifyConnectionFailure(err)
	}

	return &Store{client: client, settings: settings}, nil
}

// Close releases provider resources. A closed Store must not be reused.
func (store *Store) Close() {
	if store == nil || store.client == nil {
		return
	}
	store.client.Close()
}

// Get returns Found=false for a provider nil/miss. Provider errors remain errors
// and are never collapsed into a miss.
func (store *Store) Get(ctx context.Context, key Key) (Entry, error) {
	rendered, err := store.renderKey(key)
	if err != nil {
		return Entry{}, err
	}
	bounded, cancel := context.WithTimeout(ctx, store.settings.OperationTimeout)
	defer cancel()

	value, err := store.client.Do(bounded, store.client.B().Get().Key(rendered).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return Entry{Found: false}, nil
		}
		return Entry{}, classifyOperationFailure(err)
	}
	return Entry{Value: []byte(value), Found: true}, nil
}

// Set stores a bounded value with an explicit positive TTL. P01.05 deliberately
// does not support unbounded cache retention.
func (store *Store) Set(ctx context.Context, key Key, value []byte, ttl time.Duration) error {
	rendered, err := store.renderKey(key)
	if err != nil {
		return err
	}
	if len(value) > store.settings.MaxValueBytes {
		return safeFailure(
			codeValueInvalid,
			failure.CategoryValidation,
			"cache value is invalid",
			failure.WithDetail("cache value exceeds the configured maximum size"),
		)
	}
	if ttl < time.Millisecond || ttl > store.settings.MaxTTL {
		return safeFailure(
			codeValueInvalid,
			failure.CategoryValidation,
			"cache value is invalid",
			failure.WithDetail("cache TTL is outside the configured bounds"),
		)
	}

	bounded, cancel := context.WithTimeout(ctx, store.settings.OperationTimeout)
	defer cancel()
	if err := store.client.Do(
		bounded,
		store.client.B().Set().Key(rendered).Value(string(value)).Px(ttl).Build(),
	).Error(); err != nil {
		return classifyOperationFailure(err)
	}
	return nil
}

// Delete removes one cache entry and reports whether a value existed.
func (store *Store) Delete(ctx context.Context, key Key) (bool, error) {
	rendered, err := store.renderKey(key)
	if err != nil {
		return false, err
	}
	bounded, cancel := context.WithTimeout(ctx, store.settings.OperationTimeout)
	defer cancel()
	count, err := store.client.Do(bounded, store.client.B().Del().Key(rendered).Build()).ToInt64()
	if err != nil {
		return false, classifyOperationFailure(err)
	}
	return count > 0, nil
}

func (store *Store) renderKey(key Key) (string, error) {
	if store == nil || store.client == nil {
		return "", safeFailure(
			codeOperationFailed,
			failure.CategoryInvariant,
			"cache store is not initialized",
		)
	}
	return RenderKey(store.settings.KeyPrefix, key)
}

func classifyConnectionFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) || errors.Is(cause, context.Canceled) {
		return safeWrappedFailure(
			cause,
			codeConnectionUnavailable,
			failure.CategoryTimeout,
			"cache connection timed out",
			failure.WithRetryable(true),
		)
	}
	return safeWrappedFailure(
		cause,
		codeConnectionUnavailable,
		failure.CategoryUnavailable,
		"cache connection is unavailable",
		failure.WithRetryable(true),
	)
}

func classifyOperationFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) || errors.Is(cause, context.Canceled) {
		return safeWrappedFailure(
			cause,
			codeOperationFailed,
			failure.CategoryTimeout,
			"cache operation timed out",
			failure.WithRetryable(true),
		)
	}
	return safeWrappedFailure(
		cause,
		codeOperationFailed,
		failure.CategoryUnavailable,
		"cache operation failed",
		failure.WithRetryable(true),
	)
}
