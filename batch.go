package otelinit

import sdktrace "go.opentelemetry.io/otel/sdk/trace"

func (pvd *Provider) WithBatchSize(size int) func(*Provider) error {
	return func(pvd *Provider) error {
		pvd.batchProcessorOption = sdktrace.WithMaxExportBatchSize(size)

		return nil
	}
}
