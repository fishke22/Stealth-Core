package main

import (
	"context"
	"fmt"
	"log"
	"net"

	core "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core"
	netrunnerspb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners"
	"google.golang.org/grpc"
)

type NetRunnersServer struct {
	netrunnerspb.UnimplementedNetRunnersAdapterServer
}

func (s *NetRunnersServer) Execute(ctx context.Context, req *core.OperationRequest) (*core.OperationResponse, error) {
	log.Printf("Executing NetRunners operation: %s", req.Command)
	
	switch req.Command {
	case "setup_listener":
		return s.setupListener(ctx, req.Params)
	case "exfiltrate_data":
		return s.ExfiltrateData(ctx, req.Params)
	case "network_scan":
		return s.networkScan(ctx, req.Params)
	default:
		return nil, fmt.Errorf("unknown command: %s", req.Command)
	}
}

func (s *NetRunnersServer) setupListener(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	protocol := params["protocol"]
	if protocol == "" {
		protocol = "https"
	}

	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Listener setup for protocol %s", protocol),
		Artifacts: map[string]string{
			"protocol": protocol,
			"status":  "listening",
		},
	}, nil
}

func (s *NetRunnersServer) exfiltrate极Data(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	dataPath := params["data_path"]
	if dataPath == "" {
		return nil, fmt.Errorf("data_path parameter is required")
	}

	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Data exfiltration completed for %s", dataPath),
		Artifacts: map[string]string{
			"data_path": dataPath,
			"status":   "completed",
		},
	}, nil
}

func (s *NetRunnersServer) networkScan(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	target := params["target"]
	if target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}

	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Network scan completed for %s", target),
		Artifacts: map[string]string{
			"target": target,
			"status": "completed",
		},
	}, nil
}

func (s *NetRunnersServer) GetCapabilities(ctx context.Context, empty *core.Empty) (*core.CapabilityList, error) {
	return &core.CapabilityList{
		Commands: []string{
			"setup_listener",
			"exfiltrate_data",
			"network_scan",
		},
	}, nil
}

func (s *NetRunnersServer) HealthCheck(ctx context.Context, empty *core.Empty) (*core.HealthStatus, error) {
	return &core.HealthStatus{
		Status:  "healthy",
		Version: "1.0",
	}, nil
}

func main() {
	server := &NetRunnersServer{}
	
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	netrunnerspb.RegisterNetRunnersAdapterServer(grpcServer, server)
	
	log.Printf("NetRunners Adapter Server listening on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}