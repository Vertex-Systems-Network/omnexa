package observability

import (
	"context"
	"unicode"
	"unicode/utf8"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"go.opentelemetry.io/otel/propagation"
)

const (
	correlationHeader   = "x-correlation-id"
	maxCorrelationRunes = 128
)

type correlationContextKey struct{}

// WithCorrelationID attaches a bounded diagnostic identifier. Correlation IDs
// are never treated as authorization, tenancy, or business-state authority.
func WithCorrelationID(ctx context.Context, value string) (context.Context, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !validCorrelationID(value) {
		return ctx, safeFailure(
			codeCorrelationInvalid,
			failure.CategoryValidation,
			"observability correlation identifier is invalid",
		)
	}
	return context.WithValue(ctx, correlationContextKey{}, value), nil
}

// CorrelationIDFromContext returns the governed diagnostic correlation ID.
func CorrelationIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	value, ok := ctx.Value(correlationContextKey{}).(string)
	return value, ok && value != ""
}

func validCorrelationID(value string) bool {
	if value == "" || utf8.RuneCountInString(value) > maxCorrelationRunes {
		return false
	}
	for _, char := range value {
		if unicode.IsSpace(char) || unicode.IsControl(char) {
			return false
		}
	}
	return true
}

// Propagator returns the P01.07 vendor-neutral W3C trace-context plus bounded
// Omnexa correlation-ID propagator. Baggage is intentionally excluded so
// arbitrary caller values cannot silently become propagated telemetry data.
func Propagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		correlationPropagator{},
	)
}

// Inject writes only governed trace/correlation context into the carrier.
func Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	Propagator().Inject(ctx, carrier)
}

// Extract reads governed trace/correlation context from the carrier.
func Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return Propagator().Extract(ctx, carrier)
}

type correlationPropagator struct{}

func (correlationPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	if value, ok := CorrelationIDFromContext(ctx); ok {
		carrier.Set(correlationHeader, value)
	}
}

func (correlationPropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	value := carrier.Get(correlationHeader)
	if !validCorrelationID(value) {
		return ctx
	}
	return context.WithValue(ctx, correlationContextKey{}, value)
}

func (correlationPropagator) Fields() []string {
	return []string{correlationHeader}
}
