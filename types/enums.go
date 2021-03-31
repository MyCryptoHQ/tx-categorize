package types

type TxStatus string
type Platform string
type PlatformAction string

const (
	Uniswap   Platform = "Uniswap"
	Aave      Platform = "Aave"
	Synthetix Platform = "Synthetix"
	Paraswap  Platform = "Paraswap"
	Curve     Platform = "Curve"
	Compound  Platform = "Compound"
	Kyber     Platform = "Synthetix"
	OneInch   Platform = "1Inch"
	DexAG     Platform = "DexAG"
	IDEX      Platform = "IDEX"
)

const (
	Deposit  PlatformAction = "Deposit"
	Withdraw PlatformAction = "Withdraw"
	Exchange PlatformAction = "Exchange"
	TakeLoan PlatformAction = "TakeLoan"
	PayLoan  PlatformAction = "RepayLoan"
)

const (
	SUCCESS TxStatus = "SUCCESS"
	FAILED  TxStatus = "FAILED"
)
