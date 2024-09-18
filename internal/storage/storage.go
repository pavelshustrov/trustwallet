package storage

import "trustwallet/internal/model"

//go:generate mockery --name=Storage --case=underscore --output=./mocks
type Storage interface {
	AddAddress(address model.Address) error
	IsSubscribed(address model.Address) (bool, error)

	AddTransaction(address model.Address, tx model.Transaction) error
	GetTransactions(address model.Address) ([]model.Transaction, error)
}
