package main

import (
	"context"
	"flag"
	"os"

	"github.com/gardenbed/basil/config"
	"github.com/gardenbed/basil/graceful"
	"github.com/gardenbed/basil/health"
	"github.com/gardenbed/basil/telemetry"
	grpctelemetry "github.com/gardenbed/basil/telemetry/grpc"

	"grpc-service-horizontal/internal/controller/greeting"
	"grpc-service-horizontal/internal/gateway/github"
	"grpc-service-horizontal/internal/handler"
	"grpc-service-horizontal/internal/repository/usercache"
	"grpc-service-horizontal/internal/server"
	"grpc-service-horizontal/metadata"
)

// Configurations
var configs = struct {
	HTTPPort               uint16
	GRPCPort               uint16
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
	GRPCPort:               9090,
	Name:                   "grpc-service-horizontal",
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
			telemetry.WithOpenTelemetry(true, true, configs.OpenTelemetryCollector, nil),
		)
	}

	probe := telemetry.NewProbe(probeOpts...)
	telemetry.Set(probe)

	defer func() {
		_ = probe.Close(ctx)
	}()

	serverInterceptor := grpctelemetry.NewServerInterceptor(probe, grpctelemetry.Options{})

	// CREATE GATEWAYS

	githubGateway, err := github.NewGateway()
	if err != nil {
		probe.Logger().Error("failed to create github gateway", "error", err)
		panic(err)
	}

	// CREATE REPOSITORIES

	usercacheRepository, err := usercache.NewRepository(configs.RedisAddress)
	if err != nil {
		probe.Logger().Error("failed to create user cache repository", "error", err)
		panic(err)
	}

	// CREATE CONTROLLERS

	greetingController, err := greeting.NewController(githubGateway, usercacheRepository)
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

	grpcServer, err := server.NewGRPC(greetingHandler, server.GRPCOptions{
		Port:    configs.GRPCPort,
		Options: serverInterceptor.ServerOptions(),
	})

	if err != nil {
		probe.Logger().Error("failed to create grpc server", "error", err)
		panic(err)
	}

	// Create an HTTP health handler for health checking the service by external systems
	health.SetLogger(probe.Logger())
	health.RegisterChecker(githubGateway, usercacheRepository)
	healthHandler := health.HandlerFunc()

	httpServer, err := server.NewHTTP(healthHandler, server.HTTPOptions{
		Port: configs.HTTPPort,
	})

	if err != nil {
		probe.Logger().Error("failed to create http server", "error", err)
		panic(err)
	}

	// Gracefully, connect the clients and start the servers
	// Gracefully, retry the lost connections
	// Gracefully, disconnect the clients and shutdown the servers on termination signals
	graceful.SetLogger(probe.Logger())
	graceful.RegisterClient(githubGateway, usercacheRepository)
	graceful.RegisterServer(grpcServer, httpServer)
	code := graceful.StartAndWait()

	os.Exit(code)
}
