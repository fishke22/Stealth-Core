package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/OpenClaw-Security/Stealth-Core/core/orchestrator" // Import Orchestrator core
	corepb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core"          // Import core proto
	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting Chimera Orchestrator...")

	// Create Orchestrator instance
	orch, err := orchestrator.NewOrchestrator()
	if err != nil {
		log.Fatalf("Failed to initialize Orchestrator: %v", err)
	}
	defer orch.Close() // Ensure all gRPC connections are closed on exit

	// Get listen address
	listenAddr := os.Getenv("ORCHESTRATOR_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":50000" // Default Orchestrator listen port
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", listenAddr, err)
	}

	grpcServer := grpc.NewServer()
	corepb.RegisterOrchestratorServer(grpcServer, orch) // Register Orchestrator service (Changed proto.RegisterOrchestratorServer to corepb.RegisterOrchestratorServer)

	log.Printf("Chimera Orchestrator listening on %s", listenAddr)

	// Start gRPC server in a separate goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down Chimera Orchestrator gracefully...")
	grpcServer.GracefulStop() // Graceful shutdown of gRPC server
	log.Println("Chimera Orchestrator stopped.")
}
