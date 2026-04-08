package main

import (
	"log"
	"net"

	"github.com/OpenClaw-Security/Stealth-Core/adapters/thesilencer/server"
	thesilencer_pb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	thesilencer_pb.RegisterTheSilencerAdapterServer(s, &server.Server{})

	log.Printf("TheSilencer Adapter listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
