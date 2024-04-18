package aggregation

type TopAddresses struct {
	Address string
	Count   int
}

func GetTopAddresses(topAmount int, vault map[string]int) (top []TopAddresses, err error) {
	if topAmount > len(vault) {
		topAmount = len(vault)
	}

	top = []TopAddresses{}

	for i := 0; i < topAmount; i++ {
		var maxAddress TopAddresses

		for address, count := range vault {
			if count > maxAddress.Count {
				maxAddress.Address = address
				maxAddress.Count = count
			}
		}

		top = append(top, maxAddress)
		vault[maxAddress.Address] = 0
	}

	return
}
