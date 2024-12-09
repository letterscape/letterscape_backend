package chain

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"testing"
)

func TestFp(t *testing.T) {
	var hostname string = "http://localhost:3000"
	var originUri string = "http://localhost:3000/space/0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266/1"
	var positionId string = "1"

	websiteEncode, err := rlp.EncodeToBytes([]interface{}{"string", hostname})
	if err != nil {
		println(err)
	}

	websiteHex := hexutil.Encode(crypto.Keccak256(websiteEncode))
	websiteFp := new(big.Int)
	websiteFp, _ = websiteFp.SetString(websiteHex[2:10], 16)
	println("websiteFp: ", websiteFp.Text(16))

	originUriEncode, err := rlp.EncodeToBytes([]interface{}{"string", originUri})
	if err != nil {
		println(err)
	}
	originUriHex := hexutil.Encode(crypto.Keccak256(originUriEncode))
	originUriFp := new(big.Int)
	originUriFp, _ = originUriFp.SetString(originUriHex[2:10], 16)
	println("originUriFp: ", originUriFp.Text(16))

	positionIdFp := new(big.Int)
	positionIdFp, _ = positionIdFp.SetString(positionId, 10)

	fp := new(big.Int).Add(new(big.Int).Add(new(big.Int).Lsh(websiteFp, 64), new(big.Int).Lsh(originUriFp, 32)), positionIdFp)
	println("fp: ", fp.Text(16))
}
