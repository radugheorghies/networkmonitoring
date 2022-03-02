package blockchaintest

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"math/big"
	token "networkmonitoring/pkg/blockchaintest/contracts"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/jcelliott/turnpike.v2"
)

type BlockchainTest struct {
	wamp           *turnpike.Client // connection to wamp server
	ethClient      *ethclient.Client
	publicKey      crypto.PublicKey
	publicKeyECDSA *ecdsa.PublicKey
	privateKey     *ecdsa.PrivateKey
	errorf         *os.File
	successf       *os.File
	addresses      Addresses
	trTime         TrTime
	trSuccess      uint64
	trFailed       uint64
	instance       *token.Token
	chainID        *big.Int
	trChan         chan Transaction
}

type TrTime struct {
	sync.Mutex
	values map[time.Duration]int
}

type Addresses struct {
	sync.Mutex
	array []string
}

type Transaction struct {
	Auth    *bind.TransactOpts
	Addr    common.Address
	Val     *big.Int
	TrTime  time.Time
	Context context.Context
}
