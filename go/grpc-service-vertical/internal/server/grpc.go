package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"grpc-service-vertical/internal/idl/greetingpb"
)

const defaultGRPCPort = 9090

// grpcServer is an interface for grpc.Server struct.
type grpcServer interface {
	GracefulStop()
	Serve(net.Listener) error
}

// GRPC is a grpc server implementing the graceful.Server interface.
type GRPC struct {
	addr   string
	server grpcServer
}

// GRPCOptions are optional settings for creating a grpc server.
type GRPCOptions struct {
	// The port number for the gRPC server.
	// The default port number is 9090.
	Port uint16
	// A TLS certificate for the gRPC server identity.
	TLSCert *tls.Certificate
	// A pool of certificate authorities for verifying gRPC client identities.
	ClientCA *x509.CertPool
	// Additional options for the gRPC server.
	Options []grpc.ServerOption
}

// NewGRPC creates a new grpc Server.
func NewGRPC(greetingService greetingpb.GreetingServiceServer, opts GRPCOptions) (*GRPC, error) {
	if opts.Port == 0 {
		opts.Port = defaultGRPCPort
	}

	tlsConfig := new(tls.Config)
	if opts.TLSCert != nil {
		tlsConfig.Certificates = []tls.Certificate{*opts.TLSCert}
	}
	if opts.ClientCA != nil {
		tlsConfig.ClientCAs = opts.ClientCA
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	grpcOpts := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}
	grpcOpts = append(grpcOpts, opts.Options...)

	server := grpc.NewServer(grpcOpts...)
	greetingpb.RegisterGreetingServiceServer(server, greetingService)

	return &GRPC{
		addr:   fmt.Sprintf(":%d", opts.Port),
		server: server,
	}, nil
}

// String returns the name of the server.
func (s *GRPC) String() string {
	return "grpc-server"
}

// ListenAndServe starts listening for incoming requests synchronously.
// It blocks the current goroutine until an error is returned.
func (s *GRPC) ListenAndServe() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	// Synchronous/Blocking
	return s.server.Serve(lis)
}

// Shutdown gracefully stops the server.
// It stops accepting new conenctions and blocks the current goroutine until all the pending requests are completed.
// If the context is cancelled, an error will be returned.
func (s *GRPC) Shutdown(ctx context.Context) error {
	done := make(chan struct{}, 1)
	go func() {
		s.server.GracefulStop()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
