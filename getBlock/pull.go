package getBlock

import (
	"context"
	"fmt"
	"getBlock/getBlock/requests"
	"sync"
	"time"
)

// TODO: dry

func (g *GetBlock) pullBlocks(wg *sync.WaitGroup, fromBlock, toBlock int64, chTxs chan<- []requests.Transaction) {
	defer wg.Done()

	time.Sleep(g.rpsLimit)
	blocks, err := requests.GetBlocksByNumber(fromBlock, toBlock, g.nodeEndpoint, g.client)
	if err != nil {
		fmt.Printf("Unable to get blocks %d-%d. Reason: %s\n", fromBlock, toBlock, err.Error())
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
		ctx, _ := context.WithTimeout(context.Background(), g.timeout)
		wg.Add(1)
		g.repullBlocks(ctx, wg, failed, chTxs)
	}
}

func (g *GetBlock) repullBlocks(ctx context.Context, wg *sync.WaitGroup, requiredBlocks []string, chTxs chan<- []requests.Transaction) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		fmt.Printf("Unable to pull blocks %v. Reason: %v\n", requiredBlocks, ctx.Err())
		return
	default:
		break
	}

	if len(requiredBlocks) == 0 {
		return
	}

	if len(requiredBlocks[:len(requiredBlocks)/2]) != 0 {

		time.Sleep(g.rpsLimit)
		blocksLeft, err := requests.GetConcreteBlocksByNumber(
			requiredBlocks[:len(requiredBlocks)/2],
			g.nodeEndpoint,
			g.client,
		)
		if err != nil {
			fmt.Printf("Unable to repull blocks %v. Reason: %s\n", len(requiredBlocks)/2, err.Error())
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
	}

	if len(requiredBlocks[len(requiredBlocks)/2:]) != 0 {

		time.Sleep(g.rpsLimit)
		blocksRight, err := requests.GetConcreteBlocksByNumber(
			requiredBlocks[len(requiredBlocks)/2:],
			g.nodeEndpoint,
			g.client,
		)
		if err != nil {
			fmt.Printf("Unable to repull blocks %v. Reason: %s\n", len(requiredBlocks)/2, err.Error())
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
}
