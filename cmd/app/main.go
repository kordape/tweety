package main

import (
	"log"

	"github.com/kordape/tweety/config"
	"github.com/kordape/tweety/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
