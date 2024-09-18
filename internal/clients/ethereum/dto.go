package ethereum

import (
	"encoding/json"
	"trustwallet/internal/model"
)

type Block struct {
	Number       string              `json:"number"`
	Hash         string              `json:"hash"`
	Transactions []model.Transaction `json:"transactions"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *Error          `json:"error,omitempty"`
}
