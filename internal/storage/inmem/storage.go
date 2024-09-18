package inmem

import (
	"sync"
	"trustwallet/internal/model"
)

type InMemory struct {
	subscribedAddresses map[model.Address]bool
	transactions        map[model.Address][]model.Transaction
	mu                  *sync.RWMutex
}

func New() *InMemory {
	return &InMemory{
		subscribedAddresses: make(map[model.Address]bool),
		transactions:        make(map[model.Address][]model.Transaction),
		mu:                  &sync.RWMutex{},
	}
}

func (im *InMemory) AddAddress(address model.Address) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	im.subscribedAddresses[address] = true

	return nil
}

func (im *InMemory) IsSubscribed(address model.Address) (bool, error) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	subscribed := im.subscribedAddresses[address]

	return subscribed, nil
}

func (im *InMemory) AddTransaction(address model.Address, tx model.Transaction) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	im.transactions[address] = append(im.transactions[address], tx)

	return nil
}

func (im *InMemory) GetTransactions(address model.Address) ([]model.Transaction, error) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	txs, ok := im.transactions[address]
	if !ok {
		return []model.Transaction{}, nil
	}

	return txs, nil
}
