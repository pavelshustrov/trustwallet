package ethereum

import (
	"log"
	"sync"
	"trustwallet/internal/model"
	"trustwallet/internal/storage"
)

//go:generate mockery --name=EthereumClient --case=underscore --output=./mocks
type EthereumClient interface {
	GetLatestBlockNumber() (int64, error)
	GetTransactionsByBlockNumber(blockNumber int64) ([]model.Transaction, error)
}

type Parser struct {
	mu           *sync.RWMutex
	currentBlock int64
	client       EthereumClient
	storage      storage.Storage
}

func New(currentBlock int64, client EthereumClient, storage storage.Storage) *Parser {
	return &Parser{
		mu:           &sync.RWMutex{},
		currentBlock: currentBlock,
		client:       client,
		storage:      storage,
	}
}

func (p *Parser) GetCurrentBlock() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return int(p.currentBlock)
}

func (p *Parser) Subscribe(address model.Address) bool {
	if err := p.storage.AddAddress(address); err != nil {
		log.Println("Failed to subscribe to address", address, err)
		return false
	}

	return true
}

func (p *Parser) GetTransactions(address model.Address) []model.Transaction {
	transactions, err := p.storage.GetTransactions(address)
	if err != nil {
		log.Println("Error getting transactions from database", err)
		return nil
	}

	return transactions
}

func (p *Parser) StartParsing() error {
	latestBlock, err := p.client.GetLatestBlockNumber()
	if err != nil {
		return err
	}

	if p.currentBlock == 0 {
		p.mu.Lock()
		p.currentBlock = latestBlock
		p.mu.Unlock()
	}

	for blockNum := p.currentBlock + 1; blockNum <= latestBlock; blockNum++ {
		if err := p.parseBlock(blockNum); err != nil {
			return err
		}

		p.mu.Lock()
		p.currentBlock = blockNum
		p.mu.Unlock()
	}

	return nil
}

func (p *Parser) parseBlock(blockNumber int64) error {
	transactionsData, err := p.client.GetTransactionsByBlockNumber(blockNumber)
	if err != nil {
		return err
	}

	for _, tx := range transactionsData {
		isFromSubscribed, _ := p.storage.IsSubscribed(tx.From)
		isToSubscribed, _ := p.storage.IsSubscribed(tx.To)

		if isFromSubscribed {
			if err := p.storage.AddTransaction(tx.From, tx); err != nil {
				log.Println("Failed to add transaction", tx.From, tx.To, err)
			}
		}

		if isToSubscribed {
			if err := p.storage.AddTransaction(tx.To, tx); err != nil {
				log.Println("Failed to add transaction", tx.From, tx.To, err)
			}
		}
	}

	return nil
}
