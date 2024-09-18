package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"trustwallet/internal/clients/ethereum"
	"trustwallet/internal/model"
	ethereumParser "trustwallet/internal/parser/ethereum"
	"trustwallet/internal/storage/inmem"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	inmemStorage := inmem.New()

	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	ethereumClient := ethereum.New("https://ethereum-rpc.publicnode.com", httpClient)

	parser := ethereumParser.New(0, ethereumClient, inmemStorage)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		log.Println("Parser started")
		defer wg.Done()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := parser.StartParsing(); err != nil {
					log.Println("error parsing: ", err)
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		log.Println("Print Subscribed Transactions")
		defer wg.Done()

		var address model.Address = "0xYourEthereumAddress"
		parser.Subscribe(address)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				transactions := parser.GetTransactions(address)
				for _, transaction := range transactions {
					log.Println(transaction.BlockNumber)
				}
			}
		}
	}()

	wg.Wait()
}
