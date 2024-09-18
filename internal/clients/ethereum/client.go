package ethereum

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"trustwallet/internal/model"
)

var ErrRPC = errors.New("RPC Error")

type Client struct {
	url    string
	client *http.Client
}

func New(url string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		url:    url,
		client: httpClient,
	}
}

func (c *Client) GetLatestBlockNumber() (int64, error) {
	rawJson, err := c.call("eth_blockNumber", []interface{}{})
	if err != nil {
		return 0, err
	}

	var blockHex string
	if err := json.Unmarshal(rawJson, &blockHex); err != nil {
		return 0, err
	}

	blockNumber, err := strconv.ParseInt(blockHex[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

func (c *Client) GetTransactionsByBlockNumber(blockNumber int64) ([]model.Transaction, error) {
	blockHex := "0x" + strconv.FormatInt(blockNumber, 16)

	rawJson, err := c.call("eth_getBlockByNumber", []interface{}{blockHex, true})
	if err != nil {
		return nil, err
	}

	var blockResp Block
	if err := json.Unmarshal(rawJson, &blockResp); err != nil {
		return nil, err
	}

	return blockResp.Transactions, nil
}

func (c *Client) call(method string, params interface{}) (json.RawMessage, error) {
	rpcReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}

	reqBytes, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp Response
	err = json.Unmarshal(respBytes, &rpcResp)
	if err != nil {
		return nil, err
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("%w: %d: %s", ErrRPC, rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}
