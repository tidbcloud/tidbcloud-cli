package telemetry

import (
	"context"
	"time"
)

var contextKey = telemetryContextKey{}

type telemetryContextKey struct{}

type telemetryContextValue struct {
	startTime time.Time
}

func NewTelemetryContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, telemetryContextValue{
		startTime: time.Now(),
	})
}
