package blockchaintest

import (
	"networkmonitoring/pkg/core/env"
	"time"
)

func New() *BlockchainTest {
	return &BlockchainTest{
		trTime: TrTime{
			values: make(map[time.Duration]int),
		},
		trChan: make(chan Transaction, env.Vars.Workers),
	}
}

func (b *BlockchainTest) Run() {
	b.init()

	b.startTest()
}
