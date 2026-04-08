package redteamtools

import (
	"context"

	core "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core"
	redteam_toolspb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/redteam_tools"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client redteam_toolspb.RedteamToolsAdapterClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := redteam_toolspb.NewRedteamToolsAdapterClient(conn)
	return &Client{conn: conn, client: client}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) EnumerateSMB(target string, options map[string]string) (*redteam_toolspb.SMBEnumerateResponse, error) {
	req := &redteam_toolspb.SMBEnumerateRequest{
		Target:  target,
		Options: options,
	}
	return c.client.EnumerateSMB(context.Background(), req)
}

func (c *Client) PerformLateralMovement(target, method, credentials, payloadPath string, options map[string]string) (*redteam_toolspb.LateralMovementResponse, error) {
	req := &redteam_toolspb.LateralMovementRequest{
		Target:      target,
		Method:      method,
		Credentials: credentials,
		PayloadPath: payloadPath,
		Options:     options,
	}
	return c.client.PerformLateralMovement(context.Background(), req)
}

func (c *Client) DumpCredentials(target, method, credentials string, options map[string]string) (*redteam_toolspb.CredentialDumpResponse, error) {
	req := &redteam_toolspb.CredentialDumpRequest{
		Target:      target,
		Method:      method,
		Credentials: credentials,
		Options:     options,
	}
	return c.client.DumpCredentials(context.Background(), req)
}

func (c *Client) ScanVulnerabilities(target string, ports []string, scanProfile string, options map[string]string) (*redteam_toolspb.VulnerabilityScanResponse, error) {
	req := &redteam_toolspb.VulnerabilityScanRequest{
		Target:      target,
		Ports:       ports,
		ScanProfile: scanProfile,
		Options:     options,
	}
	return c.client.ScanVulnerabilities(context.Background(), req)
}

func (c *Client) ExecuteOperation(opReq *core.OperationRequest) (*core.OperationResponse, error) {
	return c.client.Execute(context.Background(), opReq)
}
