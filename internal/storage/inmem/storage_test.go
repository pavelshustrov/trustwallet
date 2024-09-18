package inmem_test

import (
	"reflect"
	"testing"
	"trustwallet/internal/model"
	"trustwallet/internal/storage/inmem"
)

func TestInMemory_AddAddress_IsSubscribed(t *testing.T) {
	tests := []struct {
		name       string
		addresses  []model.Address // Addresses to add
		checkAddr  model.Address   // Address to check
		wantSubbed bool
	}{
		{
			name:       "Subscribe to new address",
			addresses:  []model.Address{"0xAddress1"},
			checkAddr:  "0xAddress1",
			wantSubbed: true,
		},
		{
			name:       "Check unsubscribed address",
			addresses:  []model.Address{},
			checkAddr:  "0xAddress2",
			wantSubbed: false,
		},
		{
			name:       "Subscribe to existing address",
			addresses:  []model.Address{"0xAddress1", "0xAddress1"},
			checkAddr:  "0xAddress1",
			wantSubbed: true,
		},
		{
			name:       "Subscribe to empty address",
			addresses:  []model.Address{""},
			checkAddr:  "",
			wantSubbed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := inmem.New()

			for _, addr := range tt.addresses {
				err := im.AddAddress(addr)
				if err != nil {
					t.Errorf("AddAddress() error = %v", err)
				}
			}

			got, err := im.IsSubscribed(tt.checkAddr)
			if err != nil {
				t.Errorf("IsSubscribed() error = %v", err)
			}
			if got != tt.wantSubbed {
				t.Errorf("IsSubscribed() = %v, want %v", got, tt.wantSubbed)
			}
		})
	}
}

func TestInMemory_AddTransaction_GetTransactions(t *testing.T) {
	tests := []struct {
		name             string
		transactions     []model.Transaction
		address          model.Address
		wantTransactions []model.Transaction
	}{
		{
			name: "Transactions for address1",
			transactions: []model.Transaction{
				{
					Hash:        "0xTxHash1",
					From:        "0xAddress1",
					To:          "0xAddress2",
					Value:       "100",
					BlockNumber: "1",
				},
				{
					Hash:        "0xTxHash2",
					From:        "0xAddress2",
					To:          "0xAddress1",
					Value:       "200",
					BlockNumber: "2",
				},
				{
					Hash:        "0xTxHash3",
					From:        "",
					To:          "0xAddress1",
					Value:       "300",
					BlockNumber: "3",
				},
			},
			address: "0xAddress1",
			wantTransactions: []model.Transaction{
				{
					Hash:        "0xTxHash1",
					From:        "0xAddress1",
					To:          "0xAddress2",
					Value:       "100",
					BlockNumber: "1",
				},
				{
					Hash:        "0xTxHash2",
					From:        "0xAddress2",
					To:          "0xAddress1",
					Value:       "200",
					BlockNumber: "2",
				},
				{
					Hash:        "0xTxHash3",
					From:        "",
					To:          "0xAddress1",
					Value:       "300",
					BlockNumber: "3",
				},
			},
		},
		{
			name: "Transactions for address2",
			transactions: []model.Transaction{
				{
					Hash:        "0xTxHash1",
					From:        "0xAddress1",
					To:          "0xAddress2",
					Value:       "100",
					BlockNumber: "1",
				},
				{
					Hash:        "0xTxHash2",
					From:        "0xAddress2",
					To:          "0xAddress1",
					Value:       "200",
					BlockNumber: "2",
				},
			},
			address: "0xAddress2",
			wantTransactions: []model.Transaction{
				{
					Hash:        "0xTxHash1",
					From:        "0xAddress1",
					To:          "0xAddress2",
					Value:       "100",
					BlockNumber: "1",
				},
				{
					Hash:        "0xTxHash2",
					From:        "0xAddress2",
					To:          "0xAddress1",
					Value:       "200",
					BlockNumber: "2",
				},
			},
		},
		{
			name:             "No transactions for address3",
			transactions:     []model.Transaction{},
			address:          "0xAddress3",
			wantTransactions: []model.Transaction{},
		},
		{
			name: "Transactions for empty address",
			transactions: []model.Transaction{
				{
					Hash:        "0xTxHash3",
					From:        "",
					To:          "0xAddress1",
					Value:       "300",
					BlockNumber: "3",
				},
			},
			address: "",
			wantTransactions: []model.Transaction{
				{
					Hash:        "0xTxHash3",
					From:        "",
					To:          "0xAddress1",
					Value:       "300",
					BlockNumber: "3",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := inmem.New()

			for _, tx := range tt.transactions {
				err := im.AddTransaction(tx.From, tx)
				if err != nil {
					t.Errorf("AddTransaction() error = %v", err)
				}
				err = im.AddTransaction(tx.To, tx)
				if err != nil {
					t.Errorf("AddTransaction() error = %v", err)
				}
			}

			got, err := im.GetTransactions(tt.address)
			if err != nil {
				t.Errorf("GetTransactions() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.wantTransactions) {
				t.Errorf("GetTransactions() = %v, want %v", got, tt.wantTransactions)
			}
		})
	}
}

func TestInMemory_Concurrency(t *testing.T) {
	im := inmem.New()

	address := model.Address("0xAddress1")
	tx := model.Transaction{
		Hash:        "0xTxHash1",
		From:        address,
		To:          "0xAddress2",
		Value:       "100",
		BlockNumber: "1",
	}

	// Simulate concurrent access
	done := make(chan bool)
	go func() {
		for i := 0; i < 1000; i++ {
			_ = im.AddAddress(address)
			_ = im.AddTransaction(address, tx)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			_, _ = im.IsSubscribed(address)
			_, _ = im.GetTransactions(address)
		}
		done <- true
	}()

	<-done
	<-done

	// Verify final state
	subscribed, err := im.IsSubscribed(address)
	if err != nil {
		t.Errorf("IsSubscribed() error = %v", err)
	}
	if !subscribed {
		t.Errorf("Expected address to be subscribed")
	}

	txs, err := im.GetTransactions(address)
	if err != nil {
		t.Errorf("GetTransactions() error = %v", err)
	}
	if len(txs) != 1000 {
		t.Errorf("Expected 1000 transactions, got %d", len(txs))
	}
}
