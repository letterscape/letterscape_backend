package chain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/letterScape/backend/constants/transaction"
)

var Mapping = func() map[string](ChainOperation) {
	m := make(map[string]ChainOperation)
	m["1"] = &EthereumOpts{}
	return m
}()

type ChainOperation interface {
	QueryTxStatus(txHash common.Hash) (transaction.TxStatus, error)

	QueryTx(txHash common.Hash) (*transaction.TxData, error)

	GetWNFT(tokenId string) (*WNFT, error)

	GetTokenURI(fp string) (string, error)
}

type Context struct {
	chainOpt ChainOperation
}

func (c *Context) SetChainOpt(chainOpt ChainOperation) {
	c.chainOpt = chainOpt
}

func (c *Context) QueryTxStatus(txHash common.Hash) (transaction.TxStatus, error) {
	return c.chainOpt.QueryTxStatus(txHash)
}

func (c *Context) QueryTx(txHash common.Hash) (*transaction.TxData, error) {
	return c.chainOpt.QueryTx(txHash)
}

func (c *Context) GetWNFT(tokenId string) (*WNFT, error) {
	return c.chainOpt.GetWNFT(tokenId)
}

func (c *Context) GetTokenURI(fp string) (string, error) {
	return c.chainOpt.GetTokenURI(fp)
}
