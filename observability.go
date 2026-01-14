package kie

import "context"

// Metrics is the interface for collecting metrics.
type Metrics interface {
	// IncCounter increments a counter metric.
	IncCounter(name string, labels ...string)

	// ObserveHistogram records a value for a histogram metric.
	ObserveHistogram(name string, value float64, labels ...string)
}

// Tracer is the interface for distributed tracing.
type Tracer interface {
	// Start creates a new span.
	Start(ctx context.Context, name string) (context.Context, Span)
}

// Span represents a single unit of work in a trace.
type Span interface {
	// SetAttribute sets an attribute on the span.
	SetAttribute(key string, value interface{})

	// RecordError records an error on the span.
	RecordError(err error)

	// End finishes the span.
	End()
}

// WithMetrics enables metrics collection with the given metrics collector.
func WithMetrics(metrics Metrics) Option {
	return func(c *Client) {
		c.metrics = metrics
	}
}

// WithTracer enables distributed tracing with the given tracer.
func WithTracer(tracer Tracer) Option {
	return func(c *Client) {
		c.tracer = tracer
	}
}
