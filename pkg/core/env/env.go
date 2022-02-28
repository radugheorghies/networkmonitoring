package env

import (
	"log"

	"github.com/caarlos0/env/v6"
)

// config env configuration
type config struct {
	WAMPURL          string `env:"WAMP_URL" envDefault:"ws://localhost:8089/"`
	WAMPRealm        string `env:"WAMP_REALM" envDefault:"dataGateway"`
	NodeName         string `env:"NODE_NAME" envDefault:"My node 1"`
	Workers          int    `env:"WORKERS" envDefault:"1000"`
	ProcessPerWorker int    `env:"PROCESS_PER_WORKER" envDefault:"100"`
}

var (
	// Vars keeps parsed env variables
	Vars config
)

func init() {
	// Parse env vars
	if err := env.Parse(&Vars); err != nil {
		log.Fatalln("Failed to parse env:", err)
	}
}
