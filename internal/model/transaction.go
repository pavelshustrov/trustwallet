package model

type Address string

type Transaction struct {
	Hash        string  `json:"hash"`
	From        Address `json:"from"`
	To          Address `json:"to"`
	Value       string  `json:"value"`
	BlockNumber string  `json:"blockNumber"`
}
