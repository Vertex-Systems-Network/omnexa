package observability

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

// Backends inject vendor-neutral OpenTelemetry SDK exporter contracts. P01.07
// intentionally does not select or configure a proprietary telemetry backend.
type Backends struct {
	Trace  sdktrace.SpanExporter
	Metric sdkmetric.Exporter
}

// Provider owns isolated trace/metric SDK lifecycle. It never mutates the
// OpenTelemetry global providers or propagator.
type Provider struct {
	settings Settings
	resource *resource.Resource
	trace    *sdktrace.TracerProvider
	metric   *sdkmetric.MeterProvider

	shutdownOnce sync.Once
	shutdownErr  error
}

// NewProvider constructs the P01.07 OpenTelemetry-compatible provider baseline.
func NewProvider(settings Settings, backends Backends) (*Provider, error) {
	res := resource.NewSchemaless(
		attribute.String("service.name", settings.ServiceName),
		attribute.String("service.version", settings.ServiceVersion),
		attribute.String("deployment.environment.name", string(settings.Environment)),
		attribute.String("omnexa.phase", "P01"),
	)

	provider := &Provider{settings: settings, resource: res}
	if !settings.Enabled {
		return provider, nil
	}

	traceOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	}
	if backends.Trace != nil {
		traceOptions = append(traceOptions, sdktrace.WithBatcher(
			backends.Trace,
			sdktrace.WithBatchTimeout(settings.ExportInterval),
			sdktrace.WithExportTimeout(settings.ExportTimeout),
			sdktrace.WithMaxQueueSize(512),
			sdktrace.WithMaxExportBatchSize(128),
		))
	}
	provider.trace = sdktrace.NewTracerProvider(traceOptions...)

	metricOptions := []sdkmetric.Option{sdkmetric.WithResource(res)}
	if backends.Metric != nil {
		reader := sdkmetric.NewPeriodicReader(
			backends.Metric,
			sdkmetric.WithInterval(settings.ExportInterval),
			sdkmetric.WithTimeout(settings.ExportTimeout),
		)
		metricOptions = append(metricOptions, sdkmetric.WithReader(reader))
	}
	provider.metric = sdkmetric.NewMeterProvider(metricOptions...)
	return provider, nil
}

// Resource returns immutable process/service identity for deterministic tests
// and downstream instrumentation registration.
func (provider *Provider) Resource() *resource.Resource {
	if provider == nil || provider.resource == nil {
		return resource.Empty()
	}
	return provider.resource
}

// Tracer returns an isolated tracer and never depends on the global OTel state.
func (provider *Provider) Tracer(name string) trace.Tracer {
	if provider == nil || provider.trace == nil {
		return tracenoop.NewTracerProvider().Tracer(name)
	}
	return provider.trace.Tracer(name)
}

// Meter returns an isolated meter and never depends on the global OTel state.
func (provider *Provider) Meter(name string) metric.Meter {
	if provider == nil || provider.metric == nil {
		return metricnoop.NewMeterProvider().Meter(name)
	}
	return provider.metric.Meter(name)
}

// ForceFlush synchronously requests pending telemetry export within the bounded
// P01.07 export timeout. Exporter diagnostics remain a private cause.
func (provider *Provider) ForceFlush(ctx context.Context) error {
	if provider == nil || !provider.settings.Enabled {
		return nil
	}
	bounded, cancel := boundedContext(ctx, provider.settings.ExportTimeout)
	defer cancel()

	var failures []error
	if provider.trace != nil {
		if err := provider.trace.ForceFlush(bounded); err != nil {
			failures = append(failures, err)
		}
	}
	if provider.metric != nil {
		if err := provider.metric.ForceFlush(bounded); err != nil {
			failures = append(failures, err)
		}
	}
	if len(failures) == 0 {
		return nil
	}
	return lifecycleFailure(bounded, errors.Join(failures...), "observability force flush failed")
}

// Shutdown flushes/releases telemetry providers once, within the configured
// bounded shutdown timeout. Observability failure never changes business state.
func (provider *Provider) Shutdown(ctx context.Context) error {
	if provider == nil || !provider.settings.Enabled {
		return nil
	}
	provider.shutdownOnce.Do(func() {
		bounded, cancel := boundedContext(ctx, provider.settings.ShutdownTimeout)
		defer cancel()

		var failures []error
		if provider.trace != nil {
			if err := provider.trace.Shutdown(bounded); err != nil {
				failures = append(failures, err)
			}
		}
		if provider.metric != nil {
			if err := provider.metric.Shutdown(bounded); err != nil {
				failures = append(failures, err)
			}
		}
		if len(failures) != 0 {
			provider.shutdownErr = lifecycleFailure(bounded, errors.Join(failures...), "observability shutdown failed")
		}
	})
	return provider.shutdownErr
}

func boundedContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, timeout)
}

func lifecycleFailure(ctx context.Context, cause error, title string) error {
	if errors.Is(ctx.Err(), context.Canceled) {
		return safeWrappedFailure(cause, codeLifecycleCanceled, failure.CategoryDependency, title)
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return safeWrappedFailure(cause, codeLifecycleTimeout, failure.CategoryTimeout, title)
	}
	return safeWrappedFailure(cause, codeLifecycleFailed, failure.CategoryDependency, title)
}
