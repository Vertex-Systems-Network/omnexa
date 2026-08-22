package storage

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const (
	storageEndpointKey         = "storage_endpoint"
	storageAccessKeyKey        = "storage_access_key"
	storageSecretKeyKey        = "storage_secret_key"
	storageRegionKey           = "storage_region"
	storageBucketKey           = "storage_bucket"
	storageUsePathStyleKey     = "storage_use_path_style"
	storageConnectTimeoutKey   = "storage_connect_timeout"
	storageOperationTimeoutKey = "storage_operation_timeout"
	storageKeyPrefixKey        = "storage_key_prefix"
	storageMaxObjectBytesKey   = "storage_max_object_bytes"

	maxSimpleObjectBytes = int64(5 * 1024 * 1024 * 1024)
)

var (
	bucketPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)
	regionPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}$`)
)

// Settings is the bounded P01.06 S3-compatible provider configuration.
// Endpoint and credentials are sensitive infrastructure configuration and must
// never be rendered in public diagnostics or artifacts.
type Settings struct {
	Endpoint         string
	AccessKey        string
	SecretKey        string
	Region           string
	Bucket           string
	UsePathStyle     bool
	ConnectTimeout   time.Duration
	OperationTimeout time.Duration
	KeyPrefix        string
	MaxObjectBytes   int64
}

// Store is the bounded S3-compatible P01.06 provider adapter. It owns no CMS,
// media-library, tenant-document or business-domain semantics.
type Store struct {
	client    *s3.Client
	transport *http.Transport
	settings  Settings
}

// LoadConfiguration composes P01.06 settings over the completed P01.02 schema
// without making the current kernel startup path storage-dependent.
func LoadConfiguration(options config.Options) (config.Config, error) {
	definitions := append(config.ApplicationSchema(), storageDefinitions()...)
	resolved, err := config.Load(definitions, options)
	if err != nil {
		return config.Config{}, safeWrappedFailure(
			err,
			codeConfigurationInvalid,
			failure.CategoryValidation,
			"storage configuration is invalid",
		)
	}
	return resolved, nil
}

func storageDefinitions() []config.Definition {
	return []config.Definition{
		{Key: storageEndpointKey, Env: "OMNEXA_STORAGE_ENDPOINT", Kind: config.KindString, Sensitive: true},
		{Key: storageAccessKeyKey, Env: "OMNEXA_STORAGE_ACCESS_KEY", Kind: config.KindString, Sensitive: true},
		{Key: storageSecretKeyKey, Env: "OMNEXA_STORAGE_SECRET_KEY", Kind: config.KindString, Sensitive: true},
		{Key: storageRegionKey, Env: "OMNEXA_STORAGE_REGION", Kind: config.KindString, Default: "us-east-1"},
		{Key: storageBucketKey, Env: "OMNEXA_STORAGE_BUCKET", Kind: config.KindString},
		{Key: storageUsePathStyleKey, Env: "OMNEXA_STORAGE_USE_PATH_STYLE", Kind: config.KindBool, Default: "true"},
		{Key: storageConnectTimeoutKey, Env: "OMNEXA_STORAGE_CONNECT_TIMEOUT", Kind: config.KindDuration, Default: "3s"},
		{Key: storageOperationTimeoutKey, Env: "OMNEXA_STORAGE_OPERATION_TIMEOUT", Kind: config.KindDuration, Default: "5s"},
		{Key: storageKeyPrefixKey, Env: "OMNEXA_STORAGE_KEY_PREFIX", Kind: config.KindString, Default: "omnexa"},
		{Key: storageMaxObjectBytesKey, Env: "OMNEXA_STORAGE_MAX_OBJECT_BYTES", Kind: config.KindInt, Default: "1073741824"},
	}
}

// SettingsFromConfig validates bounded P01.06 settings without publishing the
// provider endpoint or credentials.
func SettingsFromConfig(resolved config.Config) (Settings, error) {
	endpoint, ok := resolved.String(storageEndpointKey)
	if !ok || !validEndpoint(endpoint) {
		return Settings{}, invalidSetting(storageEndpointKey)
	}
	accessKey, ok := resolved.String(storageAccessKeyKey)
	if !ok || accessKey == "" {
		return Settings{}, invalidSetting(storageAccessKeyKey)
	}
	secretKey, ok := resolved.String(storageSecretKeyKey)
	if !ok || secretKey == "" {
		return Settings{}, invalidSetting(storageSecretKeyKey)
	}
	region, ok := resolved.String(storageRegionKey)
	if !ok || !regionPattern.MatchString(region) {
		return Settings{}, invalidSetting(storageRegionKey)
	}
	bucket, ok := resolved.String(storageBucketKey)
	if !ok || !validBucket(bucket) {
		return Settings{}, invalidSetting(storageBucketKey)
	}
	usePathStyle, ok := resolved.Bool(storageUsePathStyleKey)
	if !ok {
		return Settings{}, invalidSetting(storageUsePathStyleKey)
	}
	connectTimeout, ok := resolved.Duration(storageConnectTimeoutKey)
	if !ok || connectTimeout <= 0 || connectTimeout > 30*time.Second {
		return Settings{}, invalidSetting(storageConnectTimeoutKey)
	}
	operationTimeout, ok := resolved.Duration(storageOperationTimeoutKey)
	if !ok || operationTimeout <= 0 || operationTimeout > 2*time.Minute {
		return Settings{}, invalidSetting(storageOperationTimeoutKey)
	}
	prefix, ok := resolved.String(storageKeyPrefixKey)
	if !ok || !objectNamespacePattern.MatchString(prefix) {
		return Settings{}, invalidSetting(storageKeyPrefixKey)
	}
	maxObjectBytes, ok := resolved.Int(storageMaxObjectBytesKey)
	if !ok || maxObjectBytes < 1 || int64(maxObjectBytes) > maxSimpleObjectBytes {
		return Settings{}, invalidSetting(storageMaxObjectBytesKey)
	}

	return Settings{
		Endpoint:         endpoint,
		AccessKey:        accessKey,
		SecretKey:        secretKey,
		Region:           region,
		Bucket:           bucket,
		UsePathStyle:     usePathStyle,
		ConnectTimeout:   connectTimeout,
		OperationTimeout: operationTimeout,
		KeyPrefix:        prefix,
		MaxObjectBytes:   int64(maxObjectBytes),
	}, nil
}

func invalidSetting(key string) error {
	return safeFailure(
		codeConfigurationInvalid,
		failure.CategoryValidation,
		"storage configuration is invalid",
		failure.WithDetail(key+" is missing or outside the supported P01.06 bounds"),
	)
}

func validEndpoint(value string) bool {
	parsed, err := url.Parse(value)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return false
	}
	return parsed.User == nil && parsed.RawQuery == "" && parsed.Fragment == "" && (parsed.Path == "" || parsed.Path == "/")
}

func validBucket(value string) bool {
	return bucketPattern.MatchString(value) && !strings.Contains(value, "..")
}

// NewStore constructs an S3-compatible provider and verifies access to the
// configured bucket with a bounded HeadBucket request.
func NewStore(ctx context.Context, resolved config.Config) (*Store, error) {
	settings, err := SettingsFromConfig(resolved)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   settings.ConnectTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          32,
		MaxIdleConnsPerHost:   16,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   settings.ConnectTimeout,
		ResponseHeaderTimeout: settings.OperationTimeout,
	}
	httpClient := &http.Client{Transport: transport}
	awsConfig := aws.Config{
		Region: settings.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			settings.AccessKey,
			settings.SecretKey,
			"",
		)),
		HTTPClient:   httpClient,
		BaseEndpoint: aws.String(settings.Endpoint),
		Retryer: func() aws.Retryer {
			return aws.NopRetryer{}
		},
	}
	client := s3.NewFromConfig(awsConfig, func(options *s3.Options) {
		options.UsePathStyle = settings.UsePathStyle
	})
	store := &Store{client: client, transport: transport, settings: settings}

	bounded, cancel := context.WithTimeout(ctx, settings.ConnectTimeout)
	defer cancel()
	if _, err := client.HeadBucket(bounded, &s3.HeadBucketInput{Bucket: aws.String(settings.Bucket)}); err != nil {
		transport.CloseIdleConnections()
		return nil, classifyConnectionFailure(err)
	}
	return store, nil
}

// Close releases idle HTTP provider connections. A closed Store must not be
// reused by callers.
func (store *Store) Close() {
	if store == nil || store.transport == nil {
		return
	}
	store.transport.CloseIdleConnections()
}

// Put streams one bounded object to the provider. Caller-supplied SHA-256 is
// sent as both provider checksum input and immutable Omnexa integrity metadata.
func (store *Store) Put(ctx context.Context, key Key, upload Upload) (ObjectInfo, error) {
	rendered, err := store.renderKey(key)
	if err != nil {
		return ObjectInfo{}, err
	}
	metadata, err := validateUpload(upload, store.settings.MaxObjectBytes)
	if err != nil {
		return ObjectInfo{}, err
	}
	checksumBytes, _ := hex.DecodeString(metadata[checksumMetadataKey])
	checksumBase64 := base64.StdEncoding.EncodeToString(checksumBytes)

	bounded, cancel := context.WithTimeout(ctx, store.settings.OperationTimeout)
	defer cancel()
	output, err := store.client.PutObject(bounded, &s3.PutObjectInput{
		Bucket:         aws.String(store.settings.Bucket),
		Key:            aws.String(rendered),
		Body:           newIntegrityReader(upload.Body, upload.ContentLength, metadata[checksumMetadataKey]),
		ContentLength:  aws.Int64(upload.ContentLength),
		ContentType:    optionalString(upload.ContentType),
		ChecksumSHA256: aws.String(checksumBase64),
		Metadata:       metadata,
	})
	if err != nil {
		return ObjectInfo{}, classifyOperationFailure(err)
	}

	return ObjectInfo{
		Key:           rendered,
		ContentLength: upload.ContentLength,
		ContentType:   upload.ContentType,
		FileName:      upload.FileName,
		SHA256:        metadata[checksumMetadataKey],
		ETag:          stringValue(output.ETag),
		Metadata:      publicMetadata(metadata),
	}, nil
}

// Head returns provider-neutral metadata without opening the object body.
func (store *Store) Head(ctx context.Context, key Key) (ObjectInfo, error) {
	rendered, err := store.renderKey(key)
	if err != nil {
		return ObjectInfo{}, err
	}
	bounded, cancel := context.WithTimeout(ctx, store.settings.OperationTimeout)
	defer cancel()
	output, err := store.client.HeadObject(bounded, &s3.HeadObjectInput{
		Bucket: aws.String(store.settings.Bucket),
		Key:    aws.String(rendered),
	})
	if err != nil {
		return ObjectInfo{}, classifyObjectOperationFailure(err)
	}
	return objectInfo(
		rendered,
		output.ContentLength,
		output.ContentType,
		output.ETag,
		output.LastModified,
		output.Metadata,
	)
}

// Open returns a streaming download whose terminal read verifies the stored
// SHA-256 metadata. Provider errors and missing objects remain distinct.
func (store *Store) Open(ctx context.Context, key Key) (*ObjectReader, error) {
	rendered, err := store.renderKey(key)
	if err != nil {
		return nil, err
	}
	bounded, cancel := context.WithTimeout(ctx, store.settings.OperationTimeout)
	defer cancel()
	output, err := store.client.GetObject(bounded, &s3.GetObjectInput{
		Bucket: aws.String(store.settings.Bucket),
		Key:    aws.String(rendered),
	})
	if err != nil {
		return nil, classifyObjectOperationFailure(err)
	}
	info, err := objectInfo(
		rendered,
		output.ContentLength,
		output.ContentType,
		output.ETag,
		output.LastModified,
		output.Metadata,
	)
	if err != nil {
		_ = output.Body.Close()
		return nil, err
	}
	return &ObjectReader{
		Info: info,
		Body: newVerifiedReadCloser(output.Body, info.ContentLength, info.SHA256),
	}, nil
}

// Delete removes an object. S3-compatible delete is intentionally idempotent;
// callers use Head/Open when they need explicit missing-object evidence.
func (store *Store) Delete(ctx context.Context, key Key) error {
	rendered, err := store.renderKey(key)
	if err != nil {
		return err
	}
	bounded, cancel := context.WithTimeout(ctx, store.settings.OperationTimeout)
	defer cancel()
	_, err = store.client.DeleteObject(bounded, &s3.DeleteObjectInput{
		Bucket: aws.String(store.settings.Bucket),
		Key:    aws.String(rendered),
	})
	if err != nil {
		return classifyOperationFailure(err)
	}
	return nil
}

func (store *Store) renderKey(key Key) (string, error) {
	if store == nil || store.client == nil {
		return "", safeFailure(
			codeOperationFailed,
			failure.CategoryInvariant,
			"storage provider is not initialized",
		)
	}
	return RenderKey(store.settings.KeyPrefix, key)
}

func objectInfo(
	key string,
	contentLength *int64,
	contentType *string,
	etag *string,
	lastModified *time.Time,
	metadata map[string]string,
) (ObjectInfo, error) {
	if contentLength == nil || *contentLength < 0 {
		return ObjectInfo{}, safeFailure(
			codeIntegrityFailed,
			failure.CategoryInvariant,
			"storage object integrity metadata is invalid",
		)
	}
	checksum, ok := metadata[checksumMetadataKey]
	if !ok {
		return ObjectInfo{}, safeFailure(
			codeIntegrityFailed,
			failure.CategoryInvariant,
			"storage object integrity metadata is missing",
		)
	}
	checksum, err := normalizeChecksum(checksum)
	if err != nil {
		return ObjectInfo{}, safeFailure(
			codeIntegrityFailed,
			failure.CategoryInvariant,
			"storage object integrity metadata is invalid",
		)
	}
	fileName := metadata[fileNameMetadataKey]
	if err := validateFileName(fileName); err != nil {
		return ObjectInfo{}, safeFailure(
			codeIntegrityFailed,
			failure.CategoryInvariant,
			"storage object metadata is invalid",
		)
	}
	return ObjectInfo{
		Key:           key,
		ContentLength: *contentLength,
		ContentType:   stringValue(contentType),
		FileName:      fileName,
		SHA256:        checksum,
		ETag:          stringValue(etag),
		LastModified:  timeValue(lastModified),
		Metadata:      publicMetadata(metadata),
	}, nil
}

func publicMetadata(metadata map[string]string) map[string]string {
	result := make(map[string]string, len(metadata))
	for key, value := range metadata {
		if key == checksumMetadataKey || key == fileNameMetadataKey {
			continue
		}
		result[key] = value
	}
	return result
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return aws.String(value)
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func timeValue(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return *value
}

func classifyConnectionFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) || errors.Is(cause, context.Canceled) {
		return safeWrappedFailure(
			cause,
			codeConnectionUnavailable,
			failure.CategoryTimeout,
			"storage connection timed out",
			failure.WithRetryable(true),
		)
	}
	return safeWrappedFailure(
		cause,
		codeConnectionUnavailable,
		failure.CategoryUnavailable,
		"storage connection is unavailable",
		failure.WithRetryable(true),
	)
}

func classifyObjectOperationFailure(cause error) error {
	var responseError *smithyhttp.ResponseError
	if errors.As(cause, &responseError) && responseError.HTTPStatusCode() == http.StatusNotFound {
		return safeWrappedFailure(
			cause,
			codeObjectNotFound,
			failure.CategoryNotFound,
			"storage object was not found",
		)
	}
	return classifyOperationFailure(cause)
}

func classifyOperationFailure(cause error) error {
	if errors.Is(cause, context.DeadlineExceeded) || errors.Is(cause, context.Canceled) {
		return safeWrappedFailure(
			cause,
			codeOperationFailed,
			failure.CategoryTimeout,
			"storage operation timed out",
			failure.WithRetryable(true),
		)
	}
	return safeWrappedFailure(
		cause,
		codeOperationFailed,
		failure.CategoryUnavailable,
		"storage operation failed",
		failure.WithRetryable(true),
	)
}

// drainVerified consumes an ObjectReader with bounded working memory. It is used
// by deterministic integration evidence and deliberately never returns content.
func drainVerified(reader *ObjectReader) error {
	if reader == nil || reader.Body == nil {
		return safeFailure(codeOperationFailed, failure.CategoryInvariant, "storage object reader is not initialized")
	}
	_, err := io.Copy(io.Discard, reader.Body)
	return err
}
