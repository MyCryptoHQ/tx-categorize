package types

type TxStatus string
type Protocol string
type ProtocolAction string

const (
	SUCCESS TxStatus = "SUCCESS"
	FAILED  TxStatus = "FAILED"
)
