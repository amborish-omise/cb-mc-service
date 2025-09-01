package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// DatadogMiddleware adds Datadog tracing to HTTP requests
func DatadogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start a new span for the HTTP request
		span := tracer.StartSpan("http.request",
			tracer.ResourceName(c.Request.Method+" "+c.Request.URL.Path),
			tracer.SpanType("web"),
			tracer.Tag("http.method", c.Request.Method),
			tracer.Tag("http.url", c.Request.URL.String()),
			tracer.Tag("http.user_agent", c.Request.UserAgent()),
			tracer.Tag("http.request_id", c.GetHeader("X-Request-ID")),
		)
		defer span.Finish()

		// Add span to context
		ctx := context.WithValue(c.Request.Context(), "dd_span", span)
		c.Request = c.Request.WithContext(ctx)

		// Add trace headers to response
		c.Header("X-Datadog-Trace-ID", fmt.Sprintf("%d", span.Context().TraceID()))
		c.Header("X-Datadog-Span-ID", fmt.Sprintf("%d", span.Context().SpanID()))

		// Process request
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// Add response tags to span
		span.SetTag("http.status_code", c.Writer.Status())
		span.SetTag("http.response_time", duration.String())
		span.SetTag("http.response_size", c.Writer.Size())

		// Add error tag if status code indicates error
		if c.Writer.Status() >= 400 {
			span.SetTag("error", true)
			span.SetTag("error.message", http.StatusText(c.Writer.Status()))
		}
	}
}

// GetSpanFromContext retrieves the Datadog span from context
func GetSpanFromContext(ctx context.Context) tracer.Span {
	if span, ok := ctx.Value("dd_span").(tracer.Span); ok {
		return span
	}
	return nil
}

// AddSpanTag adds a tag to the current span
func AddSpanTag(ctx context.Context, key, value string) {
	if span := GetSpanFromContext(ctx); span != nil {
		span.SetTag(key, value)
	}
}

// AddSpanError adds an error to the current span
func AddSpanError(ctx context.Context, err error) {
	if span := GetSpanFromContext(ctx); span != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
	}
}
