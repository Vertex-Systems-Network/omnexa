package cache

import (
	"strings"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

func TestLoadConfigurationRedactsSensitiveProviderValues(t *testing.T) {
	const (
		address  = "cache.internal.example:6379"
		username = "restricted-user"
		password = "restricted-password"
	)
	resolved, err := LoadConfiguration(config.Options{
		Environ: map[string]string{
			"OMNEXA_CACHE_ADDRESS":  address,
			"OMNEXA_CACHE_USERNAME": username,
			"OMNEXA_CACHE_PASSWORD": password,
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}

	redacted := resolved.Redacted()
	for _, key := range []string{cacheAddressKey, cacheUsernameKey, cachePasswordKey} {
		if redacted[key] != "<redacted>" {
			t.Fatalf("%s diagnostic = %q, want redacted", key, redacted[key])
		}
	}
	joined := strings.Join(mapValues(redacted), " ")
	for _, secret := range []string{address, username, password} {
		if strings.Contains(joined, secret) {
			t.Fatalf("redacted configuration leaked sensitive provider value %q", secret)
		}
	}
}

func TestSettingsFromConfigAppliesBoundedDefaults(t *testing.T) {
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{cacheAddressKey: "127.0.0.1:6379"},
		Strict:    true,
	})
	if err != nil {
		t.Fatalf("LoadConfiguration() error = %v", err)
	}
	settings, err := SettingsFromConfig(resolved)
	if err != nil {
		t.Fatalf("SettingsFromConfig() error = %v", err)
	}
	if settings.ConnectTimeout != 3*time.Second {
		t.Fatalf("ConnectTimeout = %v, want 3s", settings.ConnectTimeout)
	}
	if settings.OperationTimeout != 2*time.Second {
		t.Fatalf("OperationTimeout = %v, want 2s", settings.OperationTimeout)
	}
	if settings.KeyPrefix != "omnexa" {
		t.Fatalf("KeyPrefix = %q, want omnexa", settings.KeyPrefix)
	}
	if settings.MaxValueBytes != 1048576 {
		t.Fatalf("MaxValueBytes = %d, want 1048576", settings.MaxValueBytes)
	}
	if settings.MaxTTL != 24*time.Hour {
		t.Fatalf("MaxTTL = %v, want 24h", settings.MaxTTL)
	}
}

func TestSettingsFromConfigRejectsUnsafeBounds(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
	}{
		{name: "missing address", overrides: map[string]string{}},
		{name: "zero connect timeout", overrides: map[string]string{cacheAddressKey: "127.0.0.1:6379", cacheConnectTimeoutKey: "0s"}},
		{name: "long connect timeout", overrides: map[string]string{cacheAddressKey: "127.0.0.1:6379", cacheConnectTimeoutKey: "31s"}},
		{name: "zero operation timeout", overrides: map[string]string{cacheAddressKey: "127.0.0.1:6379", cacheOperationTimeoutKey: "0s"}},
		{name: "bad prefix", overrides: map[string]string{cacheAddressKey: "127.0.0.1:6379", cacheKeyPrefixKey: "Bad Prefix"}},
		{name: "zero max bytes", overrides: map[string]string{cacheAddressKey: "127.0.0.1:6379", cacheMaxValueBytesKey: "0"}},
		{name: "oversize max bytes", overrides: map[string]string{cacheAddressKey: "127.0.0.1:6379", cacheMaxValueBytesKey: "16777217"}},
		{name: "zero max ttl", overrides: map[string]string{cacheAddressKey: "127.0.0.1:6379", cacheMaxTTLKey: "0s"}},
		{name: "long max ttl", overrides: map[string]string{cacheAddressKey: "127.0.0.1:6379", cacheMaxTTLKey: "169h"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolved, err := LoadConfiguration(config.Options{Overrides: test.overrides, Strict: true})
			if err != nil {
				t.Fatalf("LoadConfiguration() error = %v", err)
			}
			_, err = SettingsFromConfig(resolved)
			if err == nil {
				t.Fatal("SettingsFromConfig() error = nil, want validation failure")
			}
			assertFailureCode(t, err, codeConfigurationInvalid)
		})
	}
}

func TestRenderKeyIsDeterministicVersionedAndScoped(t *testing.T) {
	key := Key{Namespace: "kernel.cache", Version: 2, Name: "fixture"}
	first, err := RenderKey("omnexa", key)
	if err != nil {
		t.Fatalf("RenderKey() error = %v", err)
	}
	second, err := RenderKey("omnexa", key)
	if err != nil {
		t.Fatalf("RenderKey() second error = %v", err)
	}
	if first != "omnexa:kernel.cache:v2:fixture" || second != first {
		t.Fatalf("rendered keys = %q / %q", first, second)
	}

	otherVersion, _ := RenderKey("omnexa", Key{Namespace: "kernel.cache", Version: 3, Name: "fixture"})
	otherNamespace, _ := RenderKey("omnexa", Key{Namespace: "kernel.other", Version: 2, Name: "fixture"})
	otherName, _ := RenderKey("omnexa", Key{Namespace: "kernel.cache", Version: 2, Name: "other"})
	for label, value := range map[string]string{
		"version":   otherVersion,
		"namespace": otherNamespace,
		"name":      otherName,
	} {
		if value == first {
			t.Fatalf("%s change collided with base cache key", label)
		}
	}
}

func TestRenderKeyRejectsInvalidSegments(t *testing.T) {
	tests := []struct {
		prefix string
		key    Key
	}{
		{prefix: "Bad Prefix", key: Key{Namespace: "kernel.cache", Version: 1, Name: "fixture"}},
		{prefix: "omnexa", key: Key{Namespace: "Kernel.Cache", Version: 1, Name: "fixture"}},
		{prefix: "omnexa", key: Key{Namespace: "kernel.cache", Version: 0, Name: "fixture"}},
		{prefix: "omnexa", key: Key{Namespace: "kernel.cache", Version: 1, Name: "bad:name"}},
	}
	for _, test := range tests {
		_, err := RenderKey(test.prefix, test.key)
		if err == nil {
			t.Fatal("RenderKey() error = nil, want validation failure")
		}
		assertFailureCode(t, err, codeKeyInvalid)
	}
}

func TestJSONCodecRoundTripAndCorruptionFailure(t *testing.T) {
	type fixture struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	codec := JSONCodec[fixture]{}
	input := fixture{Name: "synthetic", Count: 7}
	encoded, err := codec.Encode(input)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if decoded != input {
		t.Fatalf("decoded = %#v, want %#v", decoded, input)
	}

	_, err = codec.Decode([]byte(`{"name":`))
	if err == nil {
		t.Fatal("Decode(corrupt) error = nil")
	}
	assertFailureCode(t, err, codeSerializationFailed)
}

func TestJSONCodecRejectsUnsupportedValueWithoutLeakingIt(t *testing.T) {
	type invalid struct {
		Channel chan int `json:"channel"`
	}
	codec := JSONCodec[invalid]{}
	_, err := codec.Encode(invalid{Channel: make(chan int)})
	if err == nil {
		t.Fatal("Encode() error = nil, want unsupported-value failure")
	}
	assertFailureCode(t, err, codeSerializationFailed)
}

func assertFailureCode(t *testing.T, err error, code failure.Code) {
	t.Helper()
	if !failure.IsCode(err, code) {
		actual, _ := failure.CodeOf(err)
		t.Fatalf("failure code = %q, want %q (error %v)", actual, code, err)
	}
}

func mapValues(values map[string]string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}
