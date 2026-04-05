package main

import (
	"context"
	"log"
	"time"

	hexstrike_pb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/hexstrike"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := hexstrike_pb.NewHexstrikeServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GeneratePayload(ctx, &hexstrike_pb.GeneratePayloadRequest{
		AttackType: "xss",
		Complexity: "advanced",
		Technology: "nodejs",
		Url:        "http://example.com",
	})
	if err != nil {
		log.Fatalf("could not generate payload: %v", err)
	}
	log.Printf("Payload: %s, Risk: %s, Test Cases: %v", r.GetPayload(), r.GetRiskRating(), r.GetTestCases())
}
