package blockchaintest

import "time"

func New() *BlockchainTest {
	return &BlockchainTest{
		trTime: TrTime{
			values: make(map[time.Duration]int),
		},
	}
}

func (b *BlockchainTest) Run() {
	b.init()

	b.startTest()
}
