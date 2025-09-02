package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mastercom-service/internal/config"
	"mastercom-service/internal/handlers"
	"mastercom-service/pkg/logger"
	"mastercom-service/pkg/middleware"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize Datadog configuration
	ddConfig := config.LoadDatadogConfig()

	// Start Datadog tracer
	ddConfig.StartTracer()
	defer ddConfig.Stop()

	// Start Datadog profiler
	ddConfig.StartProfiler()
	defer ddConfig.Stop()

	// Initialize Datadog logger
	logger := logger.NewDatadogLogger()
	logger.SetLevel(getLogLevel(cfg.LogLevel))

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize handlers
	handlers.InitHandlers(logger)
	handlers.InitDocumentHandlers(logger)

	// Initialize router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.DatadogMiddleware())

	// Template-style health check endpoint
	router.GET("/__ops/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		span := tracer.StartSpan("health.check", tracer.ResourceName("HealthCheck"))
		defer span.Finish()

		logger.InfoWithSpan(span, "Health check requested", nil)

		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"service":     "mastercom-service",
			"version":     "v1.0.0",
			"dd_trace_id": span.Context().TraceID(),
			"dd_span_id":  span.Context().SpanID(),
		})
	})

	// API routes
	api := router.Group("/api/v6")
	{
		// Case Filing endpoints
		cases := api.Group("/cases")
		{
			cases.POST("", handlers.CreateCase)
			cases.GET("", handlers.ListCases)
			cases.GET("/:id", handlers.GetCase)
			cases.PUT("/:id", handlers.UpdateCase)
			cases.DELETE("/:id", handlers.DeleteCase)
		}

		// Document endpoints
		documents := api.Group("/documents")
		{
			documents.POST("", handlers.UploadDocument)
			documents.GET("/:id", handlers.GetDocument)
			documents.DELETE("/:id", handlers.DeleteDocument)
		}
	}

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting MasterCom service", logrus.Fields{
			"port":        cfg.Port,
			"environment": cfg.Environment,
			"dd_service":  ddConfig.ServiceName,
			"dd_env":      ddConfig.Environment,
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", logrus.Fields{"error": err.Error()})
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...", nil)

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", logrus.Fields{"error": err.Error()})
	}

	logger.Info("Server exited", nil)
}

func getLogLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}
