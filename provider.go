package otelinit

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// Provider is the provider for tracing.
type Provider struct {
	ctx                  context.Context
	resource             *resource.Resource
	traceExporter        sdktrace.SpanExporter
	batchProcessorOption sdktrace.BatchSpanProcessorOption
	traceProviderOption  sdktrace.TracerProviderOption
	Shutdown             func() error
}

// NewProvider creates a new provider with default options.
func NewProvider(
	ctx context.Context,
	serviceName string,
	options ...func(*Provider) error,
) (*Provider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("creating stdout exporter: %w", err)
	}

	batchProcessorOption := sdktrace.WithMaxExportBatchSize(1)

	traceProviderOption := sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample()))

	pvd := &Provider{
		ctx:                  ctx,
		resource:             res,
		traceExporter:        traceExporter,
		batchProcessorOption: batchProcessorOption,
		traceProviderOption:  traceProviderOption,
		Shutdown:             nil,
	}

	// push options
	for _, option := range options {
		if err := option(pvd); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return pvd, nil
}

// Init initializes the provider.
func (pvd *Provider) Init() {
	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(pvd.traceExporter, pvd.batchProcessorOption)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
		sdktrace.WithResource(pvd.resource),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
			jaeger.Jaeger{},
		),
	)

	pvd.Shutdown = func() error {
		// Shutdown will flush any remaining spans and shut down the exporter.
		if err := tracerProvider.Shutdown(pvd.ctx); err != nil {
			return fmt.Errorf("failed to shutdown TracerProvider: %w", err)
		}

		return nil
	}
}
