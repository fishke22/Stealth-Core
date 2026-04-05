package villager_ai

import (
	"context"

	core "github.com/OpenClaw-Security/Stealth-Core/pkg/proto"
	villagerpb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/villager_ai"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	client villagerpb.VillagerAIAdapterClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	
	client := villagerpb.NewVillagerAIAdapterClient(conn)
	return &Client{conn: conn, client: client},
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GenerateStealthyPayload(req *villagerpb.GenerateStealthyPayloadRequest) (*villagerpb.GenerateStealthyPayloadResponse, error) {
	return c.client.GenerateStealthyPayload(context.Background(), req)
}

func (c *Client) LaunchCovertPhishing(req *villagerpb.LaunchCovertPhishingRequest) (*villagerpb.LaunchCovertPhishingResponse, error) {
	return c.client.LaunchCovertPhishing(context.Background(), req)
}

func (c *Client) AutomatedExfiltration(req *villagerpb.AutomatedExfiltrationRequest) (*villagerpb.AutomatedExfiltrationResponse, error) {
	return c.client.AutomatedExfiltration(context.Background(), req)
}

func (c *Client) DeployCovertC2AndPersistence(req *villagerpb.DeployC2PersistenceRequest) (*villagerpb.DeployC2PersistenceResponse, error) {
	return c.client.DeployCovertC2AndPersistence(context.Background(), req)
}

func (c *Client) ExecuteOperation(opReq *core.OperationRequest) (*core.OperationResponse, error) {
	return c.client.Execute(context.Background(), opReq)
}
