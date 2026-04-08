package client

import (
	"context"
	"log"
	"time"

	thesilencer_pb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TheSilencerClient struct {
	client thesilencer_pb.TheSilencerAdapterClient
}

func NewTheSilencerClient(addr string) (*TheSilencerClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &TheSilencerClient{
		client: thesilencer_pb.NewTheSilencerAdapterClient(conn),
	}, nil
}

func (c *TheSilencerClient) ObfuscateFile(inputPath, technique string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.ObfuscateFile(ctx, &thesilencer_pb.ObfuscateRequest{
		InputPath:  inputPath,
		Technique:  technique,
		OutputPath: "",
	})
	if err != nil {
		return "", err
	}
	return resp.OutputPath, nil
}

func (c *TheSilencerClient) CleanTraces(paths []string, method string) (*thesilencer_pb.CleanupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.CleanTraces(ctx, &thesilencer_pb.CleanupRequest{
		Paths:  paths,
		Method: method,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *TheSilencerClient) HealthCheck() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.HealthCheck(ctx, &thesilencer_pb.Empty{})
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}
