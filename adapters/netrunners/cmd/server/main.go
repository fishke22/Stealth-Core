package main

import (
	"log"
	"net"

	"github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners"
	"google.golang.org/grpc"

	netrunnerslib "github.com/OpenClaw-Security/Stealth-Core/pkg/adapters/netrunners"
)

func main() {
	server := netrunnerslib.NewNetRunnersServer()

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	netrunners.RegisterNetRunnersAdapterServer(grpcServer, server)

	log.Printf("NetRunners Adapter Server listening on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
