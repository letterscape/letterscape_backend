package chain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/letterScape/backend/constants/transaction"
	"github.com/letterScape/backend/global"
	"log"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type EthereumOpts struct {
}

const nftAddress = "0x5147c5C1Cb5b5D3f56186C37a4bcFBb3Cd0bD5A7"
const marketAddress = "0xF2cb3cfA36Bfb95E0FD855C1b41Ab19c517FcDB9"
const spaceAddress = "0xC3549920b94a795D75E6C003944943D552C46F97"

// todo load operation put into start process

func loadMarketABI() (*abi.ABI, error) {
	marketABIFile := "chain/abi/market.abi"
	marketABIBytes, err := os.ReadFile(marketABIFile)
	if err != nil {
		log.Printf("cannot read ABI file: %v", err)
		return nil, err
	}
	marketABI, err := abi.JSON(bytes.NewReader(marketABIBytes))
	if err != nil {
		log.Printf("resolve abi failed: %v", err)
		return nil, err
	}
	return &marketABI, nil
}

func loadNftABI() (*abi.ABI, error) {
	nftABIFile := "chain/abi/lsnft.abi"
	nftABIBytes, err := os.ReadFile(nftABIFile)
	if err != nil {
		log.Printf("cannot read ABI file: %v", err)
		return nil, err
	}
	nftABI, err := abi.JSON(bytes.NewReader(nftABIBytes))
	if err != nil {
		log.Printf("resolve abi failed: %v", err)
		return nil, err
	}
	return &nftABI, nil
}

func (ether *EthereumOpts) QueryTxStatus(txHash common.Hash) (transaction.TxStatus, error) {
	client, err := ethclient.Dial(global.BlockChainConfig.RpcUrl)
	if err != nil {
		log.Fatal(err)
	}
	_, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return transaction.None, err
	}
	if isPending {
		return transaction.Pending, nil
	}

	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return transaction.None, err
	}
	if receipt.Status == 1 {
		return transaction.Success, nil
	} else {
		return transaction.Failed, nil
	}
}

func (ether *EthereumOpts) QueryTx(txHash common.Hash) (*transaction.TxData, error) {
	client, err := ethclient.Dial(global.BlockChainConfig.RpcUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal("call TransactionByHash error:", err)
		return nil, err
	}

	signer := types.NewLondonSigner(tx.ChainId())
	from, err := types.Sender(signer, tx)
	if err != nil {
		log.Fatal("resolve tx from error: ", err)
		return nil, err
	}

	data, err := parseMarketTxData(tx.Data())
	if err != nil {
		log.Fatal("resolve tx data error:", err)
		return nil, err
	}

	txData := &transaction.TxData{
		Data:    data,
		From:    from.String(),
		To:      tx.To().String(),
		Value:   tx.Value().String(),
		ChainId: tx.ChainId().String(),
		Time:    tx.Time(),
	}

	if isPending {
		txData.TxStatus = transaction.Pending
		return txData, nil
	}

	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		txData.TxStatus = transaction.None
		return txData, nil
	}
	if receipt.Status == 1 {
		txData.TxStatus = transaction.Success
		return txData, nil
	} else {
		txData.TxStatus = transaction.Failed
		return txData, nil
	}

	return txData, nil
}

func (ether *EthereumOpts) GetWNFT(tokenId string) (*WNFT, error) {
	// todo put rpc client into the global config
	client, err := ethclient.Dial(global.BlockChainConfig.RpcUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	marketAddr := common.HexToAddress(marketAddress)

	marketABI, err := loadMarketABI()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	tokenIdBig, err := hexutil.DecodeBig(tokenId)
	if err != nil {
		return nil, err
	}
	data, err := marketABI.Pack("getWNFT", tokenIdBig)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	msg := ethereum.CallMsg{To: &marketAddr, Data: data}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// todo need more test for situation of wnft not exists or burnt
	unpacked, err := marketABI.Methods["getWNFT"].Outputs.Unpack(result)
	if err != nil {
		return nil, err
	}
	unpackedString := fmt.Sprintf("%v", unpacked)
	if len(unpackedString) < 2 {
		return nil, errors.New("unpack result error")
	}
	values := strings.Split(unpackedString[2:len(unpackedString)-2], " ")
	wnft := &WNFT{}
	if len(values) != reflect.TypeOf(wnft).Elem().NumField() {
		return nil, errors.New("unpack values len error")
	}
	wnft.IsPaid, err = strconv.ParseBool(values[0])
	if err != nil {
		return nil, err
	}
	wnft.IsListed, err = strconv.ParseBool(values[1])
	if err != nil {
		return nil, err
	}
	wnft.IsExpired, err = strconv.ParseBool(values[2])
	if err != nil {
		return nil, err
	}
	wnft.Interval = new(big.Int)
	wnft.Interval.SetString(values[3], 10)

	owner := common.HexToAddress(values[4])
	wnft.Owner = &(owner)

	wnft.Deadline = new(big.Int)
	wnft.Deadline.SetString(values[5], 10)

	wnft.TokenId = new(big.Int)
	wnft.TokenId.SetString(values[6], 10)

	wnft.Price = new(big.Int)
	wnft.Price.SetString(values[7], 10)

	wnft.LastDealPrice = new(big.Int)
	wnft.LastDealPrice.SetString(values[8], 10)

	return wnft, nil
}

func (ether *EthereumOpts) GetTokenURI(fp string) (string, error) {
	client, err := ethclient.Dial(global.BlockChainConfig.RpcUrl)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	nftAddr := common.HexToAddress(nftAddress)

	nftABI, err := loadNftABI()
	if err != nil {
		log.Printf("abi load error: ", err)
		return "", err
	}

	fpBig := new(big.Int)
	fpBig.SetString(fp[2:], 16)
	log.Printf("fp: %s", fpBig.Text(16))

	data, err := nftABI.Pack("getTokenURI", fpBig)
	if err != nil {
		log.Println(err)
	}
	msg := ethereum.CallMsg{To: &nftAddr, Data: data}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}
	unpacked, err := nftABI.Methods["getTokenURI"].Outputs.Unpack(result)
	if err != nil {
		log.Println(err)
	}
	unpackedStr := fmt.Sprintf("%v", unpacked)
	tokenURI := unpackedStr[1 : len(unpackedStr)-1]
	log.Printf("fp: %s, tokenURI: %s", fp, tokenURI)
	return tokenURI, nil
}

func parseMarketTxData(data []byte) ([]interface{}, error) {
	funcSelector := data[:4]

	marketABI, err := loadMarketABI()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	method, err := marketABI.MethodById(funcSelector)
	if err != nil {
		return nil, err
	}

	args := data[4:]
	unpacked, err := method.Inputs.Unpack(args)
	if err != nil {
		return nil, err
	}

	return unpacked, nil
}
