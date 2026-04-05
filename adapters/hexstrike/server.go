package main

import (
	"context"
	"fmt"
	"log"
	"net"

	core "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core"
	hexstrike_pb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/hexstrike"
	"google.golang.org/grpc"
)

type server struct {
	hexstrike_pb.UnimplementedHexstrikeAdapterServer
}

func (s *server) Execute(ctx context.Context, req *core.OperationRequest) (*core.OperationResponse, error) {
	log.Printf("Executing Hexstrike operation: %s", req.Command)
	
	switch req.Command {
	case "generate_payload":
		return s.generatePayload(ctx, req.Params)
	case "send_phishing_email":
		return s.sendPhishingEmail(ctx, req.Params)
	case "exploit_vulnerability":
		return s.exploitVulnerability(ctx, req.Params)
	default:
		return nil, fmt.Errorf("unknown command: %s", req.Command)
	}
}

func (s *server) GeneratePayload(ctx context.Context, req *hexstrike_pb.GeneratePayloadRequest) (*hexstrike_pb.GeneratePayloadResponse, error) {
	log.Printf("Generating payload with type: %s", req.Type)
	
	// Placeholder implementation
	return &hexstrike_pb.GeneratePayloadResponse{
		PayloadPath: "/tmp/generated_payload.bin",
		Checksum:    "abc123def456",
		Size:        1024,
	}, nil
}

func (s *server) SendPhishingEmail(ctx context.Context, req *hexstrike_pb.PhishingRequest) (*hexstrike_pb.PhishingResponse, error) {
	log.Printf("Sending phishing email with template: %s", req.TemplateName)
	
	// Placeholder implementation
	return &hexstrike_pb.PhishingResponse{
		SentCount:        5,
		FailedRecipients: []string{},
		MessageId:        "msg-12345",
	}, nil
}

func (s *server) ExploitVulnerability(ctx context.Context, req *hexstrike_pb.ExploitRequest) (*hexstrike_pb.ExploitResponse, error) {
	log.Printf("Exploiting vulnerability: %s on target: %s", req.Vulnerability, req.Target)
	
	// Placeholder implementation
	return &hexstrike_pb.ExploitResponse{
		Success:   true,
		Output:    "Exploit completed successfully",
		SessionId: "session-67890",
	}, nil
}

func (s *server) generatePayload(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	log.Printf("Generating payload with type: %s", params["type"])
	
	req := &hexstrike_pb.GeneratePayloadRequest{
		Type:    params["type"],
		Lhost:   params["lhost"],
		Lport:   params["lport"],
		Options: params,
	}
	
	resp, err := s.GeneratePayload(ctx, req)
	if err != nil {
		return &core.OperationResponse{
			Success: false,
			Error:   fmt.Sprintf("Payload generation failed: %v", err),
		}, nil
	}
	
	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Payload generated: %s", resp.PayloadPath),
		Artifacts: map[string]string{
			"payload_path": resp.PayloadPath,
			"checksum":     resp.Checksum,
			"size":         fmt.Sprintf("%d", resp.Size),
		},
	}, nil
}

func (s *server) sendPhishingEmail(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	log.Printf("Sending phishing email with template: %s", params["template"])
	
	req := &hexstrike_pb.PhishingRequest{
		TemplateName:   params["template"],
		Subject:        params["subject"],
		AttachmentPath: params["attachment"],
		Recipients:     []string{params["recipients"]},
		SmtpConfig:     params,
	}
	
	resp, err := s.SendPhishingEmail(ctx, req)
	if err != nil {
		return &core.OperationResponse{
			Success: false,
			Error:   fmt.Sprintf("Phishing email failed: %v", err),
		}, nil
	}
	
	return &core.OperationResponse{
		Success: true,
		Output:  fmt.Sprintf("Phishing email sent: %d recipients", resp.SentCount),
		Artifacts: map[string]string{
			"sent_count": fmt.Sprintf("%d", resp.SentCount),
			"message_id": resp.MessageId,
		},
	}, nil
}

func (s *server) exploitVulnerability(ctx context.Context, params map[string]string) (*core.OperationResponse, error) {
	log.Printf("Exploiting vulnerability: %s", params["vulnerability"])
	
	req := &hexstrike_pb.ExploitRequest{
		Target:        params["target"],
		Vulnerability: params["vulnerability"],
		Payload:       params["payload"],
		Options:       params,
	}
	
	resp, err := s.ExploitVulnerability(ctx, req)
	if err != nil {
		return &core.OperationResponse{
			Success: false,
			Error:   fmt.Sprintf("Exploit failed: %v", err),
		}, nil
	}
	
	return &core.OperationResponse{
		Success: true,
		Output:  resp.Output,
		Artifacts: map[string]string{
			"success":    fmt.Sprintf("%v", resp.Success),
			"session_id": resp.SessionId,
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hexstrike_pb.RegisterHexstrikeAdapterServer(s, &server{})
	log.Printf("Hexstrike Adapter Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
