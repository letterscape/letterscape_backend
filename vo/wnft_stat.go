package vo

import "github.com/letterScape/backend/constants/transaction"

type DealPriceStatVO struct {
	DealTime  string `json:"dealTime"`
	DealPrice string `json:"dealPrice"`
}

type TradeRecordVO struct {
	TradeId   string `json:"tradeId"`
	WnftId    string `json:"wnftId"`
	Buyer     string `json:"buyer"`
	Seller    string `json:"seller"`
	DealPrice string `json:"dealPrice"`
	DealTime  string `json:"dealTime"`
}

type TradePageVO struct {
	Total int64            `json:"total"`
	List  *[]TradeRecordVO `json:"list"`
}

type TransactionVO struct {
	TxId       string               `json:"txId"`
	DetailId   string               `json:"detailId"`
	TxHash     string               `json:"txHash"`
	TxStatus   transaction.TxStatus `json:"txStatus"`
	TxType     string               `json:"txType"`
	CreateTime string               `json:"createTime"`
}

type TransactionPageVO struct {
	Total int64            `json:"total"`
	List  *[]TransactionVO `json:"list"`
}
