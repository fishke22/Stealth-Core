module github.com/OpenClaw-Security/Stealth-Core/core/orchestrator

go 1.22

require (
	google.golang.org/grpc v1.80.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto => ../pkg/proto
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core => ../pkg/proto/core
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/hexstrike => ../pkg/proto/hexstrike
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer => ../pkg/proto/thesilencer
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners => ../pkg/proto/netrunners
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/villager_ai => ../pkg/proto/villager_ai
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/redteam_tools => ../pkg/proto/redteam_tools
	github.com/OpenClaw-Security/Stealth-Core/pkg/workflow => ../pkg/workflow
)
