package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	core "github.com/OpenClaw-Security/Stealth-Core/pkg/proto"
	redteam_toolspb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/redteam_tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RedteamToolsServer struct {
	redteam_toolspb.UnimplementedRedteamToolsAdapterServer
	toolBasePath string // Hexstrike-redteam 的根目录
}

func NewRedteamToolsServer() *RedteamToolsServer {
	return &RedteamToolsServer{
		toolBasePath: "/tools/Hexstrike-redteam", // 默认路径
	}
}

func (s *RedteamToolsServer) Execute(ctx context.Context, req *core.OperationRequest) (*core.OperationResponse, error) {
	log.Printf("Executing RedteamTools operation: %s", req.Command)

	switch req.Command {
	case "enumerate_smb":
		return s.enumerateSMBWrapper(ctx, req.Params)
	case "perform_lateral_movement":
		return s.performLateralMovementWrapper(ctx, req.Params)
	case "dump_credentials":
		return s.dumpCredentialsWrapper(ctx, req.Params)
	case "scan_vulnerabilities":
		return s.scanVulnerabilitiesWrapper(ctx, req.Params)
	default:
		return nil, fmt.Errorf("unknown command for RedteamTools: %s", req.Command)
	}
}

// 辅助函数：运行命令行工具
func (s *RedteamToolsServer) runToolCommand(ctx context.Context, toolScript string, args []string) (string, error) {
	fullPath := filepath.Join(s.toolBasePath, toolScript)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("redteam tool not found: %s", fullPath)
	}

	cmdArgs := []string{fullPath}
	cmdArgs = append(cmdArgs, args...)

	log.Printf("Running redteam tool command: %s %v", "python3" /*or bash/specific interpreter*/, cmdArgs)
	cmd := exec.Command("python3", cmdArgs...) // 假设大部分工具是Python脚本
	
	// 设置超时
	if timeoutStr, ok := ctx.Value("timeout").(string); ok {
		if duration, err := time.ParseDuration(timeoutStr); err == nil {
			timer := time.AfterFunc(duration, func() {
				cmd.Process.Kill()
				log.Printf("Command timed out: %s", strings.Join(cmdArgs, " "))
			})
			defer timer.Stop()
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

// enumerateSMB 封装 Hexstrike-redteam 的 SMB 枚举工具 (例如 enum4linux-ng, smbmap)
func (s *RedteamToolsServer) EnumerateSMB(ctx context.Context, req *redteam_toolspb.SMBEnumerateRequest) (*redteam_toolspb.SMBEnumerateResponse, error) {
	log.Printf("RedteamTools: Enumerating SMB on target %s", req.Target)

	// 假设 Hexstrike-redteam 有一个统一的 SMB 枚举脚本，例如 `./scripts/smb_enum.sh` 或 `./tools/smb_enum.py`
	// 这里我们直接调用一个假想的命令
	output, err := s.runToolCommand(ctx, "tools/smb_enum.py", []string{"--target", req.Target, "--full", "true"})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "SMB enumeration failed: %v", err)
	}

	// 实际应解析输出，提取 shares 和 users
	shares := []string{"share1", "share2"}
	users := []string{"admin", "guest"}
	outPath := fmt.Sprintf("/tmp/smb_enum_%s.txt", req.Target)
	os.WriteFile(outputPath, []byte(output), 0644) // 保存原始输出

	return &redteam_toolspb.SMBEnumerateResponse{
		Success:   true,
		OutputPath: outputPath,
		Shares:    shares,
		Users:     users,
	}, nil
}

func (s *RedteamToolsServer) enumerateSMBWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &redteam_toolspb.SMBEnumerateRequest{
		Target:  params["target"],
		Options: params,
	}
	resp, err := s.EnumerateSMB(ctx, req)
	if err != nil {
		return &core.OperationResponse{Success: false, Error: err.Error()}, nil
	}
	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("SMB enumeration completed for %s. Shares: %v, Users: %v. Report: %s", req.Target, resp.Shares, resp.Users, resp.OutputPath),
		Artifacts: map[string]string{
			"smb_report_path": resp.OutputPath,
		},
	}, nil
}


// PerformLateralMovement 封装横向移动工具 (例如 psexec, wmiexec, ssh)
func (s *RedteamToolsServer) PerformLateralMovement(ctx context.Context, req *redteam_toolspb.LateralMovementRequest) (*redteam_toolspb.LateralMovementResponse, error) {
	log.Printf("RedteamTools: Performing lateral movement to %s via %s", req.Target, req.Method)

	// 假设 Hexstrike-redteam 有一个工具来处理横向移动，例如 `tools/lateral_move.py`
	args := []string{
		"--target", req.Target,
		"--method", req.Method,
		"--credentials", req.Credentials,
	}
	if req.PayloadPath != "" {
		args = append(args, "--payload", req.PayloadPath)
	}

	output, err := s.runToolCommand(ctx, "tools/lateral_move.py", args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lateral movement failed: %v", err)
	}

	// 实际应解析输出，判断是否成功建立会话
	success := strings.Contains(output, "Session established")
	sessionId := ""
	if success {
		sessionId = "session-" + generateRandomString(10) // 假设生成session ID
	}

	return &redteam_toolspb.LateralMovementResponse{
		Success:   success,
		Output:    output,
		SessionId: sessionId,
	}, nil
}

func (s *RedteamToolsServer) performLateralMovementWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &redteam_toolspb.LateralMovementRequest{
		Target:      params["target"],
		Method:      params["method"],
		Credentials: params["credentials"],
		PayloadPath: params["payload_path"],
	}
	resp, err := s.PerformLateralMovement(ctx, req)
	if err != nil {
		return &core.OperationResponse{Success: false, Error: err.Error()}, nil
	}
	return &core.OperationResponse{
		Success: resp.Success,
		Output:  fmt.Sprintf("Lateral movement to %s (%t). Session: %s", req.Target, resp.Success, resp.SessionId),
		Artifacts: map[string]string{
			"session_id": resp.SessionId,
		},
	}, nil
}

// DumpCredentials 封装凭据转储工具 (例如 mimikatz, secretsdump)
func (s *RedteamToolsServer) DumpCredentials(ctx context.Context, req *redteam_toolspb.CredentialDumpRequest) (*redteam_toolspb.CredentialDumpResponse, error) {
	log.Printf("RedteamTools: Dumping credentials from %s via %s", req.Target, req.Method)

	// 假设 Hexstrike-redteam 有一个工具来处理凭据转储，例如 `tools/cred_dump.py`
	args := []string{
		"--target", req.Target,
		"--method", req.Method,
	}
	if req.Credentials != "" {
		args = append(args, "--credentials", req.Credentials)
	}

	output, err := s.runToolCommand(ctx, "tools/cred_dump.py", args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Credential dump failed: %v", err)
	}

	// 实际应解析输出，提取凭据，并保存到文件
	dumpedCredentialsPath := fmt.Sprintf("/tmp/creds_dump_%s.txt", req.Target)
	os.WriteFile(dumpedCredentialsPath, []byte(output), 0644)
	
	// 简单统计行数作为凭据数量
	credCount := len(strings.Split(output, "\n")) -1

	return &redteam_toolspb.CredentialDumpResponse{
		Success:              true,
		DumpedCredentialsPath: dumpedCredentialsPath,
		CredentialCount:      int32(credCount),
	}, nil
}

func (s *RedteamToolsServer) dumpCredentialsWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &redteam_toolspb.CredentialDumpRequest{
		Target:      params["target"],
		Method:      params["method"],
		Credentials: params["credentials"],
	}
	resp, err := s.DumpCredentials(ctx, req)
	if err != nil {
		return &core.OperationResponse{Success: false, Error: err.Error()}, nil
	}
	return &core.OperationResponse{
		Success: resp.Success,
		Output:  fmt.Sprintf("Credentials dumped from %s. Count: %d. Path: %s", req.Target, resp.CredentialCount, resp.DumpedCredentialsPath),
		Artifacts: map[string]string{
			"credentials_path": resp.DumpedCredentialsPath,
		},
	}, nil
}

// ScanVulnerabilities 封装漏洞扫描工具 (例如 Nmap with Vuln scripts, Nuclei)
func (s *RedteamToolsServer) ScanVulnerabilities(ctx context.Context, req *redteam_toolspb.VulnerabilityScanRequest) (*redteam_toolspb.VulnerabilityScanResponse, error) {
	log.Printf("RedteamTools: Scanning vulnerabilities on target %s with profile %s", req.Target, req.ScanProfile)

	// 假设 Hexstrike-redteam 有一个封装 Nmap/Nuclei 的脚本，例如 `tools/vuln_scan.py`
	args := []string{
		"--target", req.Target,
		"--profile", req.ScanProfile,
	}
	if len(req.Ports) > 0 {
		args = append(args, "--ports", strings.Join(req.Ports, ","))
	}

	output, err := s.runToolCommand(ctx, "tools/vuln_scan.py", args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Vulnerability scan failed: %v", err)
	}

	reportPath := fmt.Sprintf("/tmp/vuln_scan_report_%s.json", req.Target)
	os.WriteFile(reportPath, []byte(output), 0644)
	
	// 实际应解析输出，提取漏洞数量和检测到的服务
	vulnCount := 5
	detectedServices := map[string]string{"80": "http", "443": "https"}

	return &redteam_toolspb.VulnerabilityScanResponse{
		Success:             true,
		ReportPath:          reportPath,
		VulnerabilitiesFound: int32(vulnCount),
		DetectedServices:    detectedServices,
	}, nil
}

func (s *RedteamToolsServer) scanVulnerabilitiesWrapper(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	req := &redteam_toolspb.VulnerabilityScanRequest{
		Target:      params["target"],
		ScanProfile: params["scan_profile"],
		Ports:       []string{params["ports"]}, // 简化
	}
	resp, err := s.ScanVulnerabilities(ctx, req)
	if err != nil {
		return &core.OperationResponse{Success: false, Error: err.Error()}, nil
	}
	return &core.OperationResponse{
		Success: resp.Success,
		Output:  fmt.Sprintf("Vulnerability scan completed for %s. Found %d vulnerabilities. Report: %s", req.Target, resp.VulnerabilitiesFound, resp.ReportPath),
		Artifacts: map[string]string{
			"vulnerability_report_path": resp.ReportPath,
		},
	}, nil
}


func (s *RedteamToolsServer) GetCapabilities(ctx context.Context, empty *core.Empty) (*core.CapabilityList, error) {
	return &core.CapabilityList{
		Commands: []string{
			"enumerate_smb",
			"perform_lateral_movement",
			"dump_credentials",
			"scan_vulnerabilities",
		},
	}, nil
}

func (s *RedteamToolsServer) HealthCheck(ctx context.Context, empty *core.Empty) (*core.HealthStatus, error) {
	// 检查 Hexstrike-redteam 的根目录是否存在
	if _, err := os.Stat(s.toolBasePath); os.IsNotExist(err) {
		return &core.HealthStatus{
			Status: "unhealthy",
			Version: "unknown",
		}, nil
	}
	
	// 可以进一步检查内部关键工具脚本是否存在
	return &core.HealthStatus{
		Status: "healthy",
		Version: "latest", // 从工具中获取实际版本
	}, nil
}

func main() {
	server := NewRedteamToolsServer()
	
	lis, err := net.Listen("tcp", ":50055") // 使用新端口
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	redteam_toolspb.RegisterRedteamToolsAdapterServer(grpcServer, server)
	
	log.Printf("RedteamTools Adapter Server listening on :50055")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// 辅助函数：生成随机字符串 (重复定义，实际应放在一个公共工具包)
func generateRandomString(length int) string {
	// 简化实现，实际应该使用安全的随机数生成
	return "random_id"
}
