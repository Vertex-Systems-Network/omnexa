package observability

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"go.opentelemetry.io/otel/trace"
)

const (
	redactedValue      = "<redacted>"
	redactedErrorValue = "<redacted-error>"
)

// Classification is the P01.07 logging-side view of the frozen data
// classification vocabulary. It does not create new data-authority semantics.
type Classification string

const (
	ClassificationPublic       Classification = "PUBLIC"
	ClassificationInternal     Classification = "INTERNAL"
	ClassificationConfidential Classification = "CONFIDENTIAL"
	ClassificationRestricted   Classification = "RESTRICTED"
)

// Classified marks an attribute with its existing data-classification level.
// CONFIDENTIAL/RESTRICTED values are redacted by default.
type Classified struct {
	Classification Classification
	Value          any
}

// Sensitive explicitly marks a value that must never be emitted.
type Sensitive struct{ Value any }

// SafeError returns a log-safe representation of an error. Structured Omnexa
// failures expose only stable code/category; arbitrary error text is redacted.
func SafeError(err error) any {
	if err == nil {
		return nil
	}
	if classified, ok := failure.As(err); ok {
		return slog.GroupValue(
			slog.String("code", string(classified.Code())),
			slog.String("category", string(classified.Category())),
		)
	}
	return redactedErrorValue
}

// Logger is the P01.07 structured logging boundary. Log messages themselves
// must be static/safe text; untrusted or sensitive values belong in attributes
// where classification/redaction rules can be enforced.
type Logger struct {
	logger *slog.Logger
}

// NewJSONLogger returns a structured JSON logger with stable field names and
// fail-closed attribute redaction.
func NewJSONLogger(writer io.Writer, settings Settings) *Logger {
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: settings.LogLevel,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if len(groups) != 0 {
				return attr
			}
			switch attr.Key {
			case slog.TimeKey:
				attr.Key = "timestamp"
			case slog.LevelKey:
				attr.Key = "severity"
			case slog.MessageKey:
				attr.Key = "message"
			}
			return attr
		},
	})
	return newLogger(safeHandler{next: handler}, settings)
}

func newLogger(handler slog.Handler, settings Settings) *Logger {
	base := slog.New(handler).With(
		slog.String("service", settings.ServiceName),
		slog.String("environment", string(settings.Environment)),
	)
	return &Logger{logger: base}
}

// Debug emits a debug record with governed context fields.
func (logger *Logger) Debug(ctx context.Context, message string, attrs ...slog.Attr) {
	logger.log(ctx, slog.LevelDebug, message, attrs...)
}

// Info emits an info record with governed context fields.
func (logger *Logger) Info(ctx context.Context, message string, attrs ...slog.Attr) {
	logger.log(ctx, slog.LevelInfo, message, attrs...)
}

// Warn emits a warning record with governed context fields.
func (logger *Logger) Warn(ctx context.Context, message string, attrs ...slog.Attr) {
	logger.log(ctx, slog.LevelWarn, message, attrs...)
}

// Error emits an error record with governed context fields.
func (logger *Logger) Error(ctx context.Context, message string, attrs ...slog.Attr) {
	logger.log(ctx, slog.LevelError, message, attrs...)
}

func (logger *Logger) log(ctx context.Context, level slog.Level, message string, attrs ...slog.Attr) {
	if logger == nil || logger.logger == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	contextAttrs := traceFields(ctx)
	all := make([]slog.Attr, 0, len(contextAttrs)+len(attrs))
	all = append(all, contextAttrs...)
	all = append(all, attrs...)
	logger.logger.LogAttrs(ctx, level, message, all...)
}

func traceFields(ctx context.Context) []slog.Attr {
	attrs := make([]slog.Attr, 0, 3)
	if value, ok := CorrelationIDFromContext(ctx); ok {
		attrs = append(attrs, slog.String("correlation_id", value))
	}
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		attrs = append(attrs,
			slog.String("trace_id", spanContext.TraceID().String()),
			slog.String("span_id", spanContext.SpanID().String()),
		)
	}
	return attrs
}

type safeHandler struct{ next slog.Handler }

func (handler safeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return handler.next.Enabled(ctx, level)
}

func (handler safeHandler) Handle(ctx context.Context, record slog.Record) error {
	clean := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	record.Attrs(func(attr slog.Attr) bool {
		clean.AddAttrs(sanitizeAttr(attr))
		return true
	})
	return handler.next.Handle(ctx, clean)
}

func (handler safeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clean := make([]slog.Attr, 0, len(attrs))
	for _, attr := range attrs {
		clean = append(clean, sanitizeAttr(attr))
	}
	return safeHandler{next: handler.next.WithAttrs(clean)}
}

func (handler safeHandler) WithGroup(name string) slog.Handler {
	return safeHandler{next: handler.next.WithGroup(name)}
}

func sanitizeAttr(attr slog.Attr) slog.Attr {
	if attr.Key == "" {
		return attr
	}
	if sensitiveKey(attr.Key) {
		return slog.String(attr.Key, redactedValue)
	}
	value := attr.Value.Resolve()
	if value.Kind() == slog.KindGroup {
		members := value.Group()
		clean := make([]slog.Attr, 0, len(members))
		for _, member := range members {
			clean = append(clean, sanitizeAttr(member))
		}
		return slog.Attr{Key: attr.Key, Value: slog.GroupValue(clean...)}
	}
	if value.Kind() == slog.KindAny {
		switch typed := value.Any().(type) {
		case Sensitive:
			return slog.String(attr.Key, redactedValue)
		case *Sensitive:
			return slog.String(attr.Key, redactedValue)
		case Classified:
			return classifiedAttr(attr.Key, typed)
		case *Classified:
			if typed == nil {
				return slog.Any(attr.Key, nil)
			}
			return classifiedAttr(attr.Key, *typed)
		case error:
			return slog.Any(attr.Key, SafeError(typed))
		}
	}
	return slog.Attr{Key: attr.Key, Value: value}
}

func classifiedAttr(key string, value Classified) slog.Attr {
	switch value.Classification {
	case ClassificationPublic, ClassificationInternal:
		return sanitizeAttr(slog.Any(key, value.Value))
	case ClassificationConfidential, ClassificationRestricted:
		return slog.String(key, redactedValue)
	default:
		return slog.String(key, redactedValue)
	}
}

func sensitiveKey(key string) bool {
	key = strings.ToLower(strings.TrimSpace(key))
	if key == "" {
		return true
	}
	sensitive := []string{
		"password", "passwd", "secret", "token", "authorization", "cookie",
		"api_key", "apikey", "access_key", "private_key", "credential",
	}
	for _, marker := range sensitive {
		if key == marker || strings.HasSuffix(key, "_"+marker) || strings.HasPrefix(key, marker+"_") {
			return true
		}
	}
	return false
}

// Capture stores deterministic structured records for tests without serializing
// them through an external telemetry backend.
type Capture struct {
	mu      sync.Mutex
	records []CapturedRecord
}

// CapturedRecord is a defensive test-only projection of one structured record.
type CapturedRecord struct {
	Timestamp   time.Time
	Severity    slog.Level
	Message     string
	Service     string
	Environment string
	Attributes  map[string]any
}

// NewCaptureLogger returns a logger that uses the same redaction boundary as the
// production JSON logger and stores records in memory for deterministic tests.
func NewCaptureLogger(settings Settings) (*Logger, *Capture) {
	capture := &Capture{}
	handler := safeHandler{next: &captureHandler{capture: capture}}
	return newLogger(handler, settings), capture
}

// Records returns a defensive copy of captured records.
func (capture *Capture) Records() []CapturedRecord {
	if capture == nil {
		return nil
	}
	capture.mu.Lock()
	defer capture.mu.Unlock()
	result := make([]CapturedRecord, len(capture.records))
	for index, record := range capture.records {
		copyRecord := record
		copyRecord.Attributes = cloneMap(record.Attributes)
		result[index] = copyRecord
	}
	return result
}

type captureHandler struct {
	capture *Capture
	attrs   []slog.Attr
	groups  []string
}

func (handler *captureHandler) Enabled(context.Context, slog.Level) bool { return true }

func (handler *captureHandler) Handle(_ context.Context, record slog.Record) error {
	attributes := make(map[string]any)
	for _, attr := range handler.attrs {
		addCaptured(attributes, handler.groups, attr)
	}
	record.Attrs(func(attr slog.Attr) bool {
		addCaptured(attributes, handler.groups, attr)
		return true
	})

	captured := CapturedRecord{
		Timestamp:  record.Time,
		Severity:   record.Level,
		Message:    record.Message,
		Attributes: attributes,
	}
	if service, ok := attributes["service"].(string); ok {
		captured.Service = service
		delete(attributes, "service")
	}
	if environment, ok := attributes["environment"].(string); ok {
		captured.Environment = environment
		delete(attributes, "environment")
	}
	handler.capture.mu.Lock()
	handler.capture.records = append(handler.capture.records, captured)
	handler.capture.mu.Unlock()
	return nil
}

func (handler *captureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := &captureHandler{capture: handler.capture, groups: append([]string(nil), handler.groups...)}
	clone.attrs = append(append([]slog.Attr(nil), handler.attrs...), attrs...)
	return clone
}

func (handler *captureHandler) WithGroup(name string) slog.Handler {
	clone := &captureHandler{capture: handler.capture, attrs: append([]slog.Attr(nil), handler.attrs...)}
	clone.groups = append(append([]string(nil), handler.groups...), name)
	return clone
}

func addCaptured(target map[string]any, groups []string, attr slog.Attr) {
	current := target
	for _, group := range groups {
		next, ok := current[group].(map[string]any)
		if !ok {
			next = make(map[string]any)
			current[group] = next
		}
		current = next
	}
	current[attr.Key] = capturedValue(attr.Value.Resolve())
}

func capturedValue(value slog.Value) any {
	switch value.Kind() {
	case slog.KindGroup:
		group := make(map[string]any)
		for _, attr := range value.Group() {
			group[attr.Key] = capturedValue(attr.Value.Resolve())
		}
		return group
	case slog.KindString:
		return value.String()
	case slog.KindBool:
		return value.Bool()
	case slog.KindInt64:
		return value.Int64()
	case slog.KindUint64:
		return value.Uint64()
	case slog.KindFloat64:
		return value.Float64()
	case slog.KindDuration:
		return value.Duration()
	case slog.KindTime:
		return value.Time()
	case slog.KindAny:
		return value.Any()
	default:
		return value.Any()
	}
}

func cloneMap(input map[string]any) map[string]any {
	result := make(map[string]any, len(input))
	for key, value := range input {
		if nested, ok := value.(map[string]any); ok {
			result[key] = cloneMap(nested)
			continue
		}
		result[key] = value
	}
	return result
}

var _ slog.Handler = safeHandler{}
var _ slog.Handler = (*captureHandler)(nil)
