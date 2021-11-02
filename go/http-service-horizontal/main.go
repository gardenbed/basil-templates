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

	"http-service-horizontal/internal/controller/greeting"
	"http-service-horizontal/internal/gateway/translate"
	"http-service-horizontal/internal/handler"
	"http-service-horizontal/internal/repository/greetingcache"
	"http-service-horizontal/internal/server"
	"http-service-horizontal/metadata"
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
	Name:                   "http-service-horizontal",
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

	// CREATE GATEWAYS

	translateGateway, err := translate.NewGateway()
	if err != nil {
		probe.Logger().Error("failed to create translate gateway", "error", err)
		panic(err)
	}

	// CREATE REPOSITORIES

	greetingcacheRepository, err := greetingcache.NewRepository(configs.RedisAddress)
	if err != nil {
		probe.Logger().Error("failed to create greeting cache repository", "error", err)
		panic(err)
	}

	// CREATE CONTROLLERS

	greetingController, err := greeting.NewController(translateGateway, greetingcacheRepository)
	if err != nil {
		probe.Logger().Error("failed to create greeting controller", "error", err)
		panic(err)
	}

	// CREATE HANDLERS

	greetingHandler, err := handler.NewGreetingHandler(greetingController)
	if err != nil {
		probe.Logger().Error("failed to create greetting handler", "error", err)
		panic(err)
	}

	// CREATE SERVERS

	// Create an HTTP health handler for health checking the service by external systems
	health.SetLogger(probe.Logger())
	health.RegisterChecker(translateGateway, greetingcacheRepository)
	healthHandler := health.HandlerFunc()

	httpServer, err := server.NewHTTP(healthHandler, greetingHandler, server.HTTPOptions{
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
	graceful.RegisterClient(translateGateway, greetingcacheRepository)
	graceful.RegisterServer(httpServer)
	code := graceful.StartAndWait()

	os.Exit(code)
}
