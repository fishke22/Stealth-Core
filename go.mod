module github.com/OpenClaw-Security/Stealth-Core

go 1.26.1

require (
	github.com/OpenClaw-Security/Stealth-Core/core/orchestrator v0.0.0-00010101000000-000000000000
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto v0.0.0-00010101000000-000000000000
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core v0.0.0-00010101000000-000000000000
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/hexstrike v0.0.0-00010101000000-000000000000
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners v0.0.0-00010101000000-000000000000
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/redteam_tools v0.0.0-00010101000000-000000000000
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer v0.0.0-00010101000000-000000000000
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/villager_ai v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.80.0
)

require (
	github.com/OpenClaw-Security/Stealth-Core/pkg/workflow v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/OpenClaw-Security/Stealth-Core/core/orchestrator => ./core/orchestrator
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto => ./pkg/proto
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core => ./pkg/proto/core
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/hexstrike => ./pkg/proto/hexstrike
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners => ./pkg/proto/netrunners
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/redteam_tools => ./pkg/proto/redteam_tools
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer => ./pkg/proto/thesilencer
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/villager_ai => ./pkg/proto/villager_ai
	github.com/OpenClaw-Security/Stealth-Core/pkg/workflow => ./pkg/workflow
)
