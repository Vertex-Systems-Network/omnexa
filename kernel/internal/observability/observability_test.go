package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func testSettings() Settings {
	return Settings{
		Enabled:         true,
		ServiceName:     "omnexa-kernel",
		ServiceVersion:  "test",
		Environment:     config.EnvironmentTest,
		LogLevel:        slog.LevelDebug,
		ExportInterval:  time.Second,
		ExportTimeout:   100 * time.Millisecond,
		ShutdownTimeout: 100 * time.Millisecond,
	}
}

func TestSettingsDefaultsAndEnvironmentLogLevels(t *testing.T) {
	resolved, err := LoadConfiguration(config.Options{Overrides: map[string]string{"environment": "local"}, Strict: true})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	settings, err := SettingsFromConfig(resolved)
	if err != nil {
		t.Fatalf("settings: %v", err)
	}
	if settings.ServiceName != "omnexa-kernel" || !settings.Enabled {
		t.Fatalf("unexpected defaults: %+v", settings)
	}
	if settings.LogLevel != slog.LevelDebug {
		t.Fatalf("local auto log level = %v, want debug", settings.LogLevel)
	}

	resolved, err = LoadConfiguration(config.Options{Overrides: map[string]string{"environment": "production"}, Strict: true})
	if err != nil {
		t.Fatalf("load production config: %v", err)
	}
	settings, err = SettingsFromConfig(resolved)
	if err != nil {
		t.Fatalf("production settings: %v", err)
	}
	if settings.LogLevel != slog.LevelInfo {
		t.Fatalf("production auto log level = %v, want info", settings.LogLevel)
	}
}

func TestSettingsRejectInvalidBoundsWithoutLeakingValues(t *testing.T) {
	secret := "should-not-appear"
	resolved, err := LoadConfiguration(config.Options{
		Overrides: map[string]string{
			"environment":                   "ci",
			observabilityServiceNameKey:     "INVALID NAME " + secret,
			observabilityExportTimeoutKey:   "3s",
			observabilityShutdownTimeoutKey: "5s",
			observabilityExportIntervalKey:  "30s",
			observabilityEnabledKey:         "true",
			observabilityLogLevelKey:        "info",
		},
		Strict: true,
	})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	_, err = SettingsFromConfig(resolved)
	if err == nil {
		t.Fatal("expected invalid service-name error")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error leaked configuration value: %v", err)
	}
	classified, ok := failure.As(err)
	if !ok || classified.Code() != codeConfigurationInvalid {
		t.Fatalf("error classification = %v, %v", classified, ok)
	}
}

func TestStrictConfigurationRejectsUnknownObservabilityKey(t *testing.T) {
	_, err := LoadConfiguration(config.Options{
		Environ: map[string]string{
			"OMNEXA_ENVIRONMENT":            "ci",
			"OMNEXA_OBSERVABILITY_NOT_REAL": "x",
		},
		Strict: true,
	})
	if err == nil {
		t.Fatal("expected strict unknown-key error")
	}
	if !failure.IsCode(err, codeConfigurationInvalid) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCaptureLoggerStableFieldsContextAndRedaction(t *testing.T) {
	settings := testSettings()
	logger, capture := NewCaptureLogger(settings)

	ctx, err := WithCorrelationID(context.Background(), "request-123")
	if err != nil {
		t.Fatalf("correlation: %v", err)
	}
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
		Remote:  true,
	}))

	structured, makeErr := failure.New("test.safe", failure.CategoryValidation, "safe title")
	if makeErr != nil {
		t.Fatalf("failure.New: %v", makeErr)
	}
	logger.Info(ctx, "operation completed",
		slog.String("component", "kernel"),
		slog.String("password", "super-secret"),
		slog.Any("customer_record", Classified{Classification: ClassificationConfidential, Value: "private-customer"}),
		slog.Any("internal_label", Classified{Classification: ClassificationInternal, Value: "safe-internal"}),
		slog.Any("provider_error", errors.New("provider-secret-value")),
		slog.Any("safe_error", structured),
		slog.Group("nested", slog.String("access_token", "nested-secret"), slog.String("ok", "visible")),
	)

	records := capture.Records()
	if len(records) != 1 {
		t.Fatalf("records = %d, want 1", len(records))
	}
	record := records[0]
	if record.Service != settings.ServiceName || record.Environment != string(settings.Environment) {
		t.Fatalf("identity fields = %q/%q", record.Service, record.Environment)
	}
	if record.Message != "operation completed" || record.Severity != slog.LevelInfo || record.Timestamp.IsZero() {
		t.Fatalf("unstable core fields: %+v", record)
	}
	for key, want := range map[string]any{
		"correlation_id":  "request-123",
		"trace_id":        traceID.String(),
		"span_id":         spanID.String(),
		"component":       "kernel",
		"password":        redactedValue,
		"customer_record": redactedValue,
		"internal_label":  "safe-internal",
		"provider_error":  redactedErrorValue,
	} {
		if got := record.Attributes[key]; got != want {
			t.Fatalf("attr %s = %#v, want %#v", key, got, want)
		}
	}
	nested, ok := record.Attributes["nested"].(map[string]any)
	if !ok || nested["access_token"] != redactedValue || nested["ok"] != "visible" {
		t.Fatalf("nested attrs = %#v", record.Attributes["nested"])
	}
	safe, ok := record.Attributes["safe_error"].(map[string]any)
	if !ok || safe["code"] != "test.safe" || safe["category"] != "validation" {
		t.Fatalf("safe error projection = %#v", record.Attributes["safe_error"])
	}

	encoded := stringifyCaptured(records)
	for _, secret := range []string{"super-secret", "private-customer", "provider-secret-value", "nested-secret"} {
		if strings.Contains(encoded, secret) {
			t.Fatalf("captured records leaked %q", secret)
		}
	}
}

func TestJSONLoggerUsesStableTopLevelFieldNames(t *testing.T) {
	var output bytes.Buffer
	logger := NewJSONLogger(&output, testSettings())
	logger.Warn(context.Background(), "bounded warning", slog.String("component", "kernel"))

	var record map[string]any
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatalf("decode JSON log: %v", err)
	}
	for _, key := range []string{"timestamp", "severity", "message", "service", "environment"} {
		if _, ok := record[key]; !ok {
			t.Fatalf("missing stable log field %q in %#v", key, record)
		}
	}
	for _, old := range []string{"time", "level", "msg"} {
		if _, ok := record[old]; ok {
			t.Fatalf("legacy slog field %q should be renamed", old)
		}
	}
}

func TestCorrelationAndTraceContextPropagation(t *testing.T) {
	ctx, err := WithCorrelationID(context.Background(), "corr-456")
	if err != nil {
		t.Fatalf("correlation: %v", err)
	}
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	ctx = trace.ContextWithRemoteSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		Remote:     true,
		TraceFlags: trace.FlagsSampled,
	}))

	carrier := propagation.MapCarrier{}
	Inject(ctx, carrier)
	if carrier.Get(correlationHeader) != "corr-456" || carrier.Get("traceparent") == "" {
		t.Fatalf("carrier = %#v", carrier)
	}
	if _, exists := carrier["baggage"]; exists {
		t.Fatal("P01.07 must not propagate arbitrary OTel baggage")
	}

	extracted := Extract(context.Background(), carrier)
	if value, ok := CorrelationIDFromContext(extracted); !ok || value != "corr-456" {
		t.Fatalf("extracted correlation = %q/%v", value, ok)
	}
	gotSpan := trace.SpanContextFromContext(extracted)
	if gotSpan.TraceID() != traceID || gotSpan.SpanID() != spanID {
		t.Fatalf("extracted span context = %s/%s", gotSpan.TraceID(), gotSpan.SpanID())
	}
}

func TestInvalidCorrelationIsRejectedAndNotExtracted(t *testing.T) {
	if _, err := WithCorrelationID(context.Background(), "contains space"); err == nil || !failure.IsCode(err, codeCorrelationInvalid) {
		t.Fatalf("invalid correlation error = %v", err)
	}
	carrier := propagation.MapCarrier{correlationHeader: "bad\nvalue"}
	ctx := Extract(context.Background(), carrier)
	if _, ok := CorrelationIDFromContext(ctx); ok {
		t.Fatal("invalid propagated correlation ID should be ignored")
	}
}

func TestProviderResourceIdentityAndTraceMetricCapture(t *testing.T) {
	settings := testSettings()
	traceExporter := &capturingTraceExporter{}
	metricExporter := &capturingMetricExporter{}
	provider, err := NewProvider(settings, Backends{Trace: traceExporter, Metric: metricExporter})
	if err != nil {
		t.Fatalf("new provider: %v", err)
	}
	defer func() { _ = provider.Shutdown(context.Background()) }()

	attrs := provider.Resource().Attributes()
	values := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		values[string(attr.Key)] = attr.Value.AsString()
	}
	for key, want := range map[string]string{
		"service.name":                settings.ServiceName,
		"service.version":             settings.ServiceVersion,
		"deployment.environment.name": string(settings.Environment),
		"omnexa.phase":                "P01",
	} {
		if values[key] != want {
			t.Fatalf("resource %s = %q, want %q", key, values[key], want)
		}
	}

	ctx, span := provider.Tracer("test.scope").Start(context.Background(), "test-span")
	span.SetAttributes(attribute.String("component", "kernel"))
	span.End()
	counter, err := provider.Meter("test.scope").Int64Counter("omnexa_test_counter")
	if err != nil {
		t.Fatalf("counter: %v", err)
	}
	counter.Add(ctx, 1)
	if err := provider.ForceFlush(context.Background()); err != nil {
		t.Fatalf("force flush: %v", err)
	}
	if traceExporter.Count() == 0 {
		t.Fatal("expected trace export")
	}
	if metricExporter.ExportCount() == 0 {
		t.Fatal("expected metric export")
	}
}

func TestDisabledProviderChangesNoApplicationSemantics(t *testing.T) {
	settings := testSettings()
	settings.Enabled = false
	traceExporter := &capturingTraceExporter{}
	metricExporter := &capturingMetricExporter{}
	provider, err := NewProvider(settings, Backends{Trace: traceExporter, Metric: metricExporter})
	if err != nil {
		t.Fatalf("new provider: %v", err)
	}
	_, span := provider.Tracer("disabled").Start(context.Background(), "ignored")
	span.End()
	counter, err := provider.Meter("disabled").Int64Counter("ignored_counter")
	if err != nil {
		t.Fatalf("noop counter: %v", err)
	}
	counter.Add(context.Background(), 1)
	if err := provider.ForceFlush(context.Background()); err != nil {
		t.Fatalf("disabled flush: %v", err)
	}
	if err := provider.Shutdown(context.Background()); err != nil {
		t.Fatalf("disabled shutdown: %v", err)
	}
	if traceExporter.Count() != 0 || metricExporter.ExportCount() != 0 {
		t.Fatal("disabled provider exported telemetry")
	}
}

func TestExporterFailureIsSafeAndDoesNotLeakBackendDetails(t *testing.T) {
	settings := testSettings()
	secret := "backend-secret-host"
	exporter := &capturingTraceExporter{exportErr: errors.New(secret)}
	provider, err := NewProvider(settings, Backends{Trace: exporter})
	if err != nil {
		t.Fatalf("new provider: %v", err)
	}
	_, span := provider.Tracer("test").Start(context.Background(), "span")
	span.End()
	err = provider.ForceFlush(context.Background())
	if err == nil {
		t.Fatal("expected exporter failure")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("public error leaked exporter details: %v", err)
	}
	classified, ok := failure.As(err)
	if !ok || classified.Code() != codeLifecycleFailed || classified.Category() != failure.CategoryDependency {
		t.Fatalf("classification = %#v/%v", classified, ok)
	}
}

func TestShutdownIsBoundedAndIdempotent(t *testing.T) {
	settings := testSettings()
	settings.ShutdownTimeout = 25 * time.Millisecond
	exporter := &capturingTraceExporter{blockShutdown: true}
	provider, err := NewProvider(settings, Backends{Trace: exporter})
	if err != nil {
		t.Fatalf("new provider: %v", err)
	}
	started := time.Now()
	err = provider.Shutdown(context.Background())
	elapsed := time.Since(started)
	if err == nil || !failure.IsCode(err, codeLifecycleTimeout) {
		t.Fatalf("shutdown error = %v", err)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("shutdown was not bounded: %s", elapsed)
	}
	second := provider.Shutdown(context.Background())
	if second == nil || second.Error() != err.Error() {
		t.Fatalf("idempotent shutdown error = %v, want %v", second, err)
	}
	if exporter.ShutdownCount() != 1 {
		t.Fatalf("exporter shutdown count = %d, want 1", exporter.ShutdownCount())
	}
}

func stringifyCaptured(records []CapturedRecord) string {
	data, _ := json.Marshal(records)
	return string(data)
}

type capturingTraceExporter struct {
	mu            sync.Mutex
	spans         int
	exportErr     error
	blockShutdown bool
	shutdowns     int
}

func (exporter *capturingTraceExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	exporter.mu.Lock()
	exporter.spans += len(spans)
	err := exporter.exportErr
	exporter.mu.Unlock()
	if err != nil {
		return err
	}
	return ctx.Err()
}

func (exporter *capturingTraceExporter) Shutdown(ctx context.Context) error {
	exporter.mu.Lock()
	exporter.shutdowns++
	block := exporter.blockShutdown
	exporter.mu.Unlock()
	if block {
		<-ctx.Done()
		return ctx.Err()
	}
	return nil
}

func (exporter *capturingTraceExporter) Count() int {
	exporter.mu.Lock()
	defer exporter.mu.Unlock()
	return exporter.spans
}

func (exporter *capturingTraceExporter) ShutdownCount() int {
	exporter.mu.Lock()
	defer exporter.mu.Unlock()
	return exporter.shutdowns
}

type capturingMetricExporter struct {
	mu      sync.Mutex
	exports int
}

func (*capturingMetricExporter) Temporality(kind sdkmetric.InstrumentKind) metricdata.Temporality {
	return sdkmetric.CumulativeTemporalitySelector(kind)
}

func (*capturingMetricExporter) Aggregation(kind sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return sdkmetric.DefaultAggregationSelector(kind)
}

func (exporter *capturingMetricExporter) Export(ctx context.Context, _ *metricdata.ResourceMetrics) error {
	exporter.mu.Lock()
	exporter.exports++
	exporter.mu.Unlock()
	return ctx.Err()
}

func (*capturingMetricExporter) ForceFlush(ctx context.Context) error { return ctx.Err() }
func (*capturingMetricExporter) Shutdown(ctx context.Context) error   { return ctx.Err() }

func (exporter *capturingMetricExporter) ExportCount() int {
	exporter.mu.Lock()
	defer exporter.mu.Unlock()
	return exporter.exports
}

var _ sdktrace.SpanExporter = (*capturingTraceExporter)(nil)
var _ sdkmetric.Exporter = (*capturingMetricExporter)(nil)
