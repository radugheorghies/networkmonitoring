package blockchaintest

import (
	"log"
	"networkmonitoring/pkg/core/env"
	"time"

	"gopkg.in/jcelliott/turnpike.v2"
)

func (b *BlockchainTest) init() {
	// this will initiate everything
	b.initWampServer()
}

func (b *BlockchainTest) initWampServer() {
	for {
		client, err := turnpike.NewWebsocketClient(turnpike.JSON, env.Vars.WAMPURL, nil, nil, nil)
		if err != nil {
			log.Println("Error initiating the client:", err)
			time.Sleep(time.Second * 5)
			continue
		}

		if _, err := client.JoinRealm(env.Vars.WAMPRealm, nil); err != nil {
			log.Println("Error joining to the realm:", err)
			time.Sleep(time.Second * 5)
			continue
		}
		b.wamp = client
		break
	}
}
