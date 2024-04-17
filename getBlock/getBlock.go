package getBlock

import (
	"fmt"
	"icescream/getBlockStd/getBlock/aggregation"
	"icescream/getBlockStd/getBlock/processing"
	"icescream/getBlockStd/getBlock/requests"
	"icescream/getBlockStd/getBlock/vault"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type GetBlock struct {
	nodeEndpoint string

	vault  *vault.Vault
	client *http.Client
}

func New(nodeEndpoint string) *GetBlock {
	return &GetBlock{
		nodeEndpoint: nodeEndpoint,
		vault:        vault.NewVault(),
		client:       &http.Client{},
	}
}

const depth = 100
const someWeight = 1

var maxParallels = someWeight * runtime.GOMAXPROCS(0)

const five = 5

const upperBorderOfGreed = 25 * time.Millisecond

func (g *GetBlock) GetTop5Addresses() []aggregation.TopAddresses {
	latestBlock, err := requests.GetBlockNumber(g.nodeEndpoint, g.client)
	if err != nil {
		fmt.Printf("Unable to get last block from chain. Reason: %s", err.Error())
		return nil
	}

	chTransactions := make(chan []requests.Transaction, depth)

	wgPulling := &sync.WaitGroup{}

	for i := 0; i < depth; i++ {
		wgPulling.Add(1)
		go g.processBlock(wgPulling, chTransactions)
	}

	wgProcessing := &sync.WaitGroup{}

	for currentBlock := latestBlock.Int64(); currentBlock > latestBlock.Int64()-depth; currentBlock-- {
		wgProcessing.Add(1)
		go g.transactionsToStream(wgProcessing, currentBlock, chTransactions)
		time.Sleep(upperBorderOfGreed)
	}

	wgProcessing.Wait()

	close(chTransactions)

	wgPulling.Wait()

	topFive, err := aggregation.GetTopAddresses(five, g.vault.GetVault())
	if err != nil {
		fmt.Printf("Unable to execute request. Reason: %s", err.Error())
	}

	return topFive
}

func (g *GetBlock) transactionsToStream(wg *sync.WaitGroup, currentBlock int64, chTx chan []requests.Transaction) {
	defer wg.Done()

	blockTransactions, err := requests.GetBlockByNumber(currentBlock, g.nodeEndpoint, g.client)
	if err != nil {
		fmt.Printf("Unable to get block %d. Reason: %s", currentBlock, err.Error())
	}

	chTx <- blockTransactions
}

func (g *GetBlock) processBlock(wg *sync.WaitGroup, chTransactions <-chan []requests.Transaction) {
	defer wg.Done()

	for transactions := range chTransactions {
		for _, transaction := range transactions {
			from, to := processing.ParseInput(transaction.Input)

			if from == "" && to == "" {
				continue
			}

			// TODO: remove case below. Newer use
			if to == "" {
				fmt.Printf("unknown recipient in transaction %s", transaction.Hash)
				continue
			}

			if from == "" {
				from = transaction.From
			}

			countFrom := g.vault.Get(from)
			g.vault.Set(from, countFrom+1)

			countTo := g.vault.Get(to)
			g.vault.Set(to, countTo+1)
		}
	}
}