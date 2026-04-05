package main

import (
	"context"
	"log"
	"time"

	thesilencer_pb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := thesilencer_pb.NewTheSilencerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GeneratePayload(ctx, &thesilencer_pb.GeneratePayloadRequest{
		AttackType: "sqli",
		Complexity: "basic",
		Technology: "php",
		Url:        "http://example.com",
	})
	if err != nil {
		log.Fatalf("could not generate payload: %v", err)
	}
	log.Printf("Payload: %s, Risk: %s, Test Cases: %v", r.GetPayload(), r.GetRiskRating(), r.GetTestCases())
}
