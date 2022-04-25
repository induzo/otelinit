package otelinit

import (
	"fmt"
	"io"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

func WithDiscardExporter() func(*Provider) error {
	return func(pvd *Provider) error {
		traceExporter, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
		if err != nil {
			return fmt.Errorf("creating discard exporter: %w", err)
		}

		pvd.traceExporter = traceExporter

		return nil
	}
}
