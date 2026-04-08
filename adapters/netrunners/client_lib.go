package netrunners

import (
	"context"

	core "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core"
	netrunnerspb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client netrunnerspb.NetRunnersAdapterClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := netrunnerspb.NewNetRunnersAdapterClient(conn)
	return &Client{conn: conn, client: client}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SetupListener(protocol, domain, port string) (string, error) {
	req := &netrunnerspb.ListenerRequest{
		Protocol: protocol,
		Domain:   domain,
		Port:     port,
	}

	resp, err := c.client.SetupListener(context.Background(), req)
	if err != nil {
		return "", err
	}

	return resp.ConnectionString, nil
}

func (c *Client) ExfiltrateData(dataPath, method, destination string) error {
	req := &netrunnerspb.ExfiltrationRequest{
		DataPath:    dataPath,
		Method:      method,
		Destination: destination,
	}

	_, err := c.client.ExfiltrateData(context.Background(), req)
	return err
}

func (c *Client) ExecuteOperation(opReq *core.OperationRequest) (*core.OperationResponse, error) {
	return c.client.Execute(context.Background(), opReq)
}
