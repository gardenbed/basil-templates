package main

import (
	"context"
	"flag"
	"os"

	"github.com/gardenbed/basil/config"
	"github.com/gardenbed/basil/graceful"
	"github.com/gardenbed/basil/health"
	"github.com/gardenbed/basil/httpx"
	"github.com/gardenbed/basil/telemetry"
	httptelemetry "github.com/gardenbed/basil/telemetry/http"

	"http-service-vertical/internal/client"
	"http-service-vertical/internal/server"
	"http-service-vertical/internal/service/greeting"
	"http-service-vertical/metadata"
)

// Configurations
var configs = struct {
	HTTPPort               uint16
	Name                   string
	Provider               string
	Region                 string
	Cluster                string
	Namespace              string
	LogLevel               string
	OpenTelemetryCollector string
	RedisAddress           string
}{
	// Default Values
	HTTPPort:               8080,
	Name:                   "http-service-vertical",
	Provider:               "local",
	Region:                 "local",
	Cluster:                "local",
	Namespace:              "local",
	LogLevel:               "debug",
	OpenTelemetryCollector: "localhost:55680",
	RedisAddress:           "localhost:6379",
}

func main() {
	ctx := context.Background()

	// Get configurations
	_ = config.Pick(&configs)
	flag.Parse()

	// CREATE A TELEMETRY PROBE

	probeOpts := []telemetry.Option{
		telemetry.WithLogger(configs.LogLevel),
		telemetry.WithMetadata(configs.Name, metadata.Version, map[string]string{
			"provider":  configs.Provider,
			"region":    configs.Region,
			"cluster":   configs.Cluster,
			"namespace": configs.Namespace,
		}),
	}

	if configs.OpenTelemetryCollector != "" {
		probeOpts = append(probeOpts,
			telemetry.WithOpenTelemetry(configs.OpenTelemetryCollector, nil),
		)
	}

	probe := telemetry.NewProbe(probeOpts...)
	telemetry.Set(probe)
	defer probe.Close(ctx)

	telemetryMiddleware := httptelemetry.NewMiddleware(probe, httptelemetry.Options{})

	// CREATE CLIENTS

	httpClient := client.NewHTTP()
	redisClient := client.NewRedis(configs.RedisAddress)

	// CREATE SERVICES

	greetingService, err := greeting.NewService(httpClient, redisClient)
	if err != nil {
		probe.Logger().Error("failed to create greeting service", "error", err)
		panic(err)
	}

	// CREATE SERVERS

	// Create an HTTP health handler for health checking the service by external systems
	health.SetLogger(probe.Logger())
	health.RegisterChecker(httpClient, redisClient)
	healthHandler := health.HandlerFunc()

	httpServer, err := server.NewHTTP(healthHandler, greetingService, server.HTTPOptions{
		Port: configs.HTTPPort,
		Middleware: []httpx.Middleware{
			telemetryMiddleware,
		},
	})

	if err != nil {
		probe.Logger().Error("failed to create http server", "error", err)
		panic(err)
	}

	// Gracefully, connect the clients and start the servers
	// Gracefully, retry the lost connections
	// Gracefully, disconnect the clients and shutdown the servers on termination signals
	graceful.SetLogger(probe.Logger())
	graceful.RegisterClient(httpClient, redisClient)
	graceful.RegisterServer(httpServer)
	code := graceful.StartAndWait()

	os.Exit(code)
}
