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

	"grpc-service/internal/client"
	"grpc-service/internal/server"
	"grpc-service/internal/service/greeting"
	"grpc-service/metadata"
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
	Name:                   "grpc-service",
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
	defer probe.Close(ctx)

	serverInterceptor := grpctelemetry.NewServerInterceptor(probe, grpctelemetry.Options{})

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

	grpcServer, err := server.NewGRPC(greetingService, server.GRPCOptions{
		Port:    configs.GRPCPort,
		Options: serverInterceptor.ServerOptions(),
	})

	if err != nil {
		probe.Logger().Error("failed to create grpc server", "error", err)
		panic(err)
	}

	// Create an HTTP health handler for health checking the service by external systems
	health.SetLogger(probe.Logger())
	health.RegisterChecker(httpClient, redisClient)
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
	graceful.RegisterClient(httpClient, redisClient)
	graceful.RegisterServer(grpcServer, httpServer)
	code := graceful.StartAndWait()

	os.Exit(code)
}
