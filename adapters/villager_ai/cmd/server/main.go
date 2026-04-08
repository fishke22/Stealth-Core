package main

import (
	"log"

	villagerai "github.com/OpenClaw-Security/Stealth-Core/adapters/villager_ai"
)

func main() {
	if err := villagerai.StartVillagerAIServer(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
