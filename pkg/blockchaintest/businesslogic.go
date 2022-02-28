package blockchaintest

import (
	"context"
	"crypto/ecdsa"
	"encoding/csv"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	token "networkmonitoring/pkg/blockchaintest/contracts"
	"networkmonitoring/pkg/core/env"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var wg sync.WaitGroup

func (b *BlockchainTest) startTest() {
	var err error
	var ok bool

	// starting aggregator
	if err := b.wamp.Publish("startAggregator", nil, []interface{}{}, nil); err != nil {
		log.Println("Problem occurred while publishing start commands to nodes:", err)
	}
	time.Sleep(time.Second * 2)
	// open the file to write erros and successfull opperations
	b.errorf, err = os.Create("errors.csv")
	if err != nil {
		log.Println("Error creating the errors file:", err)
		return
	}
	defer b.errorf.Close()

	b.successf, err = os.Create("success.csv")
	if err != nil {
		log.Println("Error creating the success file:", err)
		return
	}
	defer b.successf.Close()

	// reading the private key from file
	prKey, err := os.ReadFile("private.key")
	if err != nil {
		log.Println("Error reading the file:", err)
		return
	}
	privateKeyString := strings.ReplaceAll(string(prKey), " ", "")

	// take json addresses from file

	records := readCsvFile("addresses.csv")
	b.addresses.array = make([]string, len(records))
	for i := 0; i < len(records); i++ {
		b.addresses.array[i] = records[i][0]
	}

	// defining the client
	if b.ethClient, err = ethclient.Dial("https://mainnet-rpc.tlxscan.com/"); err != nil {
		log.Fatal(err)
	}

	if b.privateKey, err = crypto.HexToECDSA(privateKeyString); err != nil {
		log.Fatal(err)
	}

	b.publicKey = b.privateKey.Public()
	if b.publicKeyECDSA, ok = b.publicKey.(*ecdsa.PublicKey); !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	address := common.HexToAddress("0xa45Ee53F50769d347002C6c66f35E77aD25d84bd")
	b.instance, err = token.NewToken(address, b.ethClient)
	if err != nil {
		log.Fatal(err)
	}

	b.chainID, err = b.ethClient.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// preparing workers
	wg.Add(env.Vars.Workers)
	for i := 0; i < env.Vars.Workers; i++ {
		go b.startWorker()
	}

	wg.Wait()

	// stop the aggregator
	if err := b.wamp.Publish("stopAggregator", nil, []interface{}{}, nil); err != nil {
		log.Println("Problem occurred while publishing start commands to nodes:", err)
	}

	log.Println("TEST RESULTS:")
	log.Println("Successfull transactions:", b.trSuccess)
	log.Println("Failed transactions:", b.trFailed)

	log.Println("Transactions by time:")
	for k, v := range b.trTime.values {
		log.Println("time:", k, "- number of transactions:", v)
	}
}

func (b *BlockchainTest) startWorker() {
	defer wg.Done()

	for i := 0; i < env.Vars.ProcessPerWorker; i++ {
		trTime := time.Now()
		context := context.Background()
		// do the magic here
		fromAddress := crypto.PubkeyToAddress(*b.publicKeyECDSA)
		nonce, err := b.ethClient.PendingNonceAt(context, fromAddress)
		if err != nil {
			log.Fatal(err)
		}

		gasPrice, err := b.ethClient.SuggestGasPrice(context)
		if err != nil {
			log.Fatal(err)
		}

		auth, _ := bind.NewKeyedTransactorWithChainID(b.privateKey, b.chainID)
		auth.Nonce = big.NewInt(int64(nonce))
		auth.Value = big.NewInt(0)     // in wei
		auth.GasLimit = uint64(300000) // in units
		auth.GasPrice = gasPrice

		// choosing the address
		rand.Seed(time.Now().UnixNano())
		recipient := rand.Intn(len(b.addresses.array) - 1)

		tx, err := b.instance.Transfer(auth, common.HexToAddress(b.addresses.array[recipient]), big.NewInt(10000000000000000))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("tx sent: %s", tx.Hash().Hex())
		go b.getTrResponse(trTime, context, tx.Hash())
	}
}

func (b *BlockchainTest) getTrResponse(trTime time.Time, context context.Context, tx common.Hash) {
	// now we need to wait until we receiving the status of the transaction
	for {
		time.Sleep(100 * time.Millisecond)
		// fmt.Printf("Waiting receipt of transaction %s\n", tx.Hash().Hex())
		if !b.IsTransactionPending(context, tx) {
			receipt, err := b.ethClient.TransactionReceipt(context, tx)
			if err != nil {
				log.Println(err)
				break
			}
			if receipt.Status == 1 {
				atomic.AddUint64(&b.trSuccess, 1)
			} else {
				atomic.AddUint64(&b.trFailed, 1)
			}
			break
		}
	}

	endTime := time.Since(trTime)

	// reagister the time

	b.trTime.Lock()
	if _, ok := b.trTime.values[endTime]; ok {
		b.trTime.values[endTime] = b.trTime.values[endTime] + 1
	} else {
		b.trTime.values[endTime] = 1
	}

	b.trTime.Unlock()
}

func (b *BlockchainTest) IsTransactionPending(context context.Context, hash common.Hash) bool {
	_, pending, err := b.ethClient.TransactionByHash(context, hash)
	if err != nil {
		panic(err)
	}
	return pending
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
