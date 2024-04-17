package getBlock

import (
	"context"
	"fmt"
	"getBlock/getBlock/aggregation"
	"getBlock/getBlock/processing"
	"getBlock/getBlock/requests"
	"getBlock/getBlock/vault"
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

const depth = 50
const someWeight = 1

var maxParallels = someWeight * runtime.GOMAXPROCS(0)

const five = 5
const blockPullingStep = 100
const rpsLimiter = 10 * time.Millisecond
const repullTimeout = 60 * time.Second

func (g *GetBlock) GetTop5Addresses() []aggregation.TopAddresses {
	latestBlock, err := requests.GetBlockNumber(g.nodeEndpoint, g.client)
	if err != nil {
		fmt.Printf("Unable to get last block from chain. Reason: %s", err.Error())
		return nil
	}

	chTransactions := make(chan []requests.Transaction, depth)

	wgPulling := &sync.WaitGroup{}

	for i := 0; i < depth; i += blockPullingStep {
		wgPulling.Add(1)
		go g.processBlock(wgPulling, chTransactions)
	}

	wgProcessing := &sync.WaitGroup{}

	for currentBlock := latestBlock.Int64(); currentBlock > latestBlock.Int64()-depth; currentBlock -= blockPullingStep {
		wgProcessing.Add(1)
		go g.pullBlocks(wgProcessing, currentBlock, currentBlock+blockPullingStep, chTransactions)
		time.Sleep(rpsLimiter)
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

const triesRepull = 3

func (g *GetBlock) pullBlocks(wg *sync.WaitGroup, fromBlock, toBlock int64, chTxs chan<- []requests.Transaction) {
	defer wg.Done()

	blocks, err := requests.GetBlocksByNumber(fromBlock, toBlock, g.nodeEndpoint, g.client)
	if err != nil {
		fmt.Printf("Unable to get blocks %d-%d. Reason: %s", fromBlock, toBlock, err.Error())
	}

	failed := make([]string, 0, len(blocks))
	for _, block := range blocks {
		if block.Result.Hash == "" {
			failed = append(failed, block.ID)
			continue
		}

		chTxs <- block.Result.Transactions
	}

	if len(failed) != 0 {
		ctx, _ := context.WithTimeout(context.Background(), repullTimeout)
		wg.Add(1)
		g.repullBlocks(ctx, wg, failed, chTxs)
		//defer cancel()
	}
}

func (g *GetBlock) repullBlocks(ctx context.Context, wg *sync.WaitGroup, requiredBlocks []string, chTxs chan<- []requests.Transaction) {
	defer wg.Done()

	time.Sleep(1 * time.Second)
	select {
	case <-ctx.Done():
		fmt.Printf("Unable to repull blocks. Reason: %v\n", ctx.Err())
		return
	default:
		break
	}

	if len(requiredBlocks) == 0 {
		return
	}

	if len(requiredBlocks) == 1 {
		blocks, err := requests.GetConcreteBlocksByNumber(
			requiredBlocks,
			g.nodeEndpoint,
			g.client,
		)
		if err != nil {
			fmt.Printf("Unable to repull blocks %v. Reason: %s", len(requiredBlocks)/2, err.Error())
		}

		failed := make([]string, 0, len(blocks))
		for _, block := range blocks {
			if block.Result.Hash == "" {
				failed = append(failed, block.ID)
				continue
			}

			chTxs <- block.Result.Transactions
		}

		if len(failed) != 0 {
			wg.Add(1)
			go g.repullBlocks(ctx, wg, failed, chTxs)
		}

		return
	}

	fmt.Printf("left slice: %v\n\n", requiredBlocks[:len(requiredBlocks)/2])
	blocksLeft, err := requests.GetConcreteBlocksByNumber(
		requiredBlocks[:len(requiredBlocks)/2],
		g.nodeEndpoint,
		g.client,
	)
	if err != nil {
		fmt.Printf("Unable to repull blocks %v. Reason: %s", len(requiredBlocks)/2, err.Error())
	}

	failedLeft := make([]string, 0, len(blocksLeft))
	for _, block := range blocksLeft {
		if block.Result.Hash == "" {
			failedLeft = append(failedLeft, block.ID)
			continue
		}

		chTxs <- block.Result.Transactions
	}

	if len(failedLeft) != 0 {
		wg.Add(1)
		go g.repullBlocks(ctx, wg, failedLeft, chTxs)
	}

	fmt.Printf("right slice: %v\n\n", requiredBlocks[len(requiredBlocks)/2:])
	blocksRight, err := requests.GetConcreteBlocksByNumber(
		requiredBlocks[len(requiredBlocks)/2:],
		g.nodeEndpoint,
		g.client,
	)
	if err != nil {
		fmt.Printf("Unable to repull blocks %v. Reason: %s", len(requiredBlocks)/2, err.Error())
	}

	failedRight := make([]string, 0, len(blocksRight))
	for _, block := range blocksRight {
		if block.Result.Hash == "" {
			failedRight = append(failedRight, block.ID)
			continue
		}

		chTxs <- block.Result.Transactions
	}

	if len(failedRight) != 0 {
		wg.Add(1)
		go g.repullBlocks(ctx, wg, failedRight, chTxs)
	}
}

func (g *GetBlock) processBlock(wg *sync.WaitGroup, chTransactions <-chan []requests.Transaction) {
	defer wg.Done()

	for transactions := range chTransactions {
		for _, transaction := range transactions {
			fmt.Printf("Processing transaction %s", transaction.Hash)
			from, to := processing.ParseInput(transaction.Input)

			if from == "" && to == "" {
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
