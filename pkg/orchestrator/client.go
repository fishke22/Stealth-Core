package orchestrator

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	corepb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core"          // 导入核心 proto
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client corepb.OrchestratorClient // Changed proto.OrchestratorClient to corepb.OrchestratorClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Orchestrator at %s: %w", addr, err)
	}
	
	client := corepb.NewOrchestratorClient(conn) // Changed proto.NewOrchestratorClient to corepb.NewOrchestratorClient
	return &Client{conn: conn, client: client}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ExecuteWorkflow sends a workflow request and receives a stream of responses.
func (c *Client) ExecuteWorkflow(ctx context.Context, workflowReq *corepb.WorkflowRequest) error { // Changed proto.WorkflowRequest to corepb.WorkflowRequest
	stream, err := c.client.ExecuteWorkflow(ctx, workflowReq)
	if err != nil {
		return fmt.Errorf("failed to execute workflow: %w", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break // Stream finished
		}
		if err != nil {
			return fmt.Errorf("failed to receive workflow response: %w", err)
		}
		log.Printf("Workflow Status: OpID=%s, Status=%s, Success=%t, Output=%s, Error=%s",
			resp.OperationId, resp.Status, resp.Response.Success, resp.Response.Output, resp.Response.Error)
		// 可以根据需要处理 Artifacts
		for k, v := range resp.Response.Artifacts {
			log.Printf("  Artifact: %s = %s", k, v)
		}
	}
	return nil
}

// ExecuteOperation sends a single operation request.
func (c *Client) ExecuteOperation(ctx context.Context, opReq *corepb.OperationRequest) (*corepb.OperationResponse, error) { // Changed proto.OperationRequest, proto.OperationResponse to corepb.OperationRequest, corepb.OperationResponse
	resp, err := c.client.ExecuteOperation(ctx, opReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute operation: %w", err)
	}
	return resp, nil
}

// GetStatus retrieves the health status of all connected tools.
func (c *Client) GetStatus(ctx context.Context) (*corepb.StatusResponse, error) { // Changed proto.StatusRequest, proto.StatusResponse to corepb.StatusRequest, corepb.StatusResponse
	resp, err := c.client.GetStatus(ctx, &corepb.StatusRequest{}) // Changed proto.StatusRequest to corepb.StatusRequest
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}
	return resp, nil
}
