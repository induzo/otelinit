package otelinit

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func WithGRPCExporter(ctx context.Context, collectorTarget string) func(*Provider) error {
	return func(pvd *Provider) error {
		client := otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(collectorTarget),
			otlptracegrpc.WithInsecure(),
		)

		traceExporter, err := otlptrace.New(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to create grpc traceExporter: %w", err)
		}

		pvd.traceExporter = traceExporter

		return nil
	}
}
