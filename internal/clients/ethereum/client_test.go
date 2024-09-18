package ethereum_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"trustwallet/internal/clients/ethereum"
	"trustwallet/internal/model"
)

func TestClient_GetLatestBlockNumber(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		bodyStr := string(bodyBytes)

		assert.Contains(t, bodyStr, `"method":"eth_blockNumber"`)

		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(`{
            "jsonrpc": "2.0",
            "id": 1,
            "result": "0x10d4f"
        }`))

		assert.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := ethereum.New(server.URL, server.Client())

	blockNumber, err := client.GetLatestBlockNumber()
	assert.NoError(t, err)
	assert.Equal(t, int64(68943), blockNumber)
}

func TestClient_GetTransactionsByBlockNumber(t *testing.T) {
	blockNumberHex := "0x1"

	mockTransactions := []model.Transaction{
		{
			Hash:        "0xTxHash1",
			From:        "0xFromAddress",
			To:          "0xToAddress",
			Value:       "0x38d7ea4c68000",
			BlockNumber: "0x1",
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		blockResponse := ethereum.Block{
			Number:       blockNumberHex,
			Hash:         "0xBlockHash",
			Transactions: mockTransactions,
		}

		response := ethereum.Response{
			JSONRPC: "2.0",
			ID:      1,
			Result:  marshalJSON(t, blockResponse),
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(marshalJSON(t, response))
		assert.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := ethereum.New(server.URL, server.Client())

	transactions, err := client.GetTransactionsByBlockNumber(1)
	assert.NoError(t, err)
	assert.Equal(t, mockTransactions, transactions)
}

func marshalJSON(t *testing.T, v interface{}) []byte {
	data, err := json.Marshal(v)
	assert.NoError(t, err)
	return data
}

func TestClient_GetLatestBlockNumber_HTTPError(t *testing.T) {
	// Mock server that returns an HTTP error
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := ethereum.New(server.URL, server.Client())

	_, err := client.GetLatestBlockNumber()
	assert.Error(t, err)
}

func TestClient_GetLatestBlockNumber_RPCError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
            "jsonrpc": "2.0",
            "id": 1,
            "error": {
                "code": -32601,
                "message": "Method not found"
            }
        }`))
		assert.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := ethereum.New(server.URL, server.Client())

	_, err := client.GetLatestBlockNumber()
	assert.Error(t, err)
	assert.ErrorIs(t, err, ethereum.ErrRPC)
}

func TestClient_GetLatestBlockNumber_JSONUnmarshalError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`Invalid JSON`))
		assert.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := ethereum.New(server.URL, server.Client())

	_, err := client.GetLatestBlockNumber()
	assert.Error(t, err)
}
