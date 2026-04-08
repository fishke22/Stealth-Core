module github.com/OpenClaw-Security/Stealth-Core/pkg/adapters/netrunners

go 1.26.1

replace github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core => /home/kali/.openclaw/workspace/chimera_project/pkg/proto/core

replace github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners => /home/kali/.openclaw/workspace/chimera_project/pkg/proto/netrunners

require (
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core v0.0.0-00010101000000-000000000000
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.80.0
)

require (
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
