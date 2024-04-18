package getBlock

import (
	"fmt"
	"getBlock/getBlock/aggregation"
	"getBlock/getBlock/requests"
	"getBlock/getBlock/vault"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	defaultTimeout  = 300 * time.Second
	defaultRpsLimit = 50 * time.Millisecond
)

type GetBlock struct {
	nodeEndpoint string
	timeout      time.Duration
	rpsLimit     time.Duration
	depth        int64
	pullingStep  int64

	vault  *vault.Vault
	client *http.Client
}

func New(endpointEnv, apiKeyEnv, timeoutEnv, depthEnv, pullingStepEnv, rpsLimitEnv string) *GetBlock {
	timeout, err := time.ParseDuration(timeoutEnv)
	if err != nil {
		fmt.Printf("Unable to parse GET_BLOCK_TIMEOUT (has set up default 300s). Reason: %s\n", err.Error())
		timeout = defaultTimeout
	}

	depth, err := strconv.Atoi(depthEnv)
	if err != nil {
		depth = 100
	}

	pullingStep, err := strconv.Atoi(pullingStepEnv)
	if err != nil {
		fmt.Printf("Unable parse GET_BLOCK_PULLING_STEP (has set up default 50). Reason: %s\n", err.Error())
		pullingStep = 10
	}

	rpsLimit, err := time.ParseDuration(rpsLimitEnv)
	if err != nil {
		fmt.Printf("Unable parse GET_BLOCK_RPS_LIMIT (has set up default 50ms). Reason: %s\n", err.Error())
		rpsLimit = defaultRpsLimit
	}

	return &GetBlock{
		timeout:      timeout,
		nodeEndpoint: "https://" + endpointEnv + "/" + apiKeyEnv,
		depth:        int64(depth),
		pullingStep:  int64(pullingStep),
		rpsLimit:     rpsLimit,

		vault:  vault.NewVault(),
		client: &http.Client{},
	}
}

const top5 = 5
const someWeight = 5

var maxParallels = someWeight * runtime.GOMAXPROCS(0)

func (g *GetBlock) GetTop5Users() []aggregation.TopAddresses {
	latestBlock, err := requests.GetBlockNumber(g.nodeEndpoint, g.client)
	if err != nil {
		fmt.Printf("Unable to get last block from chain. Reason: %s", err.Error())
		return nil
	}

	chTransactions := make(chan []requests.Transaction, g.depth)

	wgProcessing := &sync.WaitGroup{}
	for i := 0; i < maxParallels; i++ {
		wgProcessing.Add(1)
		go g.processBlock(wgProcessing, chTransactions)
	}

	latestBlockInt := latestBlock.Int64()

	wgPulling := &sync.WaitGroup{}
	for currentBlock := latestBlockInt; currentBlock > latestBlockInt-g.depth; currentBlock -= g.pullingStep {
		wgPulling.Add(1)
		go g.pullBlocks(wgPulling, currentBlock-g.pullingStep, currentBlock, chTransactions)

		time.Sleep(g.rpsLimit)
	}

	wgPulling.Wait()

	close(chTransactions)

	wgProcessing.Wait()

	topFive, err := aggregation.GetTopAddresses(top5, g.vault.GetVault())
	if err != nil {
		fmt.Printf("Unable to execute request. Reason: %s", err.Error())
	}

	return topFive
}
