package transaction

import "time"

type TxObject string

const (
	Wnft TxObject = "wnft"
)

// txType enum
const (
	MINT = iota + 1
	LIST
	BUY
	BURN
	HOLDFEE
)

func Type2String(txType int) string {
	switch txType {
	case MINT:
		return "Mint"
	case LIST:
		return "List"
	case BUY:
		return "Buy"
	case BURN:
		return "Burn"
	case HOLDFEE:
		return "Holdfee"
	default:
		return "Unknown"
	}
}

type TxStatus string

const (
	None    TxStatus = "none"
	Pending TxStatus = "pending"
	Mined   TxStatus = "mined"
	Failed  TxStatus = "failed"
	Success TxStatus = "success"
)

type TxData struct {
	TxStatus TxStatus
	Data     []interface{}
	From     string
	To       string
	Value    string
	ChainId  string
	Time     time.Time
}
