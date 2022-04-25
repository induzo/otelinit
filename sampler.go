package otelinit

import sdktrace "go.opentelemetry.io/otel/sdk/trace"

func (pvd *Provider) WithSampler(opt sdktrace.TracerProviderOption) func(*Provider) error {
	return func(pvd *Provider) error {
		pvd.traceProviderOption = opt

		return nil
	}
}
