package ethereum_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"trustwallet/internal/model"
	"trustwallet/internal/parser/ethereum"
	"trustwallet/internal/parser/ethereum/mocks"
	storagemocks "trustwallet/internal/storage/mocks"
)

func TestParser_GetCurrentBlock(t *testing.T) {
	parser := ethereum.New(123456, nil, nil)

	expectedBlock := 123456
	actualBlock := parser.GetCurrentBlock()

	assert.Equal(t, expectedBlock, actualBlock, "GetCurrentBlock() should return the current block number")
}

func TestParser_Subscribe(t *testing.T) {
	mockStorage := storagemocks.NewStorage(t)
	parser := ethereum.New(0, nil, mockStorage)

	testAddress := model.Address("0xTestAddress")

	mockStorage.On("AddAddress", testAddress).Return(nil)

	success := parser.Subscribe(testAddress)

	assert.True(t, success, "Subscribe() should return true on success")
	mockStorage.AssertExpectations(t)
}

func TestParser_Subscribe_Failure(t *testing.T) {
	mockStorage := storagemocks.NewStorage(t)
	parser := ethereum.New(0, nil, mockStorage)

	testAddress := model.Address("0xTestAddress")
	mockError := errors.New("storage error")

	mockStorage.On("AddAddress", testAddress).Return(mockError)

	success := parser.Subscribe(testAddress)

	assert.False(t, success, "Subscribe() should return false on failure")
	mockStorage.AssertExpectations(t)
}

func TestParser_GetTransactions(t *testing.T) {
	mockStorage := storagemocks.NewStorage(t)
	parser := ethereum.New(0, nil, mockStorage)

	testAddress := model.Address("0xTestAddress")
	expectedTransactions := []model.Transaction{
		{
			Hash:        "0xHash1",
			From:        testAddress,
			To:          "0xAddress2",
			Value:       "100",
			BlockNumber: "1",
		},
		{
			Hash:        "0xHash2",
			From:        "0xAddress3",
			To:          testAddress,
			Value:       "200",
			BlockNumber: "2",
		},
	}

	mockStorage.On("GetTransactions", testAddress).Return(expectedTransactions, nil)

	actualTransactions := parser.GetTransactions(testAddress)

	assert.Equal(t, expectedTransactions, actualTransactions, "GetTransactions() should return the correct transactions")
	mockStorage.AssertExpectations(t)
}

func TestParser_GetTransactions_Error(t *testing.T) {
	mockStorage := storagemocks.NewStorage(t)
	parser := ethereum.New(0, nil, mockStorage)

	testAddress := model.Address("0xTestAddress")
	mockError := errors.New("storage error")

	mockStorage.On("GetTransactions", testAddress).Return(nil, mockError)

	actualTransactions := parser.GetTransactions(testAddress)

	assert.Nil(t, actualTransactions, "GetTransactions() should return nil when storage returns an error")
	mockStorage.AssertExpectations(t)
}

func TestParser_StartParsing_InitialBlock(t *testing.T) {
	mockClient := mocks.NewEthereumClient(t)
	parser := ethereum.New(0, mockClient, nil)

	mockClient.On("GetLatestBlockNumber").Return(int64(100), nil)

	err := parser.StartParsing()

	assert.NoError(t, err)
	assert.Equal(t, 100, parser.GetCurrentBlock(), "currentBlock should be set to latestBlock when it is 0")
	mockClient.AssertExpectations(t)
}

func TestParser_StartParsing_ParseBlocks(t *testing.T) {
	mockClient := mocks.NewEthereumClient(t)
	mockStorage := storagemocks.NewStorage(t)
	parser := ethereum.New(98, mockClient, mockStorage)

	mockClient.On("GetLatestBlockNumber").Return(int64(100), nil)

	txsBlock99 := []model.Transaction{
		{
			Hash:        "0xHash99",
			From:        "0xSubscribedAddress",
			To:          "0xAddress2",
			Value:       "100",
			BlockNumber: "99",
		},
	}
	txsBlock100 := []model.Transaction{
		{
			Hash:        "0xHash100",
			From:        "0xAddress3",
			To:          "0xSubscribedAddress",
			Value:       "200",
			BlockNumber: "100",
		},
	}

	mockClient.On("GetTransactionsByBlockNumber", int64(99)).Return(txsBlock99, nil)
	mockClient.On("GetTransactionsByBlockNumber", int64(100)).Return(txsBlock100, nil)

	// Simulate subscribed address
	mockStorage.On("IsSubscribed", model.Address("0xSubscribedAddress")).Return(true, nil).Twice()
	mockStorage.On("IsSubscribed", mock.Anything).Return(false, nil)

	// Expect AddTransaction to be called
	mockStorage.On("AddTransaction", model.Address("0xSubscribedAddress"), mock.AnythingOfType("model.Transaction")).Return(nil).Twice()

	err := parser.StartParsing()

	assert.NoError(t, err)
	assert.Equal(t, 100, parser.GetCurrentBlock(), "currentBlock should be updated to latestBlock")
	mockClient.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestParser_StartParsing_ClientError(t *testing.T) {
	mockClient := mocks.NewEthereumClient(t)

	parser := ethereum.New(98, mockClient, nil)

	mockError := errors.New("client error")
	mockClient.On("GetLatestBlockNumber").Return(int64(0), mockError)

	err := parser.StartParsing()

	assert.EqualError(t, err, "client error")
	mockClient.AssertExpectations(t)
}
