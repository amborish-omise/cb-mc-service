package config

import (
	"fmt"
	"os"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

type DatadogConfig struct {
	Enabled     bool    `mapstructure:"DD_ENABLED"`
	ServiceName string  `mapstructure:"DD_SERVICE"`
	Environment string  `mapstructure:"DD_ENV"`
	Version     string  `mapstructure:"DD_VERSION"`
	AgentHost   string  `mapstructure:"DD_AGENT_HOST"`
	AgentPort   string  `mapstructure:"DD_AGENT_PORT"`
	LogLevel    string  `mapstructure:"DD_LOG_LEVEL"`
	TraceRate   float64 `mapstructure:"DD_TRACE_SAMPLE_RATE"`
}

func LoadDatadogConfig() *DatadogConfig {
	config := &DatadogConfig{
		Enabled:     getEnvBool("DD_ENABLED", true),
		ServiceName: getEnv("DD_SERVICE", "mastercom-service"),
		Environment: getEnv("DD_ENV", "development"),
		Version:     getEnv("DD_VERSION", "1.0.0"),
		AgentHost:   getEnv("DD_AGENT_HOST", "localhost"),
		AgentPort:   getEnv("DD_AGENT_PORT", "8126"),
		LogLevel:    getEnv("DD_LOG_LEVEL", "info"),
		TraceRate:   getEnvFloat("DD_TRACE_SAMPLE_RATE", 1.0),
	}

	return config
}

func (c *DatadogConfig) StartTracer() {
	if !c.Enabled {
		return
	}

	// Set Datadog tracer options
	tracer.Start(
		tracer.WithService(c.ServiceName),
		tracer.WithEnv(c.Environment),
		tracer.WithAgentAddr(c.AgentHost+":"+c.AgentPort),
		tracer.WithLogStartup(false),
	)
}

func (c *DatadogConfig) StartProfiler() {
	if !c.Enabled {
		return
	}

	// Start Datadog profiler
	profiler.Start(
		profiler.WithService(c.ServiceName),
		profiler.WithEnv(c.Environment),
		profiler.WithAgentAddr(c.AgentHost+":"+c.AgentPort),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.BlockProfile,
			profiler.MutexProfile,
		),
	)
}

func (c *DatadogConfig) Stop() {
	if !c.Enabled {
		return
	}

	profiler.Stop()
	tracer.Stop()
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if f, err := parseFloat(value); err == nil {
			return f
		}
	}
	return defaultValue
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
