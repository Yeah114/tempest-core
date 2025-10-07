package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Yeah114/tempest-core/internal/app"
	"github.com/Yeah114/tempest-core/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		address = "0.0.0.0"
		port    = 20919
	)
	flag.StringVar(&address, "a", address, "Bind tempest-core service to a specific TCP/IPv4 address")
	flag.IntVar(&port, "p", port, "Bind tempest-core service to a specific TCP/IPv4 port")
	flag.Parse()

	listenAddr := fmt.Sprintf("%s:%d", address, port)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", listenAddr, err)
	}

	state := app.NewFatalderState()
	services := server.NewServices(state)

	grpcServer := grpc.NewServer()
	services.Register(grpcServer)
	reflection.Register(grpcServer)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("received signal %s, shutting down", sig)
		grpcServer.GracefulStop()
	}()

	log.Printf("tempest-core listening on %s", listenAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server stopped: %v", err)
	}
}
