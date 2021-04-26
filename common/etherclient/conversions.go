package etherclient

import (
	"encoding/hex"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func ConvertFromWei(amount big.Int) float64 {
	inputAmountFloat := new(big.Float).SetInt(&amount)
	output := new(big.Float).Quo(inputAmountFloat, big.NewFloat(math.Pow10(int(18))))
	outputFloat, _ := output.Float64()
	return outputFloat
}

func ConvertFromBase(amount big.Int, decimal int) float64 {
	inputAmountFloat := new(big.Float).SetInt(&amount)
	output := new(big.Float).Quo(inputAmountFloat, big.NewFloat(math.Pow10(int(decimal))))
	outputFloat, _ := output.Float64()
	return outputFloat
}

func CalculateUserReadableGas(gasUsed *big.Int, gasPrice *big.Int) float64 {
	return ConvertFromWei(*new(big.Int).Mul(gasUsed, gasPrice))
}

func ConvertFromHexToBigInt(amount []byte) *big.Int {
	z := new(big.Int)
	z.SetBytes(amount)
	return z
}

func ConvertHashToString(hash common.Hash) string {
	trimmedHash := common.TrimLeftZeroes(hash.Bytes())
	return common.BytesToAddress(trimmedHash).String()
}

func ConvertFromHexByteToHexString(amount []byte) string {
	trimmedAmt := common.TrimLeftZeroes(amount)
	enc := make([]byte, len(trimmedAmt)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], trimmedAmt)
	return string(enc)
}
