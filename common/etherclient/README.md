# mycryptohq/infrastructure/common/etherclient

A simple shared pkg to simplify ethereum node interaction and conversion

To use, first initialize an ethclient. Then use it to make calls.
```go
client := etherclient.MakeETHClient("https://localhost:8000")
bal, err := etherclient.GetBalance(client, address)
if err != nil {
	return float64(0), err
}
return etherclient.ConvertFromWei(bal), nil
```