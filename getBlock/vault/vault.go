package vault

import (
	"strings"
	"sync"
)

type Vault struct {
	sync.RWMutex
	vault map[string]int
}

func NewVault() *Vault {
	return &Vault{
		vault: make(map[string]int, 1000),
	}
}

func (v *Vault) Set(key string, value int) {
	key = strings.ToLower(key)
	v.Lock()
	defer v.Unlock()

	v.vault[key] = value
}

func (v *Vault) Get(key string) int {
	key = strings.ToLower(key)
	v.RLock()
	defer v.RUnlock()

	return v.vault[key]
}

func (v *Vault) GetVault() map[string]int {
	v.RLock()
	defer v.RUnlock()

	return v.vault
}
