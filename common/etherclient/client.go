package etherclient

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func MakeETHClient(nodeEndpoint string) *ethclient.Client {
	configEndpoint := nodeEndpoint
	fmt.Println(configEndpoint)
	client, err := ethclient.Dial(configEndpoint)
	if err != nil {
		log.Fatalf("Could not connect to eth client")
	}
	return client
}
