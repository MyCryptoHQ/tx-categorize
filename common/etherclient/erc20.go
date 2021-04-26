package etherclient

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func BalanceOf(client *ethclient.Client, tokenObject TokenObject) (*big.Int, error) {

	tokenCaller, err := NewTokenCaller(tokenObject.Contract, client)
	if err != nil {
		return big.NewInt(0), err
	}

	balance, err := tokenCaller.BalanceOf(nil, tokenObject.Wallet)
	if err != nil {
		return big.NewInt(0), err
	}
	return balance, nil
}

func TotalSupply(client *ethclient.Client, tokenObject TokenObject) (*big.Int, error) {

	tokenCaller, err := NewTokenCaller(tokenObject.Contract, client)
	if err != nil {
		return big.NewInt(0), err
	}

	totalSupply, err := tokenCaller.TotalSupply(nil)
	if err != nil {
		return big.NewInt(0), err
	}
	return totalSupply, nil
}

func Decimals(client *ethclient.Client, tokenObject TokenObject) (uint8, error) {

	tokenCaller, err := NewTokenCaller(tokenObject.Contract, client)
	if err != nil {
		return 0, err
	}

	decimals, err := tokenCaller.Decimals(nil)
	if err != nil {
		return 0, err
	}
	return decimals, nil
}
