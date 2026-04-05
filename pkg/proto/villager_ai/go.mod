module github.com/OpenClaw-Security/Stealth-Core/pkg/proto/villager_ai

go 1.22

require (
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core v0.0.0-00010101000000-000000000000 // indirect
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/hexstrike v0.0.0-00010101000000-000000000000 // indirect
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners v0.0.0-00010101000000-000000000000 // indirect
	github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer v0.0.0-00010101000000-000000000000 // indirect
	google.golang.org/grpc v1.80.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/OpenClaw-Security/Stealth-Core/pkg/proto/core => ../core
replace github.com/OpenClaw-Security/Stealth-Core/pkg/proto/hexstrike => ../hexstrike
replace github.com/OpenClaw-Security/Stealth-Core/pkg/proto/netrunners => ../netrunners
replace github.com/OpenClaw-Security/Stealth-Core/pkg/proto/thesilencer => ../thesilencer
