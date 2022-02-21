package node

import (
	"log"
	"time"

	"gopkg.in/jcelliott/turnpike.v2"
)

func (n *Node) init() {
	// this will initiate everything
	n.initWampServer()

	n.initMonitoring()
}

func (n *Node) initWampServer() {
	for client, err := turnpike.NewWebsocketClient(turnpike.JSON, wampURL, nil, nil, nil); ; {
		if err != nil {
			log.Println("Error initiating the client:", err)
			time.Sleep(time.Second * 5)
			continue
		}

		if _, err := client.JoinRealm(wampRealm, nil); err != nil {
			log.Println("Error joining to the realm:", err)
			time.Sleep(time.Second * 5)
			continue
		}
		n.wamp = client
		break
	}
}
