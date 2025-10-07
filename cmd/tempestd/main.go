package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Yeah114/tempest-core/launcher"
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

	exited := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server, err := launcher.Start(ctx, launcher.Options{
		Address: address,
		Port:    port,
		Callback: func() {
			log.Printf("tempest-core server has stopped")
			close(exited)
		},
	})
	if err != nil {
		log.Fatalf("failed to start launcher: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("received signal %s, shutting down", sig)
		server.Stop()
	}()

	log.Printf("tempest-core listening on %s", listenAddr)
	<-exited
}