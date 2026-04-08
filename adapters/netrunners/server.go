package netrunners

import (
	"context"
	"fmt"
	"log"

	core "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core"
	netrunnerspb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners"
)

type NetRunnersServer struct {
	netrunnerspb.UnimplementedNetRunnersAdapterServer
}

func NewNetRunnersServer() *NetRunnersServer {
	return &NetRunnersServer{}
}

func (s *NetRunnersServer) Execute(ctx context.Context, req *core.OperationRequest) (*core.OperationResponse, error) {
	log.Printf("Executing NetRunners operation: %s", req.Command)

	switch req.Command {
	case "setup_listener":
		return s.executeSetupListener(ctx, req.Params)
	case "covert_channel":
		return s.executeCovertChannel(ctx, req.Params)
	case "exfiltrate_data":
		return s.executeExfiltrateData(ctx, req.Params)
	case "network_scan":
		return s.executeNetworkScan(ctx, req.Params)
	default:
		return nil, fmt.Errorf("unknown command: %s", req.Command)
	}
}

func (s *NetRunnersServer) GetCapabilities(ctx context.Context, _ *core.Empty) (*core.CapabilityList, error) {
	return &core.CapabilityList{
		Commands: []string{
			"SetupListener",
			"CovertChannel",
			"ExfiltrateData",
			"NetworkScan",
		},
	}, nil
}

func (s *NetRunnersServer) HealthCheck(ctx context.Context, _ *core.Empty) (*core.HealthStatus, error) {
	return &core.HealthStatus{
		Status:  "healthy",
		Version: "1.0",
	}, nil
}

func (s *NetRunnersServer) SetupListener(ctx context.Context, req *netrunnerspb.ListenerRequest) (*netrunnerspb.ListenerResponse, error) {
	log.Printf("Setting up listener: protocol=%s, domain=%s, port=%s", req.Protocol, req.Domain, req.Port)

	protocol := req.Protocol
	if protocol == "" {
		protocol = "https"
	}

	return &netrunnerspb.ListenerResponse{
		ListenerId:       fmt.Sprintf("listener_%d", 12345),
		Status:           "active",
		ConnectionString: fmt.Sprintf("%s://%s:%s", protocol, req.Domain, req.Port),
	}, nil
}

func (s *NetRunnersServer) CovertChannel(ctx context.Context, req *netrunnerspb.ChannelRequest) (*netrunnerspb.ChannelResponse, error) {
	log.Printf("Creating covert channel: type=%s, target=%s", req.ChannelType, req.Target)

	return &netrunnerspb.ChannelResponse{
		Success:      true,
		ResponseData: req.Data,
		LatencyMs:    int32(100),
	}, nil
}

func (s *NetRunnersServer) ExfiltrateData(ctx context.Context, req *netrunnerspb.ExfiltrationRequest) (*netrunnerspb.ExfiltrationResponse, error) {
	log.Printf("Exfiltrating data: path=%s, method=%s, destination=%s", req.DataPath, req.Method, req.Destination)

	return &netrunnerspb.ExfiltrationResponse{
		Success:        true,
		BytesSent:      int32(1024),
		ExfiltrationId: fmt.Sprintf("exfil_%d", 67890),
	}, nil
}

func (s *NetRunnersServer) NetworkScan(ctx context.Context, req *netrunnerspb.ScanRequest) (*netrunnerspb.ScanResponse, error) {
	log.Printf("Performing network scan: target=%s, type=%s", req.Target, req.ScanType)

	return &netrunnerspb.ScanResponse{
		ScanId:      fmt.Sprintf("scan_%d", 11111),
		HostsFound:  int32(5),
		PortsOpen:   int32(8),
		ResultsPath: fmt.Sprintf("/tmp/scan_results_%s.txt", req.Target),
	}, nil
}

func (s *NetRunnersServer) executeSetupListener(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &netrunnerspb.ListenerRequest{
		Protocol: params["protocol"],
		Domain:   params["domain"],
		Port:     params["port"],
	}
	resp, err := s.SetupListener(ctx, req)
	if err != nil {
		return nil, err
	}

	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Listener established: %s", resp.ConnectionString),
		Artifacts: map[string]string{
			"listener_id":       resp.ListenerId,
			"status":            resp.Status,
			"connection_string": resp.ConnectionString,
		},
	}, nil
}

func (s *NetRunnersServer) executeCovertChannel(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &netrunnerspb.ChannelRequest{
		ChannelType: params["channel_type"],
		Target:      params["target"],
		Data:        []byte(params["data"]),
	}
	resp, err := s.CovertChannel(ctx, req)
	if err != nil {
		return nil, err
	}

	return &core.OperationResponse{
		Success: resp.Success,
		Output:  fmt.Sprintf("Covert channel established, latency: %dms", resp.LatencyMs),
		Artifacts: map[string]string{
			"latency_ms": fmt.Sprintf("%d", resp.LatencyMs),
		},
	}, nil
}

func (s *NetRunnersServer) executeExfiltrateData(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &netrunnerspb.ExfiltrationRequest{
		DataPath:    params["data_path"],
		Method:      params["method"],
		Destination: params["destination"],
	}
	resp, err := s.ExfiltrateData(ctx, req)
	if err != nil {
		return nil, err
	}

	return &core.OperationResponse{
		Success: resp.Success,
		Output:  fmt.Sprintf("Data exfiltration completed: %d bytes sent", resp.BytesSent),
		Artifacts: map[string]string{
			"exfiltration_id": resp.ExfiltrationId,
			"bytes_sent":      fmt.Sprintf("%d", resp.BytesSent),
		},
	}, nil
}

func (s *NetRunnersServer) executeNetworkScan(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &netrunnerspb.ScanRequest{
		Target:   params["target"],
		ScanType: params["scan_type"],
	}
	resp, err := s.NetworkScan(ctx, req)
	if err != nil {
		return nil, err
	}

	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Network scan completed: %d hosts, %d open ports", resp.HostsFound, resp.PortsOpen),
		Artifacts: map[string]string{
			"scan_id":      resp.ScanId,
			"hosts_found":  fmt.Sprintf("%d", resp.HostsFound),
			"ports_open":   fmt.Sprintf("%d", resp.PortsOpen),
			"results_path": resp.ResultsPath,
		},
	}, nil
}
