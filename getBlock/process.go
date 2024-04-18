package getBlock

import (
	"getBlock/getBlock/parse"
	"getBlock/getBlock/requests"
	"sync"
)

func (g *GetBlock) processBlock(wg *sync.WaitGroup, chTransactions <-chan []requests.Transaction) {
	defer wg.Done()

	for transactions := range chTransactions {
		for _, transaction := range transactions {
			from, to := parse.Input(transaction.Input)

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
