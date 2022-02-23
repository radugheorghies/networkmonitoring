package env

import (
	"log"

	"github.com/caarlos0/env/v6"
)

// config env configuration
type config struct {
	WAMPURL   string `env:"WAMP_URL" envDefault:"ws://localhost:8087/"`
	WAMPRealm string `env:"WAMP_REALM" envDefault:"dataGateway"`
	NodeName  string `env:"NODE_NAME" envDefault:"My node 1"`
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
