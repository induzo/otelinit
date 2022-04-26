// This package allows you to init and enable tracing in your app
package otelinit

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestInitProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		options []func(*provider) error
		wantErr bool
	}{
		{
			name:    "expecting traces if it is a correct writer",
			options: nil,
			wantErr: false,
		},
		{
			name: "expecting error at provider new",
			options: []func(*provider) error{
				func(pvd *provider) error { return errors.New("error") },
			},
			wantErr: true,
		},
		{
			name: "expecting error at provider init",
			options: []func(*provider) error{
				func(pvd *provider) error {
					pvd.resourceOptions = append(
						pvd.resourceOptions,
						resource.WithDetectors(&testBadDetector{schemaURL: "https://opentelemetry.io/schemas/1.4.0"}),
					)
					pvd.resourceOptions = append(
						pvd.resourceOptions,
						resource.WithDetectors(&testBadDetector{schemaURL: "https://opentelemetry.io/schemas/1.3.0"}),
					)

					return nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			ctx := context.Background()

			tt.options = append(tt.options, WithWriterTraceExporter(buf))
			tt.options = append(tt.options, WithBatchSize(1))

			sd, err := InitProvider(
				ctx, "test",
				tt.options...,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitProvider() expected error %t got %v", tt.wantErr, err)

				return
			}

			if !tt.wantErr {
				defer func() {
					_ = sd()
				}()

				tracer := otel.Tracer("test-tracer")

				// work begins
				_, span := tracer.Start(ctx, "t")
				span.End()

				time.Sleep(10 * time.Millisecond)

				trs, _ := io.ReadAll(buf)
				if len(trs) == 0 {
					t.Errorf("no traces")
				}
			}
		})
	}
}