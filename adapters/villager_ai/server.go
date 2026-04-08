package villagerai

import (
	"context"
	"fmt"
	"log"
	"net"

	core "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core"
	thesilencer "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer"
	villagerpb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/villager_ai"
	"google.golang.org/grpc"
)

type VillagerAIServer struct {
	villagerpb.UnimplementedVillagerAIAdapterServer
}

func NewVillagerAIServer() (*VillagerAIServer, error) {
	return &VillagerAIServer{}, nil
}

func (s *VillagerAIServer) Execute(ctx context.Context, req *core.OperationRequest) (*core.OperationResponse, error) {
	log.Printf("Executing Villager-AI fusion operation: %s", req.Command)

	switch req.Command {
	case "generate_stealthy_payload":
		return s.generateStealthyPayloadWrapper(ctx, req.Params)
	case "launch_covert_phishing":
		return s.launchCovertPhishingWrapper(ctx, req.Params)
	case "automated_exfiltration":
		return s.automatedExfiltrationWrapper(ctx, req.Params)
	case "deploy_c2_persistence":
		return s.deployC2PersistenceWrapper(ctx, req.Params)
	default:
		return nil, fmt.Errorf("unknown command for Villager-AI: %s", req.Command)
	}
}

// generateStealthyPayload 融合功能：生成载荷 -> 混淆 -> 检测
func (s *VillagerAIServer) GenerateStealthyPayload(ctx context.Context, req *villagerpb.GenerateStealthyPayloadRequest) (*villagerpb.GenerateStealthyPayloadResponse, error) {
	log.Printf("Villager-AI: Generating stealthy payload of type %s with obfuscation %s", req.BasePayloadType, req.ObfuscationTechnique)

	// 模擬生成基礎載荷
	originalHash := "abc123"

	// 模擬混淆載荷
	obfuscatedHash := "def456"

	// 模擬檢測結果
	detectionResp := &thesilencer.DetectionResponse{
		Detected: false,
		Score:    0,
	}

	return &villagerpb.GenerateStealthyPayloadResponse{
		FinalPayloadPath: "/tmp/obfuscated_payload.bin",
		OriginalHash:     originalHash,
		ObfuscatedHash:   obfuscatedHash,
		DetectionResults: detectionResp,
	}, nil
}

// 封装成 core.OperationResponse
func (s *VillagerAIServer) generateStealthyPayloadWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &villagerpb.GenerateStealthyPayloadRequest{
		BasePayloadType:      params["base_payload_type"],
		ObfuscationTechnique: params["obfuscation_technique"],
		HexstrikeOptions:     map[string]string{"lhost": params["lhost"], "lport": params["lport"]},
		ThesilencerOptions:   map[string]string{"level": "high"},
	}
	resp, err := s.GenerateStealthyPayload(ctx, req)
	if err != nil {
		return &core.OperationResponse{Success: false, Error: err.Error()}, nil
	}
	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Generated stealthy payload: %s", resp.FinalPayloadPath),
		Artifacts: map[string]string{
			"final_payload_path": resp.FinalPayloadPath,
			"detection_score":    fmt.Sprintf("%d", resp.DetectionResults.Score),
		},
	}, nil
}

// launchCovertPhishing 融合功能：投递载荷
func (s *VillagerAIServer) LaunchCovertPhishing(ctx context.Context, req *villagerpb.LaunchCovertPhishingRequest) (*villagerpb.LaunchCovertPhishingResponse, error) {
	log.Printf("Villager-AI: Launching covert phishing campaign '%s' with payload %s", req.CampaignName, req.PayloadPath)

	// 模擬發送釣魚郵件
	return &villagerpb.LaunchCovertPhishingResponse{
		Success:          true,
		SentCount:        10,
		FailedRecipients: []string{},
		CampaignId:       fmt.Sprintf("campaign-%s-%d", req.CampaignName, 1234567890),
	}, nil
}

func (s *VillagerAIServer) launchCovertPhishingWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &villagerpb.LaunchCovertPhishingRequest{
		CampaignName:          params["campaign_name"],
		PayloadPath:           params["payload_path"],
		PhishingTemplate:      params["phishing_template"],
		Recipients:            []string{params["recipient"]},
		HexstrikeEmailOptions: map[string]string{"subject": params["subject"], "smtp_server": params["smtp_server"]},
	}
	resp, err := s.LaunchCovertPhishing(ctx, req)
	if err != nil {
		return &core.OperationResponse{Success: false, Error: err.Error()}, nil
	}
	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Launched phishing campaign %s. Sent to %d recipients.", resp.CampaignId, resp.SentCount),
		Artifacts: map[string]string{
			"campaign_id": resp.CampaignId,
			"sent_count":  fmt.Sprintf("%d", resp.SentCount),
		},
	}, nil
}

// automatedExfiltration 融合功能：数据渗出
func (s *VillagerAIServer) AutomatedExfiltration(ctx context.Context, req *villagerpb.AutomatedExfiltrationRequest) (*villagerpb.AutomatedExfiltrationResponse, error) {
	log.Printf("Villager-AI: Automating exfiltration of %s to %s via %s", req.SourceDataPath, req.DestinationEndpoint, req.ExfiltrationMethod)

	// 模擬數據滲出
	return &villagerpb.AutomatedExfiltrationResponse{
		Success:                true,
		ExfiltrationReportPath: "report.json",
		TotalBytesExfiltrated:  1024,
	}, nil
}

func (s *VillagerAIServer) automatedExfiltrationWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &villagerpb.AutomatedExfiltrationRequest{
		SourceDataPath:            params["source_data_path"],
		DestinationEndpoint:       params["destination_endpoint"],
		ExfiltrationMethod:        params["exfiltration_method"],
		NetrunnersOptions:         map[string]string{"encryption": "aes256"},
		ThesilencerCleanupOptions: map[string]string{"method": "secure_delete"},
	}
	resp, err := s.AutomatedExfiltration(ctx, req)
	if err != nil {
		return &core.OperationResponse{Success: false, Error: err.Error()}, nil
	}
	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Automated exfiltration completed. Bytes sent: %d", resp.TotalBytesExfiltrated),
		Artifacts: map[string]string{
			"total_bytes_exfiltrated": fmt.Sprintf("%d", resp.TotalBytesExfiltrated),
		},
	}, nil
}

// DeployCovertC2AndPersistence 融合功能：部署 C2 和持久化
func (s *VillagerAIServer) DeployCovertC2AndPersistence(ctx context.Context, req *villagerpb.DeployC2PersistenceRequest) (*villagerpb.DeployC2PersistenceResponse, error) {
	log.Printf("Villager-AI: Deploying covert C2 and persistence with method %s", req.PersistenceMethod)

	// 模擬部署隱蔽 C2 監聽器
	c2Info := "tcp://example.com:443"

	// 模擬實現持久化
	persistenceReport := "Persistence established successfully"

	return &villagerpb.DeployC2PersistenceResponse{
		Success:           true,
		C2ListenerInfo:    c2Info,
		PersistenceReport: persistenceReport,
	}, nil
}

func (s *VillagerAIServer) deployC2PersistenceWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &villagerpb.DeployC2PersistenceRequest{
		PayloadPath:                 params["payload_path"],
		PersistenceMethod:           params["persistence_method"],
		C2Protocol:                  params["c2_protocol"],
		NetrunnersC2Options:         map[string]string{"domain": params["c2_domain"], "port": params["c2_port"]},
		HexstrikePersistenceOptions: map[string]string{"technique": params["persistence_method"]},
	}
	resp, err := s.DeployCovertC2AndPersistence(ctx, req)
	if err != nil {
		return &core.OperationResponse{Success: false, Error: err.Error()}, nil
	}
	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Deployed C2 listener: %s. Persistence report: %s", resp.C2ListenerInfo, resp.PersistenceReport),
		Artifacts: map[string]string{
			"c2_listener":        resp.C2ListenerInfo,
			"persistence_report": resp.PersistenceReport,
		},
	}, nil
}

func (s *VillagerAIServer) GetCapabilities(ctx context.Context, empty *core.Empty) (*core.CapabilityList, error) {
	return &core.CapabilityList{
		Commands: []string{
			"generate_stealthy_payload",
			"launch_covert_phishing",
			"automated_exfiltration",
			"deploy_c2_persistence",
		},
	}, nil
}

func (s *VillagerAIServer) HealthCheck(ctx context.Context, empty *core.Empty) (*core.HealthStatus, error) {
	return &core.HealthStatus{
		Status:  "healthy",
		Version: "1.0-fusion",
	}, nil
}

// StartVillagerAIServer starts the Villager-AI adapter gRPC server
func StartVillagerAIServer() error {
	server, err := NewVillagerAIServer()
	if err != nil {
		return fmt.Errorf("failed to create Villager-AI server: %w", err)
	}

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	villagerpb.RegisterVillagerAIAdapterServer(grpcServer, server)

	log.Printf("Villager-AI Adapter Server listening on :50054")
	return grpcServer.Serve(lis)
}

// calculateFileHash computes MD5 hash of a file
func calculateFileHash(filePath string) (string, error) {
	// This is a simplified implementation for demonstration
	return "dummy_hash", nil
}
