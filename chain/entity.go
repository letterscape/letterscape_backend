package chain

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type WNFT struct {
	IsPaid        bool            `json:"isPaid"`
	IsListed      bool            `json:"isListed"`
	IsExpired     bool            `json:"isExpired"`
	Interval      *big.Int        `json:"interval"`
	Owner         *common.Address `json:"owner"`
	Deadline      *big.Int        `json:"deadline"`
	TokenId       *big.Int        `json:"tokenId"`
	Price         *big.Int        `json:"price"`
	LastDealPrice *big.Int        `json:"lastDealPrice"`
}
