package main

import (
	"log"

	redteamtools "github.com/OpenClaw-Security/Stealth-Core/adapters/redteam_tools"
)

func main() {
	if err := redteamtools.StartRedteamToolsServer(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
